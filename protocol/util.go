package protocol

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/olsdavis/goelan/util"
)

func Uvarint(n uint32) []byte {
	buf := make([]byte, 5)
	l := binary.PutUvarint(buf, uint64(n))
	return buf[:l]
}

var oneVarintBuf []byte = []byte{1}

func Varint(n int32) []byte {
	if n == 1 {  // TODO: fix this
		return oneVarintBuf
	}
	buf := make([]byte, 5)
	l := binary.PutVarint(buf, int64(n))
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

	ret := br.Buf[br.read]
	br.read++
	return ret, nil
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
