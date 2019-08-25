package main

import (
	"fmt"
	"gotcp/Gotcp"
	"gotcp/Igotcp"
	"io"
	"net"
	"strconv"
	"syscall"
	"time"
)

func main() {

	// Increase resources limitations
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}


	c:= make(chan struct{})

	for i:=0; i< 500;i++ {
		go func(i int) {
			var (
				conn net.Conn
				err  error
				b    []byte
				msgHead Igotcp.IMessage
			)
			if conn, err = net.Dial("tcp4", "127.0.0.1:9888"); err != nil {
				fmt.Println(err)
				return
			}

			for {
				msgPkg := Gotcp.NewMsgPack()
				if b, err = msgPkg.Pack(Gotcp.NewMessage(1, []byte("我是client " + strconv.Itoa(i))));err != nil{
					fmt.Println(err)
					return
				}

				if _, err = conn.Write(b); err != nil {
					fmt.Println(err)
					return
				}

				headData := make([]byte, msgPkg.GetHeadLen())
				_, err = io.ReadFull(conn, headData)

				if err!=nil {
					fmt.Println(err)
					return
				}

				msgHead, err = msgPkg.Unpack(headData)

				if msgHead.GetLen() > 0 {
					msg := msgHead.(*Gotcp.Message)
					msg.Data = make([]byte, msg.GetLen())
					_, err = io.ReadFull(conn, msg.Data)
					fmt.Println("---> Recv MsgID: ", msg.Id, ", datalen = ", msg.Len, "data = ", string(msg.Data))
				}

				if err != nil {
					fmt.Println(err)
					return
				}
				time.Sleep(2 * time.Second)
			}
		}(i)

		time.Sleep(time.Millisecond*200)
	}
	<-c
}
