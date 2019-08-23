![](./logo.png)

`gotcp` 是采用golang,以面向接口方式开发的,轻量级的,TCP框架.


## Features

- Golang开发-并发处理请求
- 面向接口-灵活扩展
- Request封装-消息数据与请求的包装
- Message封装-粘包问题处理
- Router封装-多路由模式,灵活注册路由来处理不同消息
- 读写分离
- 消息队列及多任务处理(work pool)
- 链接管理-在创建后与销毁之前灵活处理

## Getting Started

### Installing

```sh
$ go get github.com/pzhen/gotcp
```

### Usage

项目中配置配置文件

conf/gotcp.json
```json
{
  "Name": "MyTcpServer",
  "Address": "127.0.0.1:9888",
  "MaxConn": 1,
  "MaxPkgSize": 512,
  "Network":"tcp4",
  "WorkPoolSize":10
 }
```

server.go
```
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
```

## License

