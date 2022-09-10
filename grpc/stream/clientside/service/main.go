package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	api "grpc/api/clientside"
	"io"
	"log"
	"net"
)

var host string
var port string

func init() {
	flag.StringVar(&host, "host", "127.0.0.1", "host")
	flag.StringVar(&port, "port", "8080", "port")
	flag.Parse()
}

type ServiceSide struct {
	api.UnimplementedClientSideStreamServer
}

func (s *ServiceSide) SendStreamToService(req api.ClientSideStream_SendStreamToServiceServer) error {
	for {
		res, err := req.Recv()
		if err == io.EOF {
			req.SendAndClose(&api.ServiceReply{
				Code:  "0",
				Value: "OK",
			})
		}
		if err != nil {
			break
		}
		log.Println(res.StreamValue)
	}
	return nil
}

func main() {
	addr := fmt.Sprintf("%s:%s", host, port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	srv := grpc.NewServer()
	api.RegisterClientSideStreamServer(srv, &ServiceSide{})
	srv.Serve(listener)
}
