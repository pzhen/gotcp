package main

import (
	"fmt"
	_ "github.com/mkevac/debugcharts"
	"gotcp/Gotcp"
	"gotcp/Igotcp"
	_ "net/http/pprof"
)

type MyRouter1 struct {
	Gotcp.BaseRouter
}

type MyRouter2 struct {
	Gotcp.BaseRouter
}

//业务处理前
func (m *MyRouter1) BeforeHandle(r Igotcp.IRequest) {
	if err := r.GetConnector().Send(1, []byte("MyRouter1======>BeforeHandle")); err != nil {
		fmt.Println(err)
	}
}

//业务处理
func (m *MyRouter1) Handle(r Igotcp.IRequest) {
	if err := r.GetConnector().Send(1, []byte("MyRouter1======>Handle")); err != nil {
		fmt.Println(err)
	}
}

//业务处理后
func (m *MyRouter1) AfterHandle(r Igotcp.IRequest) {
	if err := r.GetConnector().Send(1, []byte("MyRouter1======>AfterHandle")); err != nil {
		fmt.Println(err)
	}
}

//直接处理业务
func (m *MyRouter2) Handle(r Igotcp.IRequest) {
	if err := r.GetConnector().Send(1, []byte("MyRouter2======>handle")); err != nil {
		fmt.Println(err)
	}
}

//连接器创建后做一些处理
func onConnect(connector Igotcp.IConnector) {
	connector.Send(1, []byte("doBefore"))
}

//连接器销毁前做一些处理
func offConnect(connector Igotcp.IConnector) {
	connector.Send(2, []byte("doAfter"))
}

func main() {
	//init server
	srv := Gotcp.InitServer()

	//注册路由
	srv.AddRouter(1, &MyRouter1{})
	srv.AddRouter(2, &MyRouter2{})

	//注册Hook
	srv.SetOnConnStart(onConnect)
	srv.SetOnConnStop(offConnect)

	//启动服务
	srv.Run()
}
