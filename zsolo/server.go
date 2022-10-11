package zsolo

import (
	"log"
	"net"
)

// Server
// Server 核心结构
type Server struct {
	Name    string
	NetWork string
	Host    string
	Port    string
	CM      *ConnManager // 链接控制
	DP      *DataPack    // 数据解包、打包
	HD      *Handler     // 消息处理器
}

func NewServer(name, network, host, port string) *Server {
	return &Server{
		Name:    name,
		NetWork: network,
		Host:    host,
		Port:    port,
		CM:      NewConnManager(),
		DP:      NewDataPack(),
		HD:      NewHandler(),
	}
}

func (s *Server) Run() {
	tcpAddr, err := net.ResolveTCPAddr(s.NetWork, net.JoinHostPort(s.Host, s.Port))
	if err != nil {
		log.Fatal(err)
	}
	listener, err := net.ListenTCP(s.NetWork, tcpAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Server Name [%s], Running Port [%s]", s.Name, s.Port)
	for {
		tcpConn, err := listener.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}
		// TODO how to gen connection id
		conn := NewConnection(s, tcpConn, "")
		go conn.Start()
	}
}

func (s *Server) Stop() {
	s.CM.ClearAll()
}
