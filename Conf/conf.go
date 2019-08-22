package Conf

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
)

type Conf struct {
	Env             string
	Name            string
	Version         string
	Address         string
	Network         string
	MaxConn         int
	MaxPkgSize      uint32
	WorkPoolSize    uint32
	MaxWorkPoolSize uint32
}

var SrvConf *Conf

func init() {
	var (
		err  error
		data []byte
	)

	SrvConf = &Conf{
		"Test",
		"Gotcp",
		"v1.0",
		"0.0.0.0:8999",
		"tcp4",
		5000000,
		512,
		10,
		1024,
	}

	func() {
		data, err = ioutil.ReadFile("conf/gotcp.json")
		err = json.Unmarshal(data, &SrvConf)
	}()

	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", errors.WithStack(err))
		os.Exit(1)
	}
}
