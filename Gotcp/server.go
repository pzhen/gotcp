package Gotcp

import (
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"gotcp/Conf"
	"gotcp/Igotcp"
	"hash/crc32"
	"net"
	"os"
)

type Server struct {
	Address     string
	Network     string
	MaxConn     int
	Handle      Igotcp.IHandle
	Manager     Igotcp.IManager
	OnConnStop  func(conn Igotcp.IConnector)
	OnConnStart func(conn Igotcp.IConnector)
}

func InitServer() (srv Igotcp.IServer) {
	srv = &Server{
		Address:     Conf.SrvConf.Address,
		Network:     Conf.SrvConf.Network,
		MaxConn:     Conf.SrvConf.MaxConn,
		Handle:      NewMsgHandle(),
		Manager:     NewManager(),
		OnConnStop:  nil,
		OnConnStart: nil,
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
		s.OnConnStart(c)
	}
}

func (s *Server) CallOnConnStop(c Igotcp.IConnector) {
	if s.OnConnStop != nil {
		s.OnConnStop(c)
	}
}

func (s *Server) Start() {
	go func() {
		debugPrint("Worker pool is starting, worker num : %d", Conf.SrvConf.WorkPoolSize)
		s.Handle.StartWorkerPool()

		var (
			addr     *net.TCPAddr
			err      error
			listener *net.TCPListener
		)

		func() {
			addr, err = net.ResolveTCPAddr(s.Network, s.Address)
			listener, err = net.ListenTCP(s.Network, addr)
		}()

		if err != nil {
			debugPrintError("%+v", errors.WithStack(err))
			os.Exit(1)
		}

		for {
			var (
				conn      *net.TCPConn
				connector Igotcp.IConnector
			)

			if conn, err = listener.AcceptTCP(); err != nil {
				debugPrintError("%+v", errors.WithStack(err))
				continue
			}

			if s.Manager.Len() >= Conf.SrvConf.MaxConn {
				conn.Close()
				debugPrintError("%+v", errors.WithStack(errors.Errorf("Too Many Connections , MaxConn is %d", Conf.SrvConf.MaxConn)))
				continue
			}

			cid, _ := uuid.NewV4()
			connector = NewConnector(s, conn, cid.String(), s.Handle)
			debugPrint("UUID=%s, HashCode=%d", cid.String(), crc32.ChecksumIEEE([]byte(cid.String())))
			connector.Start()
		}
	}()
}

func (s *Server) Stop() {
	s.Manager.ClearConn()
	debugPrint("Server stoped HTTP on %s", s.Address)
}

func (s *Server) Run() {
	c := make(chan struct{})
	debugPrint("Listening and serving HTTP on %s", s.Address)
	s.Start()
	<-c
}
