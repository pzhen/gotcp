package Igotcp

type IMessage interface {
	GetId() uint32
	GetLen() uint32
	GetData() []byte
	SetId(uint32)
	SetLen(uint32)
	SetData([]byte)
}
