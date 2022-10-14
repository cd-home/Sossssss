package main

import (
	"log"
	"zinx"
)

type PingPong struct {
}

func (p *PingPong) Handle(request *zine.Request) {
	err := request.Conn.SendMessage(1001, []byte("pong from server: "+string(request.Msg.Data)))
	if err != nil {
		log.Println(err)
	}
}

func main() {
	s := zine.NewServer("s", "tcp", "localhost", "8080")
	s.HD.AddRouter(0, &PingPong{})
	s.Run()
}
