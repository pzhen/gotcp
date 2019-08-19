package main

import (
	"gotcp/Gotcp"
	"gotcp/Igotcp"
	"log"
)

type MyRouter struct {
	Gotcp.BaseRouter
}

func (m *MyRouter) BeforeHandle(r Igotcp.IRequest) {
	_, e := r.GetConnector().GetTCPConnection().Write([]byte("before ...\n"))
	if e != nil {
		log.Println(e)
	}
}

func (m *MyRouter) Handle(r Igotcp.IRequest) {
	_, e := r.GetConnector().GetTCPConnection().Write([]byte("handle ...\n"))
	if e != nil {
		log.Println(e)
	}
}

func (m *MyRouter) AfterHandle(r Igotcp.IRequest) {
	_, e := r.GetConnector().GetTCPConnection().Write([]byte("after ...\n"))
	if e != nil {
		log.Println(e)
	}
}

func main() {
	srv := Gotcp.InitServer()
	srv.AddRouter(&MyRouter{})
	srv.Run()
}
