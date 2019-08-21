//连接器
//负责 为每个tcp链接生成一个连接器
package Gotcp

import (
	"errors"
	"gotcp/Conf"
	"gotcp/Igotcp"
	"io"
	"log"
	"net"
)

type Connector struct {
	conn     *net.TCPConn
	connID   int
	isClosed bool
	exitChan chan bool
	handle   Igotcp.IHandle
	msgChan  chan []byte
	instance Igotcp.IServer
}

//为每个连接生成一个连接器
func NewConnector(srv Igotcp.IServer, conn *net.TCPConn, connID int, r Igotcp.IHandle) (connector Igotcp.IConnector) {
	connector = &Connector{
		conn:     conn,
		connID:   connID,
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
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	// 销毁socket之前
	c.instance.CallOnConnStop(c)
	c.conn.Close()
	c.exitChan <- true
	c.instance.GetManager().Remove(c)
	close(c.exitChan)
	close(c.msgChan)
}

func (c *Connector) Read() {
	log.Println("[Info] read  goroutine is starting...")
	defer log.Println("[Info] read  goroutine is quit...")
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
		if _, err = io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			log.Println("[Error] io.ReadFull(c.GetTCPConnection(), headData) : ", err)
			break
		}

		if msg, err = mpkg.Unpack(headData); err != nil {
			log.Println("[Error] mpkg.Unpack(headData) : ", err)
			break
		}

		if msg.GetLen() > 0 {
			buf = make([]byte, msg.GetLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), buf); err != nil {
				log.Println("[Error]io.ReadFull(c.GetTCPConnection(), buf) : ", err)
				break
			}
			msg.SetData(buf)
		}

		req = &Request{
			connector: c,
			msg:       msg,
		}

		if Conf.SrvConf.WorkPoolSize > 0 {
			c.handle.SendMsgToTaskQueue(req)
		} else {
			go c.handle.DoMsgHandler(req)
		}

		log.Printf("[Info] server receive data : msgId : %d , msgLen : %d msgData : %s \n", msg.GetId(), msg.GetLen(), msg.GetData())
	}
}

func (c *Connector) Write() {
	log.Println("[Info] write goroutine is starting...")
	defer log.Println("[Info] write goroutine is quit...")
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.conn.Write(data); err != nil {
				log.Println("[Error] c.conn.Write(data) : ", err)
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

func (c *Connector) GetConnID() uint32 {
	return uint32(c.connID)
}

func (c *Connector) GetInstance() Igotcp.IServer {
	return c.instance
}

func (c *Connector) GetRemoteAddr() net.Addr {
	return nil
}
