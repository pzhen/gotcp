package Igotcp

type IMsgPack interface {
	GetHeadLen() int
	Pack(msg IMessage) ([]byte, error)
	Unpack([]byte) (IMessage, error)
}
