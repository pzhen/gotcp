package Igotcp

import "net"

type IConnector interface {
	Start()
	Stop()
	Read()
	Write()
	GetTCPConnection() *net.TCPConn
	GetUUIDHashCode() uint32
	GetUUID() string
	GetRemoteAddr() net.Addr
	Send(msgId uint32, data []byte) error
	GetInstance() IServer
}
