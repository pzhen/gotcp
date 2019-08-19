package Conf

import (
	"encoding/json"
	"gotcp/Igotcp"
	"io/ioutil"
	"log"
)

type GConf struct {
	Name            string
	Host            string
	Version         string
	IPVersion       string
	Port            int
	MaxConn         int
	MaxPkgSize      uint32
	WorkPoolSize    uint32
	MaxWorkPoolSize uint32
	TcpServer       Igotcp.IServer
}

var G_Conf *GConf

func init() {
	var (
		err  error
		data []byte
	)

	G_Conf = &GConf{
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

	if data, err = ioutil.ReadFile("conf/gotcp.json"); err != nil {
		log.Fatalln("[Panic] ioutil.ReadFile(conf/gotcp.json) : ", err)
	}

	if err = json.Unmarshal(data, &G_Conf); err != nil {
		log.Fatalln("[Panic] json.Unmarshal(data, &G_Conf) : ", err)
	}
}
