package protocol

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/olsdavis/goelan/log"
	"github.com/olsdavis/goelan/util"
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
	LoginSuccessPacketId         = 0x02
	EncryptionResponsePacketId   = 0x01
	EncryptionRequestPacketId    = 0x01
	// Play state
	ChatPacketId       = 0x10
	KickPlayerPacketId = 0x1A
	KeepAlivePacketId  = 0x1F
	ChunkDataPacketId  = 0x20
	JoinGamePacketId   = 0x23

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

/*** Util ***/

func Uvarint(n uint32) []byte {
	buf := make([]byte, 5)
	l := binary.PutUvarint(buf, uint64(n))

	return buf[:l]
}

type ByteReader struct {
	Buf  []byte
	read int
}

func (br *ByteReader) SetData(data []byte) {
	br.Buf = data
}

func (br *ByteReader) ReadByte() (byte, error) {
	if br.read >= len(br.Buf) {
		return 0, errors.New("out of bounds")
	}

	defer func() {
		br.read++
	}()
	return br.Buf[br.read], nil
}

func (br *ByteReader) Read(p []byte) (int, error) {
	remaining := br.Buf[br.read:]
	if len(remaining) == 0 {
		return 0, ReadAllError
	}
	tr := util.Min(len(remaining), len(p))
	for i := 0; i < tr; i++ {
		p[i] = remaining[i]
	}
	br.read += tr
	return tr, nil
}

func (br *ByteReader) Reset() {
	br.Buf = emptyArray
	br.read = 0
}

func (br ByteReader) String() string {
	return fmt.Sprintf("{read:%v, buf:%v}", br.read, br.Buf)
}
