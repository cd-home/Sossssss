package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	api "grpc/api/serviceside"
	"log"
	"net"
	"strconv"
	"time"
)

var host string
var port string

func init() {
	flag.StringVar(&host, "host", "127.0.0.1", "host")
	flag.StringVar(&port, "port", "8080", "port")
	flag.Parse()
}

type ServiceSideStream struct {
	api.UnimplementedServiceSideStreamServer
}

func (s *ServiceSideStream) GetStreamValueFromService(req *api.ClientRequest, srv api.ServiceSideStream_GetStreamValueFromServiceServer) error {
	for i := 0; i < 10; i++ {
		err := srv.Send(&api.ServiceReply{StreamValue: strconv.Itoa(i) + req.Data})
		time.Sleep(time.Second)
		if err != nil {
			return err
		}
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
	api.RegisterServiceSideStreamServer(srv, &ServiceSideStream{})
	err = srv.Serve(listener)
	if err != nil {
		log.Fatal(err)
		return
	}
}
