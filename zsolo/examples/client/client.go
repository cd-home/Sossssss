package main

import (
	"io"
	"log"
	"net"
	"time"
	"zinx"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Println(err)
	}
	for {
		dp := zsolo.NewDataPack()
		msg := zsolo.NewMessage(1, []byte("Hello"))
		data, err := dp.Pack(msg)
		if err != nil {
			return
		}

		// write
		conn.Write(data)

		// Read
		dp2 := zsolo.NewDataPack()
		head := make([]byte, dp2.HeadLen())
		_, err = io.ReadFull(conn, head)
		if err != nil {
			log.Println(err)
			break
		}
		msgHead, err := dp2.UnPack(head)
		if err != nil {
			log.Println(err)
			break
		}
		if msgHead.MsgLen() > 0 {
			msgHead.Data = make([]byte, msgHead.MsgLen())
			_, err = io.ReadFull(conn, msgHead.Data)
			if err != nil {
				log.Println(err)
				break
			}
			log.Println(string(msgHead.Data))
		}
		time.Sleep(2 * time.Second)
	}
}
