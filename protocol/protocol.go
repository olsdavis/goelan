package protocol

import (
	"encoding/binary"
	"errors"
	"github.com/olsdavis/goelan/log"
	"sync"
)

var (
	ByteOrder  = binary.BigEndian
	emptyArray []byte
	// ByteReader

	ReadAllError = errors.New("reached the end of this reader's buffer")
	readerPool   = sync.Pool{New: func() interface{} {
		return &ByteReader{
			emptyArray,
			0,
		}
	}}
)

type Callback func()

type Packet struct {
	ID uint64
}

type RawPacket struct {
	Packet
	Data     *ByteReader
	Callback Callback // the callback is a function called when the packet is sent
}

func NewRawPacket(id uint64, data []byte, callback Callback) *RawPacket {
	b := readerPool.Get().(*ByteReader)
	b.SetData(data)
	return &RawPacket{
		Packet{id},
		b,
		callback,
	}
}

func (r *RawPacket) ReadByte() byte {
	b, err := r.Data.ReadByte()
	if err != nil {
		panic(err)
	}
	return b
}

func (r *RawPacket) ReadUnsignedByte() byte {
	var ubyte uint8
	err := binary.Read(r.Data, ByteOrder, &ubyte)
	if err != nil {
		log.Error("Could not read unsigned short:", err)
	}
	return ubyte
}

func (r *RawPacket) ReadBoolean() bool {
	return r.ReadByte() == 1
}

func (r *RawPacket) ReadUnsignedShort() uint16 {
	var short uint16
	err := binary.Read(r.Data, ByteOrder, &short)
	if err != nil {
		log.Error("Could not read unsigned short:", err)
	}
	return short
}

func (r *RawPacket) ReadVarint() int32 {
	i, err := binary.ReadVarint(r.Data)
	if err != nil {
		log.Error("Could not read varint:", err)
	}
	return int32(i)
}

func (r *RawPacket) ReadUnsignedVarint() uint32 {
	i, err := binary.ReadUvarint(r.Data)
	if err != nil {
		log.Error("Could not read unsigned varint:", err)
	}
	return uint32(i)
}

func (r *RawPacket) ReadLong() int64 {
	var long int64
	err := binary.Read(r.Data, ByteOrder, &long)
	if err != nil {
		log.Error("Could not read long:", err)
	}
	return long
}

func (r *RawPacket) ReadByteArray() []byte {
	size := r.ReadUnsignedVarint()
	buf := make([]byte, size)
	r.Data.Read(buf)
	return buf
}

func (r *RawPacket) ReadString() string {
	return string(r.ReadByteArray())
}

func (r *RawPacket) Release() {
	if r.Data == nil {
		panic("data is nil!")
	}
	r.Data.Reset()
	readerPool.Put(r.Data)
}

const (
	/*** PACKET IDs ***/

	// Handshake state
	HandshakePacketId = 0x00
	PingPacketId      = 0x01
	// Login state
	LoginStartPacketId           = 0x00
	LoginStateDisconnectPacketId = 0x00
	EncryptionRequestPacketId    = 0x01
	EncryptionResponsePacketId   = 0x01
	LoginSuccessPacketId         = 0x02
	// Play state
	ClientSettingsPacketId    = 0x04
	PluginMessagePacketId     = 0x09
	ChatPacketId              = 0x10
	KeepAliveIncomingPacketId = 0x0B
	KickPlayerPacketId        = 0x1A
	KeepAliveOutgoingPacketId = 0x1F
	ChunkDataPacketId         = 0x20
	JoinGamePacketId          = 0x23

	/*** PACKET CONSTS ***/
	HandshakeStatusNextState = 1
	HandshakeLoginNextState  = 2
)

// Chat
type MessageMode int

const (
	// chat message (only for players)
	ChatMessageMode MessageMode = iota
	// system message (what you should use)
	DefaultMessageMode
	// action bar message
	ActionBarMode
)
