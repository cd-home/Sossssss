package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc/api"
	"log"
)

// Token 定制认证方式
type Token struct {
	AppID     string
	AppSecret string
}

// GetRequestMetadata 获取当前请求认证所需的元数据（metadata）
func (t *Token) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"app_id": t.AppID, "app_secret": t.AppSecret}, nil
}

// RequireTransportSecurity 是否需要基于 TLS 认证进行安全传输
func (t *Token) RequireTransportSecurity() bool {
	return false
}

func main() {
	token := &Token{
		AppID:     "Foo",
		AppSecret: "Bar",
	}
	conn, _ := grpc.Dial("127.0.0.1:8081",
		// 禁用认证
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// 实际上是设置到metadata [header]
		grpc.WithPerRPCCredentials(token),
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
