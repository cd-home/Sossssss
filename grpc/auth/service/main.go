package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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

	// authorization bearer|basic
	srv := grpc.NewServer([]grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			Check,
		),
	}...)

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

// Check 一元拦截器, 普通模式
func Check(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// 从 ctx 中就可以获取到 metadata
	md, _ := metadata.FromIncomingContext(ctx)
	if _, ok := md["app_id"]; !ok {
		return nil, status.Error(codes.Unauthenticated, "认证未通过")
	}
	fmt.Println(md)
	return handler(ctx, req)
}
