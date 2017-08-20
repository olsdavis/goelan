package protocol

//TODO: Handle errors

import (
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

// NewResponse creates a new response.
func NewResponse() *Response {
	return &Response{data: new(bytes.Buffer)}
}

// WriteBoolean writes the given boolean to the current response.
func (r *Response) WriteBoolean(b bool) *Response {
	if b {
		return r.WriteByte(1)
	} else {
		return r.WriteByte(0)
	}
}

// WriteChat writes a Chat JSON Object to the current response.
func (r *Response) WriteChat(obj string) *Response {
	return r.WriteJSON(Chat{obj})
}

// WriteJSON writes the given interface as a JSON string to the current response.
// (Currently ignores failure.)
func (r *Response) WriteJSON(obj interface{}) *Response {
	j, _ := json.Marshal(obj)
	return r.WriteByteArray(j)
}

// WriteByte writes the given byte to the current response.
func (r *Response) WriteByte(b byte) *Response {
	r.data.Write([]byte{b})
	return r
}

// WriteUnsignedByte writes the given unsigned byte to the current response.
func (r *Response) WriteUnsignedByte(b uint8) *Response {
	binary.Write(r.data, ByteOrder, b)
	return r
}

// WriteVarInt writes the given VarInt to the current response.
func (r *Response) WriteVarint(i uint32) *Response {
	_, err := r.data.Write(Uvarint(i))
	if err != nil {
		panic(err)
	}
	return r
}

// WriteInt writes the given integer to the current response.
func (r *Response) WriteInt(i int32) *Response {
	binary.Write(r.data, ByteOrder, i)
	return r
}

// WriteLong writes the given long to the current response.
func (r *Response) WriteLong(l int64) *Response {
	binary.Write(r.data, ByteOrder, l)
	return r
}

// WriteByteArray writes the given byte array to the current response.
func (r *Response) WriteByteArray(b []byte) *Response {
	r.WriteVarint(uint32(len(b)))
	r.data.Write(b)
	return r
}

// WriteString writes the given string to the current response.
func (r *Response) WriteString(str string) *Response {
	return r.WriteByteArray([]byte(str))
}

// Returns the raw packet created from the written bytes and the provided id.
func (r *Response) ToRawPacket(id uint64) *RawPacket {
	return NewRawPacket(id, r.data.Bytes(), nil)
}

// Clear clears the data from the response's buffer.
func (r *Response) Clear() {
	r.data = new(bytes.Buffer)
}
