package Gotcp

import "gotcp/Igotcp"

type Request struct {
	connector Igotcp.IConnector
	msg       Igotcp.IMessage
}

func (r *Request) GetConnector() Igotcp.IConnector {
	return r.connector
}

func (r *Request) GetMsgData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgId() uint32 {
	return r.msg.GetId()
}
