package zine

import (
	"errors"
	"sync"
)

type ConnManager struct {
	mu          sync.Mutex
	Connections map[string]*Connection
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		Connections: make(map[string]*Connection),
	}
}

func (cm *ConnManager) Add(conn *Connection) {
	cm.mu.Lock()
	cm.Connections[conn.ID] = conn
	cm.mu.Unlock()
}

func (cm *ConnManager) Get(connID string) (*Connection, error) {
	cm.mu.Lock()
	if conn, ok := cm.Connections[connID]; ok {
		return conn, nil
	}
	cm.mu.Unlock()
	return nil, errors.New("connection Not Found")
}

func (cm *ConnManager) Remove(conn *Connection) {
	cm.mu.Lock()
	delete(cm.Connections, conn.ID)
	conn.Close()
	cm.mu.Unlock()
}

func (cm *ConnManager) ClearAll() {
	cm.mu.Lock()
	for connID, conn := range cm.Connections {
		// 关闭链接
		conn.Close()
		// 移除链接
		delete(cm.Connections, connID)
	}
	cm.mu.Unlock()
}
