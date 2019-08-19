package Gotcp

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func TestDataPack(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err: ", err)
		return
	}
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept error", err)
			}

			go func(conn net.Conn) {
				dp := NewMsgPack()
				for {
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head error")
						break
					}

					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("server unpacke err ", err)
						return
					}
					if msgHead.GetLen() > 0 {
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetLen())
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data err: ", err)
							return
						}
						fmt.Println("---> Recv MsgID: ", msg.Id, ", datalen = ", msg.Len, "data = ", string(msg.Data))
					}
				}

			}(conn)
		}
	}()

	/*
	   模拟客户端
	*/
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err: ", err)
		return
	}

	//创建一个封包对象 dp
	var dp = NewMsgPack() //模拟粘包过程，封装两个msg一同发送
	//封装第一个msg1包
	var msg1 = &Message{
		Id:   1,
		Len:  4,
		Data: []byte{'1', '2', '3', '4'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 error11111", err)
		return
	}
	//封装第二个msg2包
	msg2 := &Message{
		Id:   2,
		Len:  7,
		Data: []byte{'5', '6', '7', '8', '9', 'a', 'b'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg1 error", err)
		return
	}
	//将两个包粘在一起
	sendData1 = append(sendData1, sendData2...)
	conn.Write(sendData1)
	select {}
}
