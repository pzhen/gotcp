package main

import (
	"fmt"
	"gotcp/Gotcp"
	"io"
	"net"
	"time"
)

func main() {
	var (
		conn net.Conn
		err  error
		b    []byte
	)
	if conn, err = net.Dial("tcp", "127.0.0.1:9888"); err != nil {
		return
	}

	for {
		msgPkg := Gotcp.NewMsgPack()
		if b, err = msgPkg.Pack(Gotcp.NewMessage(1, []byte("nihao1111..."))); err != nil {
			fmt.Println("pack error")
			return
		}

		if _, err = conn.Write(b); err != nil {
			fmt.Println("write error", err)
			return
		}

		headData := make([]byte, msgPkg.GetHeadLen())
		_, err := io.ReadFull(conn, headData)
		if err != nil {
			fmt.Println("read head error")
			break
		}

		msgHead, err := msgPkg.Unpack(headData)
		if err != nil {
			fmt.Println("server unpacke err ", err)
			return
		}
		if msgHead.GetLen() > 0 {
			msg := msgHead.(*Gotcp.Message)
			msg.Data = make([]byte, msg.GetLen())
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("server unpack data err: ", err)
				return
			}
			fmt.Println("---> Recv MsgID: ", msg.Id, ", datalen = ", msg.Len, "data = ", string(msg.Data))
		}

		time.Sleep(2 * time.Second)
	}

}
