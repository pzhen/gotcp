//对消息message进行包装
//防止粘包,采用TLV格式进行封包解包
package Gotcp

import (
	"bytes"
	"encoding/binary"
	"errors"
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
	dataBuff := bytes.NewBuffer([]byte{})
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetLen()); err != nil {
		return nil, err
	}

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetId()); err != nil {
		return nil, err
	}

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

func (p *msgPkg) Unpack(binaryData []byte) (Igotcp.IMessage, error) {
	dataBuff := bytes.NewReader(binaryData)
	var (
		id  uint32
		len uint32
	)
	if err := binary.Read(dataBuff, binary.LittleEndian, &len); err != nil {
		return nil, err
	}
	if err := binary.Read(dataBuff, binary.LittleEndian, &id); err != nil {
		return nil, err
	}

	msg := &Message{
		Len: len,
		Id:  id,
	}

	if Conf.SrvConf.MaxPkgSize > 0 && msg.Len > Conf.SrvConf.MaxPkgSize {
		return nil, errors.New("too Large msg data recv ")
	}

	return msg, nil
}
