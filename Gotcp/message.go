package Gotcp

import "gotcp/Igotcp"

type Message struct {
	Id   uint32
	Len  uint32
	Data []byte
}

func NewMessage(id uint32, data []byte) Igotcp.IMessage {
	return &Message{
		Id:   id,
		Len:  uint32(len(data)),
		Data: data,
	}
}

func (m *Message) GetId() uint32 {
	return m.Id
}

func (m *Message) GetLen() uint32 {
	return m.Len
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetId(id uint32) {
	m.Id = id
}

func (m *Message) SetData(data []byte) {
	m.Data = data
}

func (m *Message) SetLen(len uint32) {
	m.Len = len
}
