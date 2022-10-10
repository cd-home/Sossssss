package main

import "zinx"

func main() {
	s := zinx.NewServer("s", "tcp", "localhost", "8080")
	s.Start()
}
