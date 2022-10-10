package zinx

import (
	"log"
	"net"
)

type Server struct {
	Name    string
	NetWork string
	Host    string
	Port    string
}

func NewServer(name, network, host, port string) *Server {
	return &Server{
		Name:    name,
		NetWork: network,
		Host:    host,
		Port:    port,
	}
}

func (s *Server) Start() {
	tcpAddr, err := net.ResolveTCPAddr(s.NetWork, net.JoinHostPort(s.Host, s.Port))
	if err != nil {
		log.Fatal(err)
	}
	listener, err := net.ListenTCP(s.NetWork, tcpAddr)
	log.Printf("Server Name [%s], Running Port [%s]", s.Name, s.Port)
	if err != nil {
		log.Fatal(err)
	}
	for {
		tcpConn, err := listener.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}
		go func() {
			buf := make([]byte, 512)
			n, err := tcpConn.Read(buf)
			if err != nil {
				log.Println(err)
				return
			}
			_, err = tcpConn.Write(buf[:n])
			if err != nil {
				log.Println(err)
				return
			}
			return
		}()
	}
}
