package main

import (
	"fmt"
	"gotcp/Gotcp"
	"gotcp/Igotcp"
	"os"
	"runtime/trace"
)

type MyRouter1 struct {
	Gotcp.BaseRouter
}

func (m *MyRouter1) Handle(r Igotcp.IRequest) {
	if err := r.GetConnector().Send(1, []byte("11111 handle ..........")); err != nil {
		fmt.Println(err)
	}
}


type MyRouter2 struct {
	Gotcp.BaseRouter
}

func (m *MyRouter2) Handle(r Igotcp.IRequest) {
	if err := r.GetConnector().Send(1, []byte("22222 handle .........")); err != nil {
		fmt.Println(err)
	}
}

func doBefore(connector Igotcp.IConnector)  {
	fmt.Println("=============>doBefore")
	connector.Send(1, []byte("doBefore"))
}

func doAfter(connector Igotcp.IConnector)  {
	fmt.Println("=============>doAfter")
	connector.Send(2, []byte("doAfter"))
}

func main() {

	go func() {
		f, err := os.Create("trace.out")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		err = trace.Start(f)
		if err != nil {
			panic(err)
		}
		defer trace.Stop()
	}()


	srv := Gotcp.InitServer()
	srv.AddRouter(1, &MyRouter1{})
	srv.AddRouter(2, &MyRouter2{})

	srv.SetOnConnStart(doBefore)
	srv.SetOnConnStop(doAfter)
	srv.Run()
}
