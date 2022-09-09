package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"grpc/api"
	"grpc/etcdv3/register"
	"log"
	"net"
	"strings"
)

var endpoints string
var service string
var lease int64
var host string
var port string

func init() {
	flag.StringVar(&endpoints, "endpoints", "127.0.0.1:2379", "endpoints")
	flag.StringVar(&host, "host", "127.0.0.1", "host")
	flag.StringVar(&port, "port", "8080", "port")
	flag.StringVar(&service, "service", "hello", "service")
	flag.Int64Var(&lease, "lease", 5, "lease")
	flag.Parse()
}

type Greeter struct {
	api.UnimplementedGreeterServer
}

func (g *Greeter) SayHello(ctx context.Context, request *api.HelloRequest) (*api.HelloReply, error) {
	return &api.HelloReply{Message: "He" + port}, nil
}

func main() {
	_endpoints := strings.Split(endpoints, ",")
	addr := fmt.Sprintf("%s:%s", host, port)
	sr, err := register.New(_endpoints, service, addr, lease)
	if err != nil {
		log.Fatal(err)
	}
	defer sr.Revoke()

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()
	api.RegisterGreeterServer(grpcServer, &Greeter{})
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}
