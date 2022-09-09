package main

import (
	"context"
	"google.golang.org/grpc"
	"grpc/api"
	"net"
	"os"
	"os/signal"
	"syscall"
)

// Greeter
// protoc --go_out=. --go-grpc_out=. ./proto/*.proto
type Greeter struct {
	api.UnimplementedGreeterServer
}

// SayHello
// implemented service methods
func (g *Greeter) SayHello(ctx context.Context, request *api.HelloRequest) (*api.HelloReply, error) {
	return &api.HelloReply{Message: "He"}, nil
}

func main() {
	listener, _ := net.Listen("tcp", "127.0.0.1:8081")

	// 创建服务
	//var opts []grpc.ServerOption
	//srv := grpc.NewServer(opts...)
	srv := grpc.NewServer()
	// Greeter注册到gRPC服务
	api.RegisterGreeterServer(srv, &Greeter{})

	// 启动服务
	go func() {
		_ = srv.Serve(listener)
	}()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	srv.GracefulStop()
}
