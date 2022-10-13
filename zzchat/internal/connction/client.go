package connction

import (
	"nhooyr.io/websocket"
)

type Client struct {
	uid   uint32
	conn  *websocket.Conn
	msgCh chan []byte
}

func NewClient(uid uint32, conn *websocket.Conn) *Client {
	return &Client{
		conn:  nil,
		msgCh: make(chan []byte),
	}
}

func (c *Client) Read() {

}

func (c *Client) Write() {

}

func (c *Client) Close() {

}
