package main

import (
	"fmt"
	"gotcp/Gotcp"
	"gotcp/Igotcp"
	"io"
	"net"
	"strconv"
	"time"
)

func main() {
	c:= make(chan struct{})

	for i:=0; i< 100;i++ {
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

				func(){
					b, err = msgPkg.Pack(Gotcp.NewMessage(1, []byte("我是client " + strconv.Itoa(i))))
					_, err = conn.Write(b)
				}()

				if err != nil {
					fmt.Println(err)
					return
				}

				func(){
					headData := make([]byte, msgPkg.GetHeadLen())
					_, err = io.ReadFull(conn, headData)
					msgHead, err = msgPkg.Unpack(headData)

					if msgHead.GetLen() > 0 {
						msg := msgHead.(*Gotcp.Message)
						msg.Data = make([]byte, msg.GetLen())
						_, err = io.ReadFull(conn, msg.Data)
						fmt.Println("---> Recv MsgID: ", msg.Id, ", datalen = ", msg.Len, "data = ", string(msg.Data))
					}
				}()

				if err != nil {
					fmt.Println(err)
					break
				}
				time.Sleep(2 * time.Second)
			}
		}(i)
	}

	<-c




}
