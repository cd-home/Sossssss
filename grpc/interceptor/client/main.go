package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc/api"
	"log"
	"time"
)

func main() {
	conn, _ := grpc.Dial("127.0.0.1:8081",
		// 禁用认证
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// 设置一个
		grpc.WithUnaryInterceptor(HelloClient),
		// 设置多个
		grpc.WithChainUnaryInterceptor(
			HelloClient,
			UnaryContextTimeout,
		),
	)

	defer conn.Close()

	// 获取gRPC客户端
	// client stub
	client := api.NewGreeterClient(conn)
	// blocking/synchronous mode 阻塞同步模式
	reply, err := client.SayHello(context.Background(), &api.HelloRequest{Name: "Nike"})
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(reply.Message)
}

func defaultContextTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	var cancel context.CancelFunc
	if _, ok := ctx.Deadline(); !ok {
		defaultTimeout := 2 * time.Second
		ctx, cancel = context.WithTimeout(ctx, defaultTimeout)
	}

	return ctx, cancel
}

// UnaryContextTimeout 客户端拦截器
// grpc.UnaryClientInterceptor 客户端拦截器方法原型
// type UnaryClientInterceptor func(ctx context.Context, method string, req, reply interface{}, cc *ClientConn, invoker UnaryInvoker, opts ...CallOption) error
// type StreamClientInterceptor func(ctx context.Context, desc *StreamDesc, cc *ClientConn, method string, streamer Streamer, opts ...CallOption) (ClientStream, error)
func UnaryContextTimeout(ctx context.Context, method string, req, resp interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	ctx_, cancel := defaultContextTimeout(ctx)
	if cancel != nil {
		defer cancel()
	}
	return invoker(ctx_, method, req, resp, cc, opts...)
}

func HelloClient(ctx context.Context, method string, req, resp interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	log.Println("Client Start1")
	err := invoker(ctx, method, req, resp, cc, opts...)
	log.Println("Client End1")
	return err
}

/*func UnaryContextTimeout() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, resp interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx, cancel := defaultContextTimeout(ctx)
		if cancel != nil {
			defer cancel()
		}
		return invoker(ctx, method, req, resp, cc, opts...)
	}
}*/
