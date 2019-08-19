package Igotcp

//对连接器的封装
type IRequest interface {
	GetConnector() IConnector
	GetMsgData() []byte
	GetMsgId() uint32
}
