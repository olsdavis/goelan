package protocol

//TODO: Handle errors

import (
	"../log"
	"bytes"
	"encoding/binary"
	"encoding/json"
)

// Represents server list ping response.
type ServerListPing struct {
	Ver  Version `json:"version"`
	Pl   Players `json:"players"`
	Desc Chat    `json:"description"`
	Fav  string  `json:"favicon,omitempty"`
}

type Version struct {
	Name     string `json:"name"`
	Protocol uint32 `json:"protocol"`
}

type Players struct {
	Max    uint `json:"max"`
	Online uint `json:"online"`
}

type Chat struct {
	Text string `json:"text"`
}

type Response struct {
	data *bytes.Buffer
}

// Creates a new response.
func NewResponse() *Response {
	log.Debug("New response!")
	return &Response{new(bytes.Buffer)}
}

// Writes a boolean.
func (r *Response) WriteBoolean(b bool) *Response {
	if b {
		return r.WriteByte(1)
	} else {
		return r.WriteByte(0)
	}
}

// Writes a Chat JSON Object.
func (r *Response) WriteChat(obj string) *Response {
	return r.WriteJSON(Chat{obj})
}

// Writes the given object as a JSON string.
func (r *Response) WriteJSON(obj interface{}) *Response {
	j, _ := json.Marshal(obj)
	return r.WriteByteArray(j)
}

// Writes the given byte.
func (r *Response) WriteByte(b byte) *Response {
	r.data.Write([]byte{b})
	return r
}

func (r *Response) WriteUnsignedByte(b uint8) *Response {
	binary.Write(r.data, ByteOrder, b)
	return r
}

// Writes a varint.
func (r *Response) WriteVarint(i uint32) *Response {
	_, err := r.data.Write(Uvarint(i))
	if err != nil {
		panic(err)
	}
	return r
}

// Writes an integer.
func (r *Response) WriteInt(i int) *Response {
	binary.Write(r.data, ByteOrder, i)
	return r
}

// Writes a long.
func (r *Response) WriteLong(l int64) *Response {
	binary.Write(r.data, ByteOrder, l)
	return r
}

// Writes a byte array.
func (r *Response) WriteByteArray(b []byte) *Response {
	r.WriteVarint(uint32(len(b)))
	r.data.Write(b)
	return r
}

// Writes a string.
func (r *Response) WriteString(str string) *Response {
	log.Debug("String to write:", str)
	return r.WriteByteArray([]byte(str))
}

// Returns the raw packet created from the written bytes and the provided id.
func (r *Response) ToRawPacket(id uint64) *RawPacket {
	return NewRawPacket(id, r.data.Bytes(), nil)
}

func (r *Response) Clear() {
	r.data.Reset()
}
