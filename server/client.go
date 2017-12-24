package server

import (
	"bytes"
	"encoding/binary"
	"github.com/olsdavis/goelan/log"
	"github.com/olsdavis/goelan/player"
	"github.com/olsdavis/goelan/protocol"
	. "github.com/olsdavis/goelan/util"
	"io"
	"net"
	"sync"
)

type ConnectionState int

var (
	emptyArray []byte
)

const (
	HandshakeState ConnectionState = iota
	LoginState
	PlayState
)

type FullReader struct {
	R      io.Reader
	oneBuf []byte
}

func NewFullReader(reader io.Reader) FullReader {
	return FullReader{
		reader,
		make([]byte, 1),
	}
}

func (reader FullReader) ReadByte() (byte, error) {
	_, err := reader.R.Read(reader.oneBuf)
	return reader.oneBuf[0], err
}

func (reader FullReader) Read(p []byte) (int, error) {
	return reader.R.Read(p)
}

// Connection struct represents a connected client.
type Connection struct {
	server          *Server
	Writer          io.WriteCloser
	Reader          FullReader
	writeChan       chan *protocol.RawPacket
	exitChan        chan int
	PacketHandler   stateHandler // the handler which depends on player's state
	ProtocolVersion uint32
	ConnectionState ConnectionState // current connection's state (handshake, login or play)

	VerifyToken    []byte // the verify token used in authentication
	VerifyUsername string // the verify username used in authentication
	SharedSecret   []byte // used for encrypting and decrypting data

	Player *player.Player

	PendingKeepAlives            *PendingList
	PendingTeleportConfirmations *PendingList

	connected bool
	sync.Mutex
}

// NewConnection creates a new connection from the given ReadWriter.
func NewConnection(socket net.Conn, server *Server) *Connection {
	return &Connection{
		server:                       server,
		Writer:                       socket,
		Reader:                       NewFullReader(socket),
		writeChan:                    make(chan *protocol.RawPacket),
		exitChan:                     make(chan int, 1),
		ConnectionState:              HandshakeState,
		VerifyToken:                  emptyArray,
		VerifyUsername:               "",
		Player:                       nil,
		PendingKeepAlives:            NewPendingList(),
		PendingTeleportConfirmations: NewPendingList(),
		connected:                    true,
		Mutex:                        sync.Mutex{},
	}
}

// Next reads the next packet and returns a RawPacket.
func (c *Connection) Next() (*protocol.RawPacket, error) {
	size, err := binary.ReadUvarint(c.Reader)

	if err != nil {
		return nil, err
	}

	if size <= 1 {
		return nil, nil
	}

	buffer := make([]byte, size)
	_, err = io.ReadAtLeast(c.Reader, buffer, int(size))
	if err != nil {
		return nil, err
	}
	id, offset := binary.Uvarint(buffer)
	rawPacket := protocol.NewRawPacket(id, buffer[offset:], nil)
	return rawPacket, nil
}

// Write enqueues the given packet to the current connection.
func (c *Connection) Write(packet *protocol.RawPacket) {
	if packet == nil {
		return
	}
	c.writeChan <- packet
}

// write receives packets from the writeChan channel and sends
// them to client.
func (c *Connection) write() {
	for {
		select {
		case <-c.exitChan:
			c.exitChan = nil
			return
		case packet := <-c.writeChan:
			if c.exitChan == nil {
				break
			}

			if packet == nil {
				continue
			}

			_, err := c.Writer.Write(toByteArray(packet))

			// omit this error
			if err != nil {
				break
			}

			if packet.Callback != nil {
				packet.Callback()
			}
		}
	}
}

// Creates a byte array from the given raw packet. Releases the packet at the end.
func toByteArray(packet *protocol.RawPacket) []byte {
	send := new(bytes.Buffer)
	send.Write(protocol.Uvarint(uint32(packet.ID)))
	send.Write(packet.Data.Buf)
	packet.Release()
	return append(protocol.Uvarint(uint32(send.Len())), send.Bytes()...)
}

// GetServer returns client's server.
func (c *Connection) GetServer() *Server {
	return c.server
}

// IsConnected returns true if the client is connected.
// Locks the Mutex for the operation and unlocks it at the end.
func (c *Connection) IsConnected() bool {
	defer c.Unlock()
	c.Lock()
	return c.connected
}

// SetConnected sets the connected field's value to b.
// Locks the Mutex for the operation and unlocks it at the end.
func (c *Connection) SetConnected(b bool) {
	c.Lock()
	c.connected = b
	c.Unlock()
}

// Sends the given message with the given mode.
// Message's mode depends on what you want to send:
// - ChatMessageMode (mode 0): used for players only;
// - DefaultMessageMode (mode 1): what you should use (system messages);
// - ActionBarMode (mode 2): if you want to send messages above the hotbar, use this mode.
func (c *Connection) SendMessage(message string, mode protocol.MessageMode) {
	response := protocol.NewResponse()
	response.WriteJSON(protocol.ChatComponent{Text: message})
	response.WriteUnsignedByte(byte(mode))
	c.Write(response.ToRawPacket(protocol.OutgoingChatPacketId))
}

// AddPlayers sends to the current client the packet which adds
// to his player list the given players.
func (c *Connection) AddPlayers(players []*player.Player) {
	packet := protocol.NewResponse()
	packet.WriteByte(protocol.PlayerListItemActionAddPlayer)
	packet.WriteUVarint(uint32(len(players)))
	for _, pl := range players {
		profile := pl.Profile

		packet.WriteUUID(*profile.RealUUID)
		packet.WriteString(profile.Name)
		packet.WriteVarint(int32(len(profile.Properties)))
		for _, property := range profile.Properties {
			packet.WriteString(property.Name)
			packet.WriteString(property.Value)
			if property.Signature == "" {
				packet.WriteBoolean(false)
			} else {
				packet.WriteBoolean(true)
				packet.WriteString(property.Signature)
			}
		}
		packet.WriteVarint(int32(pl.GameMode))
		packet.WriteVarint(0)
		packet.WriteBoolean(true)
		packet.WriteJSON(protocol.ChatComponent{Text: profile.Name})
	}
	c.Write(packet.ToRawPacket(protocol.PlayerListItemPacketId))
}

// Disconnect disconnects the current client for the given reason. (May be empty.)
func (c *Connection) Disconnect(reason string) {
	if !c.IsConnected() {
		return
	}

	if reason != "" {
		var packetId uint64

		switch c.ConnectionState {
		case PlayState:
			packetId = protocol.KickPlayerPacketId
		case LoginState:
			packetId = protocol.LoginStateDisconnectPacketId
		}

		response := protocol.NewResponse()
		response.WriteJSON(protocol.ChatComponent{Text: reason})
		rp := response.ToRawPacket(packetId)
		// waits the packet to be send; it prevents us from writing messages
		// while the socket is being closed
		done := make(chan int, 1)
		rp.Callback = func() {
			close(done)
		}
		c.Write(rp)
		<-done
	}

	c.exitChan <- 0
	err := c.Writer.Close()
	c.SetConnected(false)

	if err != nil {
		log.Error("Could not properly close the connection with a client:", err)
	}

	close(c.writeChan)
}
