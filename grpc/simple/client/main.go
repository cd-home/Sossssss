package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc/api"
	"log"
)

func main() {
	conn, _ := grpc.Dial("127.0.0.1:8081",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	defer conn.Close()

	// 获取gRPC客户端
	client := api.NewGreeterClient(conn)
	reply, err := client.SayHello(context.Background(), &api.HelloRequest{Name: "Nike"})
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(reply.Message)
}
