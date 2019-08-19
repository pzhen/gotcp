package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	var (
		conn net.Conn
		err  error
		buf  []byte
		cnt  int
	)
	if conn, err = net.Dial("tcp", "127.0.0.1:9888"); err != nil {
		return
	}

	for {
		if _, err = conn.Write([]byte("hello gotcp ...")); err != nil {
			return
		}

		buf = make([]byte, 512)
		if cnt, err = conn.Read(buf); err != nil {
			return
		}

		fmt.Println(string(buf[:cnt]))

		time.Sleep(2 * time.Second)
	}

}
