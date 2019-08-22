package Gotcp

import (
	"bytes"
	"encoding/binary"
	"github.com/pkg/errors"
	"gotcp/Conf"
	"gotcp/Igotcp"
)

type msgPkg struct{}

func NewMsgPack() *msgPkg {
	return &msgPkg{}
}

func (p *msgPkg) GetHeadLen() int {
	return 8
}

func (p *msgPkg) Pack(msg Igotcp.IMessage) ([]byte, error) {
	var (
		err      error
		dataBuff *bytes.Buffer
	)

	func() {
		dataBuff = bytes.NewBuffer([]byte{})
		err = binary.Write(dataBuff, binary.LittleEndian, msg.GetLen())
		err = binary.Write(dataBuff, binary.LittleEndian, msg.GetId())
		err = binary.Write(dataBuff, binary.LittleEndian, msg.GetData())
	}()

	if err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

func (p *msgPkg) Unpack(binaryData []byte) (Igotcp.IMessage, error) {
	var (
		id       uint32
		len      uint32
		err      error
		dataBuff *bytes.Reader
	)

	func() {
		dataBuff = bytes.NewReader(binaryData)
		err = binary.Read(dataBuff, binary.LittleEndian, &len)
		err = binary.Read(dataBuff, binary.LittleEndian, &id)
	}()

	if err != nil {
		return nil, err
	}

	msg := &Message{
		Len: len,
		Id:  id,
	}

	if Conf.SrvConf.MaxPkgSize > 0 && msg.Len > Conf.SrvConf.MaxPkgSize {
		return nil, errors.New("Msgpkg too large")
	}
	return msg, nil
}
