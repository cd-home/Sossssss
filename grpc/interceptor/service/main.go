package main

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"grpc/api"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	/*
		    // 设置一个拦截器
			srv := grpc.NewServer([]grpc.ServerOption{
				grpc.UnaryInterceptor(HelloInterceptor),
			}...)

			// 设置多个
			srv := grpc.NewServer([]grpc.ServerOption{
				grpc.ChainUnaryInterceptor(
					HelloInterceptor,
					AccessLog,
				),
			}...)
	*/

	// 可以通过第三方包设置多个(内部递归调用的)

	srv := grpc.NewServer([]grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			HelloInterceptor,
			AccessLog,
		)),
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

// HelloInterceptor 一元拦截器, 普通模式
func HelloInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// do check op
	log.Println("你好")
	resp, err := handler(ctx, req)
	log.Println("再见")
	return resp, err
}

// AccessLog 多个拦截器的话, 是洋葱模型
func AccessLog(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	requestLog := "access request log: method: %s, begin_time: %d, request: %v"
	beginTime := time.Now().Local().Unix()
	log.Printf(requestLog, info.FullMethod, beginTime, req)

	resp, err := handler(ctx, req)

	responseLog := "access response log: method: %s, begin_time: %d, end_time: %d, response: %v"
	endTime := time.Now().Local().Unix()
	log.Printf(responseLog, info.FullMethod, beginTime, endTime, resp)
	return resp, err
}
