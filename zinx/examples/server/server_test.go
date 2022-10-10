package main

import (
	"net"
	"testing"
)

func TestServer(t *testing.T) {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		t.Error(err)
	}
	_, err = conn.Write([]byte("Hello Server"))
	if err != nil {
		t.Error(err)
	}
	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(buf[:n]))
}
