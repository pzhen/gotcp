package main

import (
	"fmt"
	"gotcp/Gotcp"
	"gotcp/Igotcp"
)

type MyRouter struct {
	Gotcp.BaseRouter
}

func (m *MyRouter) Handle(r Igotcp.IRequest) {
	if err := r.GetConnector().Send(1, []byte("123456....")); err != nil {
		fmt.Println(err)
	}
}

func main() {
	srv := Gotcp.InitServer()
	srv.AddRouter(&MyRouter{})
	srv.Run()
}
