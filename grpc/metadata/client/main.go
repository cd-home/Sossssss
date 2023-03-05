package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"grpc/api"
	"log"
)

func main() {
	conn, _ := grpc.Dial("127.0.0.1:8081",
		// 禁用认证
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	defer conn.Close()

	// 获取gRPC客户端
	// client stub
	client := api.NewGreeterClient(conn)
	// blocking/synchronous mode 阻塞同步模式

	// 设置metadata
	md := metadata.New(map[string]string{"go": "1.18", "call": "rpc"})
	// 单个添加
	//ctx := metadata.AppendToOutgoingContext(context.Background(), "app", "001")
	// 添加 md
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// metadata 是存储在头部的
	reply, err := client.SayHello(ctx, &api.HelloRequest{Name: "Nike"})
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(reply.Message)
}
