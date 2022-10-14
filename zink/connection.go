package zine

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"sync"
	"time"
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
	c := &Connection{
		TcpServer: s,
		ID:        id,
		Conn:      conn,
		msgCh:     make(chan []byte),
		state:     true,
	}
	return c
}

func (c *Connection) Start() {
	c.ctx, c.cancel = context.WithCancel(context.Background())
	go c.Send()
	go c.Recv()
	select {
	case <-c.ctx.Done():
		_ = c.Conn.Close()
		c.state = false
		close(c.msgCh)
		return
	}
}

func (c *Connection) Send() {
	log.Printf("[%s] Connection Send G is Running.", c.ID)
	defer log.Printf("[%s] Connection Send G was Stop.", c.ID)
	for {
		select {
		case data, ok := <-c.msgCh:
			if !ok {
				return
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
	log.Printf("[%s] Connection Recv G is Running.", c.ID)
	defer log.Printf("[%s] Connection Recv G was Stop.", c.ID)
	defer c.Close()
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			head := make([]byte, c.TcpServer.DP.HeadLen())
			_, err := io.ReadFull(c.Conn, head)
			// 客户端关闭
			if err == io.EOF {
				log.Printf("[%s] Client Connection Close", c.ID)
				c.Close()
				c.TcpServer.CM.Remove(c)
				return
			}
			if err != nil {
				log.Printf("[%s] Connection Read Data Error", c.ID)
				return
			}
			msg, err := c.TcpServer.DP.UnPack(head)
			if err != nil {
				log.Printf("[%s] Connection UpPack Head Data Error", c.ID)
			}
			var data []byte
			if msg.MsgLen() > 0 {
				data = make([]byte, msg.MsgLen())
				if _, err := io.ReadFull(c.Conn, data); err != nil {
					log.Printf("[%s] Connection Read Body Data Error", c.ID)
					return
				}
			}
			msg.Data = data
			request := &Request{
				Conn: c,
				Msg:  msg,
			}
			// Handle Message
			// TODO work pool
			go c.TcpServer.HD.Handle(request)
		}
	}
}

func (c *Connection) Close() {
	c.cancel()
}

func (c *Connection) SendMessage(msgID uint32, data []byte) error {
	if !c.state {
		return errors.New("connection closed")
	}
	msg, err := c.TcpServer.DP.Pack(NewMessage(msgID, data))
	if err != nil {
		return err
	}
	timeout := time.NewTimer(time.Second)
	defer timeout.Stop()
	select {
	case <-timeout.C:
		return errors.New("send to client timeout")
	case c.msgCh <- msg:
		return nil
	}
}
