package zinx

import (
	"context"
	"log"
	"net"
	"sync"
)

type Connection struct {
	ctx       context.Context
	cancel    context.CancelFunc
	mu        sync.Mutex
	TcpServer *Server
	ID        string
	Conn      net.Conn
	msgCh     chan []byte
	props     map[string]string
	state     bool
}

func NewConnection(s *Server, conn net.Conn, id string) *Connection {
	return &Connection{
		TcpServer: s,
		ID:        id,
		Conn:      conn,
		msgCh:     make(chan []byte, 4096),
		state:     true,
	}
}

func (c *Connection) Send() {
	log.Printf("[%s] Connection Send G is Running...", c.ID)
	for {
		select {
		case data, ok := <-c.msgCh:
			if !ok {
				break
			}
			_, err := c.Conn.Write(data)
			if err != nil {
				log.Println("send data error ", err)
				return
			}
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *Connection) Recv() {
	log.Printf("[%s] Connection Recv G is Running...", c.ID)
	
}
