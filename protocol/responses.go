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

func NewResponse() *Response {
	return &Response{new(bytes.Buffer)}
}

func (r *Response) WriteChat(obj string) *Response {
	return r.WriteJSON(Chat{obj})
}

func (r *Response) WriteJSON(obj interface{}) *Response {
	json, _ := json.Marshal(obj)
	return r.WriteByteArray(json)
}

func (r *Response) WriteVarint(i uint32) *Response {
	_, err := r.data.Write(Uvarint(i))
	if err != nil {
		panic(err)
	}
	return r
}

func (r *Response) WriteLong(l int64) *Response {
	binary.Write(r.data, binary.BigEndian, l)
	return r
}

func (r *Response) WriteByteArray(b []byte) *Response {
	r.WriteVarint(uint32(len(b)))
	r.data.Write(b)
	return r
}

func (r *Response) WriteString(str string) *Response {
	return r.WriteByteArray([]byte(str))
}

func (r *Response) ToRawPacket(id uint64) *RawPacket {
	return NewRawPacket(id, r.data.Bytes(), nil)
}
