package Conf

import (
	"encoding/json"
	"gotcp/Igotcp"
	"io/ioutil"
)

type Conf struct {
	Env     string
	Name    string
	Version string

	Host      string
	IPVersion string
	Port      int

	MaxConn         int
	MaxPkgSize      uint32
	WorkPoolSize    uint32
	MaxWorkPoolSize uint32
	TcpServer       Igotcp.IServer
}

var SrvConf *Conf

func init() {
	var (
		err  error
		data []byte
	)

	SrvConf = &Conf{
		Env:             "Test",
		TcpServer:       nil,
		MaxPkgSize:      512,
		WorkPoolSize:    10,
		MaxWorkPoolSize: 1024,
		Port:            8999,
		MaxConn:         1000,
		Version:         "v1.0",
		IPVersion:       "tcp4",
		Name:            "gotcp",
		Host:            "0.0.0.0",
	}

	func(){
		data, err = ioutil.ReadFile("conf/gotcp.json")
		err = json.Unmarshal(data, &SrvConf)
	}()

	if err != nil {
		panic(err)
	}
}
