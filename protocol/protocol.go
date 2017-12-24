package protocol

import (
	"encoding/binary"
	"errors"
	"github.com/olsdavis/goelan/log"
	"sync"
	"math"
	"github.com/olsdavis/goelan/world"
	"reflect"
	"github.com/olsdavis/goelan/util"
)

const (
	MaxByteArrayLength = 32767
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

// NewRawPacket creates a RawPacket from the packet's id, the data and
// the given Callback.
func NewRawPacket(id uint64, data []byte, callback Callback) *RawPacket {
	b := readerPool.Get().(*ByteReader)
	b.SetData(data)
	return &RawPacket{
		Packet{id},
		b,
		callback,
	}
}

// ReadByte reads a byte and returns it.
func (r *RawPacket) ReadByte() byte {
	b, err := r.Data.ReadByte()
	if err != nil {
		panic(err)
	}
	return b
}

// ReadUnsignedByte reads an unsigned byte and returns it.
func (r *RawPacket) ReadUnsignedByte() byte {
	var ubyte uint8
	err := binary.Read(r.Data, ByteOrder, &ubyte)
	if err != nil {
		log.Error("Could not read unsigned short:", err)
	}
	return ubyte
}

// ReadBoolean reads a boolean and returns it.
func (r *RawPacket) ReadBoolean() bool {
	return r.ReadByte() == 1
}

// ReadUnsignedShort reads an unsigned short and returns it.
func (r *RawPacket) ReadUnsignedShort() uint16 {
	var short uint16
	err := binary.Read(r.Data, ByteOrder, &short)
	if err != nil {
		log.Error("Could not read unsigned short:", err)
	}
	return short
}

// ReadVarint reads a Varint and returns it.
func (r *RawPacket) ReadVarint() int32 {
	i, err := binary.ReadVarint(r.Data)
	if err != nil {
		log.Error("Could not read varint:", err)
	}
	return int32(i)
}

// ReadUnsignedVarint reads an unsigned Varint and returns it.
func (r *RawPacket) ReadUnsignedVarint() uint32 {
	i, err := binary.ReadUvarint(r.Data)
	if err != nil {
		log.Error("Could not read unsigned varint:", err)
	}
	return uint32(i)
}

// ReadFloat reads a float32 and returns it.
func (r *RawPacket) ReadFloat() float32 {
	buf := make([]byte, 4)
	r.Data.Read(buf)
	return math.Float32frombits(ByteOrder.Uint32(buf))
}

// ReadDouble reads a float64 and returns it.
func (r *RawPacket) ReadDouble() float64 {
	buf := make([]byte, 8)
	r.Data.Read(buf)
	return math.Float64frombits(ByteOrder.Uint64(buf))
}

// ReadLong reads an int64 and returns it.
func (r *RawPacket) ReadLong() int64 {
	var long int64
	err := binary.Read(r.Data, ByteOrder, &long)
	if err != nil {
		log.Error("Could not read long:", err)
	}
	return long
}

// ReadByteArrayMax reads a byte array which's length
// cannot exceed max.
func (r *RawPacket) ReadByteArrayMax(max uint32) []byte {
	size := r.ReadUnsignedVarint()
	if size > max {
		// TODO: Kick client
		return []byte{0}
	}
	buf := make([]byte, size)
	r.Data.Read(buf)
	return buf
}

// ReadByteArray reads a byte array which's length cannot
// exceed MaxByteArrayLength.
func (r *RawPacket) ReadByteArray() []byte {
	return r.ReadByteArrayMax(MaxByteArrayLength)
}

// ReadStringMax reads a string which's length cannot exceed max.
func (r *RawPacket) ReadStringMax(max uint32) string {
	return string(r.ReadByteArrayMax(max))
}

// ReadString reads a string which's length cannot exceed MaxByteArrayLength.
func (r *RawPacket) ReadString() string {
	return string(r.ReadByteArray())
}

// ReadStructure reads the data as the given structure,
// and sets its pointer's value.
func (r *RawPacket) ReadStructure(impl *interface{}) {
	s := *impl
	switch s.(type) {
	case int8:
		*impl = r.ReadByte()
	case uint8:
		*impl = r.ReadUnsignedByte()
	case bool:
		*impl = r.ReadBoolean()
	case int:
	case int32:
		*impl = r.ReadVarint()
	case int64:
		//TODO: implement
		panic("unimplemented: varlong")
	case uint32:
		*impl = r.ReadUnsignedVarint()
	case float32:
		*impl = r.ReadFloat()
	case float64:
		*impl = r.ReadDouble()
	case string:
		*impl = r.ReadString()
	case []byte:
		*impl = r.ReadByteArray()
	case world.Location:
		//TODO: implement
		panic("unimplemented: world.Location")
	case util.UUID:
		//TODO: implement
		panic("unimplemented: util.UUID")
	default:
		t := reflect.ValueOf(impl)
		if t.CanInterface() {
		}
	}
}

// Release releases RawPacket's data and puts it back to the pool.
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
	TeleportConfirmPacketId               = 0x00
	IncomingChatPacketId                  = 0x02
	ClientStatusPacketId                  = 0x03
	ClientSettingsPacketId                = 0x04
	SpawnPlayerPacketId                   = 0x05
	OutgoingAnimationPacketId             = 0x06
	ClickWindowPacketId                   = 0x07
	CloseWindowPacketId                   = 0x08
	PluginMessagePacketId                 = 0x09
	KeepAliveIncomingPacketId             = 0x0B
	IncomingPlayerPositionAndLookPacketId = 0x0E
	OutgoingChatPacketId                  = 0x0F
	KickPlayerPacketId                    = 0x1A
	IncomingAnimationPacketId             = 0x1D
	KeepAliveOutgoingPacketId             = 0x1F
	ChunkDataPacketId                     = 0x20
	JoinGamePacketId                      = 0x23
	PlayerAbilitiesPacketId               = 0x2C
	PlayerListItemPacketId                = 0x2E
	OutgoingPlayerPositionAndLookPacketId = 0x2F

	/*** PACKET CONSTS ***/
	HandshakeStatusNextState = 1
	HandshakeLoginNextState  = 2

	PlayerListItemActionAddPlayer         = 0
	PlayerListItemActionUpdateGamemode    = 1
	PlayerListItemActionUpdateLatency     = 2
	PlayerListItemActionUpdateDisplayName = 3
	PlayerListItemActionRemovePlayer      = 4
)
