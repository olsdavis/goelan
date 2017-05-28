package server

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/olsdavis/goelan/log"
	"github.com/olsdavis/goelan/player"
	"github.com/olsdavis/goelan/protocol"
	"io"
	"net"
	"sync"
)

type ConnectionState int

var (
	ErrNotEnoughBytes = errors.New("not enough bytes to read")
	EmtpyArray        []byte
)

const (
	HandshakeState = iota
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

	connected bool
	sync.Mutex
}

// Creates a new connection from the given ReadWriter.
func NewConnection(socket net.Conn, server *Server) *Connection {
	return &Connection{
		server:          server,
		Writer:          socket,
		Reader:          NewFullReader(socket),
		writeChan:       make(chan *protocol.RawPacket),
		exitChan:        make(chan int, 1),
		ConnectionState: HandshakeState,
		VerifyToken:     EmtpyArray,
		VerifyUsername:  "",
		Player:          nil,
		connected:       true,
		Mutex:           sync.Mutex{},
	}
}

// Reads the next packet and returns the Packet object from the
// "protocol" package.
func (c *Connection) Next() (*protocol.RawPacket, error) {
	size, err := binary.ReadUvarint(c.Reader)

	if err != nil {
		return nil, err
	}

	if size <= 1 {
		return nil, nil
	}

	buffer := make([]byte, size)
	read, err := io.ReadAtLeast(c.Reader, buffer, int(size))
	if err != nil {
		// Error comes from here
		return nil, err
	} else if read < int(size) {
		return nil, ErrNotEnoughBytes
	}

	id, offset := binary.Uvarint(buffer)
	rawPacket := protocol.NewRawPacket(id, buffer[offset:], nil)

	return rawPacket, nil
}

// Write sends the given packet to the current connection.
func (c *Connection) Write(packet *protocol.RawPacket) {
	if packet == nil {
		return
	}
	c.writeChan <- packet
}

// Receives packets from the writeChan and sends them to client.
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
			} else {
				if packet.ID == protocol.JoinGamePacketId {
					log.Debug("Very here:", packet.Data)
				}
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
	defer packet.Release()
	send := new(bytes.Buffer)
	send.Write(protocol.Uvarint(uint32(packet.ID)))
	send.Write(packet.Data.Buf)
	return append(protocol.Uvarint(uint32(send.Len())), send.Bytes()...)
}

// Returns connection's server.
func (c *Connection) GetServer() *Server {
	return c.server
}

// Returns true if the client is connected.
func (c *Connection) IsConnected() bool {
	defer c.Unlock()
	c.Lock()
	return c.connected
}

// Sets the connected field's value.
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
	response.WriteChat(message)
	response.WriteByte(byte(mode))
	c.Write(response.ToRawPacket(protocol.ChatPacketId))
}

// Disconnects the current client for the given reason. (May be empty.)
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
		response.WriteChat(reason)
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
