package Gotcp

import (
	"github.com/pkg/errors"
	"gotcp/Conf"
	"gotcp/Igotcp"
	"hash/crc32"
	"io"
	"net"
	"sync"
)

type Connector struct {
	conn     *net.TCPConn
	uuid     string
	isClosed bool
	exitChan chan bool
	handle   Igotcp.IHandle
	msgChan  chan []byte
	instance Igotcp.IServer
	mutex     sync.Mutex
}

func NewConnector(srv Igotcp.IServer, conn *net.TCPConn, cid string, r Igotcp.IHandle) (connector Igotcp.IConnector) {
	connector = &Connector{
		conn:     conn,
		uuid:     cid,
		isClosed: false,
		handle:   r,
		exitChan: make(chan bool, 1),
		msgChan:  make(chan []byte),
		instance: srv,
	}

	connector.GetInstance().GetManager().Add(connector)
	return
}

func (c *Connector) Start() {
	go c.Read()
	go c.Write()
	c.instance.CallOnConnStart(c)
}

func (c *Connector) Stop() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if ! c.isClosed {
		c.isClosed = true
		c.instance.CallOnConnStop(c)
		c.conn.Close()
		c.exitChan <- true
		c.instance.GetManager().Remove(c)
		close(c.exitChan)
		close(c.msgChan)
	}
}

func (c *Connector) Read() {
	defer debugPrint("Goroutine read quit from UUID=%s", c.GetUUID())
	defer c.Stop()
	for {
		var (
			mpkg     Igotcp.IMsgPack
			msg      Igotcp.IMessage
			headData []byte
			buf      []byte
			err      error
			req      Igotcp.IRequest
		)

		mpkg = NewMsgPack()
		headData = make([]byte, mpkg.GetHeadLen())

		if _, err = io.ReadFull(c.GetTCPConnection(), headData);
		err != nil && err != io.EOF {
			c.Stop()
			debugPrintWarn(err.Error())
			break
		}

		if msg, err = mpkg.Unpack(headData); err != nil {
			c.Stop()
			debugPrintWarn(err.Error())
			break
		}

		buf = make([]byte, msg.GetLen())
		if _, err = io.ReadFull(c.GetTCPConnection(), buf);
		err != nil && err != io.EOF {
			c.Stop()
			debugPrintWarn(err.Error())
			break
		}
		msg.SetData(buf)

		req = &Request{
			connector: c,
			msg:       msg,
		}

		if Conf.SrvConf.WorkPoolSize > 0 {
			c.handle.SendMsgToTaskQueue(req)
		} else {
			go c.handle.DoMsgHandler(req)
		}

		//debugPrint("Receive : msgId : %d , msgLen : %d msgData : %s", msg.GetId(), msg.GetLen(), msg.GetData())
	}
}

func (c *Connector) Write() {
	defer debugPrint("Goroutine write quit from UUID=%s", c.GetUUID())
	defer c.Stop()
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.conn.Write(data); err != nil {
				debugPrintError("%+v", errors.WithStack(err))
				return
			}
		case <-c.exitChan:
			return
		}
	}
}

// 开发者发送消息，将消息封包后交给管道由 write() 发送给client
func (c *Connector) Send(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("connection is closed")
	}
	var (
		msgPkg    Igotcp.IMsgPack
		err       error
		binaryMsg []byte
	)

	msgPkg = NewMsgPack()
	if binaryMsg, err = msgPkg.Pack(NewMessage(msgId, data)); err != nil {
		return err
	}

	c.msgChan <- binaryMsg
	return nil
}

func (c *Connector) GetTCPConnection() *net.TCPConn {
	return c.conn
}

func (c *Connector) GetUUIDHashCode() uint32 {
	return crc32.ChecksumIEEE([]byte(c.uuid))
}

func (c *Connector) GetUUID() string {
	return c.uuid
}

func (c *Connector) GetInstance() Igotcp.IServer {
	return c.instance
}

func (c *Connector) GetRemoteAddr() net.Addr {
	return nil
}
