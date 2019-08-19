package Gotcp

import (
	"fmt"
	"gotcp/Conf"
	"gotcp/Helper"
	"gotcp/Igotcp"
	"log"
	"net"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
	MaxConn   int
	Handle    Igotcp.IHandle
	Manager   Igotcp.IManager

	OnConnStart func(conn Igotcp.IConnector)
	OnConnStop  func(conn Igotcp.IConnector)
}

func init() {
	Helper.PrintLog()
}

// 初始化gotcp服务
func InitServer() (srv Igotcp.IServer) {
	srv = &Server{
		Name:      Conf.G_Conf.Name,
		IP:        Conf.G_Conf.Host,
		Port:      Conf.G_Conf.Port,
		MaxConn:   Conf.G_Conf.MaxConn,
		IPVersion: Conf.G_Conf.IPVersion,
		Handle:    NewMsgHandle(),
		Manager:   NewManager(),
	}
	return
}

func (s *Server) GetManager() Igotcp.IManager {
	return s.Manager
}

func (s *Server) AddRouter(msgId uint32, router Igotcp.IRouter) {
	s.Handle.AddRouter(msgId, router)
}

func (s *Server) SetOnConnStart(f func(Igotcp.IConnector)) {
	s.OnConnStart = f
}

func (s *Server) SetOnConnStop(f func(Igotcp.IConnector)) {
	s.OnConnStop = f
}

func (s *Server) CallOnConnStart(c Igotcp.IConnector) {
	if s.OnConnStart != nil {
		log.Println("[Info] call on start connect")
		s.OnConnStart(c)
	}

}

func (s *Server) CallOnConnStop(c Igotcp.IConnector) {
	if s.OnConnStop != nil {
		log.Println("[Info] call on stop connect")
		s.OnConnStop(c)
	}
}

//启动TCP服务
func (s *Server) Start() {
	go func() {
		// 开启工作池
		s.Handle.StartWorkerPool()

		var (
			addr     *net.TCPAddr
			err      error
			listener *net.TCPListener
			cid      int
		)

		if addr, err = net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port)); err != nil {
			log.Fatalln("[Error] net.ResolveTCPAddr() : ", err)
			return
		}

		if listener, err = net.ListenTCP(s.IPVersion, addr); err != nil {
			log.Fatalln("[Error] net.ListenTCP(s.IPVersion, addr) : ", err)
		}

		for {
			var (
				conn      *net.TCPConn
				connector Igotcp.IConnector
			)

			if conn, err = listener.AcceptTCP(); err != nil {
				log.Println("[Warning] listener.AcceptTCP() : ", err)
				continue
			}

			if s.Manager.Len() >= Conf.G_Conf.MaxConn {
				//TODO 错误包
				conn.Close()
				log.Println("[Error] too Many Connections MaxConn : ", Conf.G_Conf.MaxConn)
				continue
			}

			// TODO 暂时链接ID用自增来处理
			cid++
			//获取连接器
			connector = NewConnector(s, conn, cid, s.Handle)
			log.Printf("[Info] welcome clientID : %d connect ...", cid)

			//启动连接器
			connector.Start()
		}
	}()
}

func (s *Server) Stop() {
	s.Manager.ClearConn()
	log.Printf("[Stop] gotcp %s is stoped at IP :%s, Port %d ......\n", Conf.G_Conf.Version, Conf.G_Conf.Host, Conf.G_Conf.Port)
}

func (s *Server) Run() {
	c := make(chan struct{})
	s.Start()
	log.Printf("[Start] gotcp %s is running at IP :%s, Port %d ......\n", Conf.G_Conf.Version, Conf.G_Conf.Host, Conf.G_Conf.Port)
	<-c
}
