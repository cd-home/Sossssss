package connction

import (
	"sync"
)

type Connection struct {
	mu      sync.Mutex
	Clients map[uint32]*Client
}

func NewConnection() *Connection {
	return &Connection{
		Clients: make(map[uint32]*Client),
	}
}

func (c *Connection) Add() {

}

func (c *Connection) Delete() {

}

func (c *Connection) Clear() {

}
