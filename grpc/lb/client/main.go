package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"grpc/api"
	"grpc/etcdv3/discovery"
	"log"
	"strings"
)

var endpoints string
var service string

func init() {
	flag.StringVar(&endpoints, "endpoints", "127.0.0.1:2379", "endpoints")
	flag.StringVar(&service, "service", "hello", "service")
	flag.Parse()
}

func main() {
	_endpoints := strings.Split(endpoints, ",")
	r := discovery.New(_endpoints)
	resolver.Register(r)
	conn, err := grpc.Dial(
		// URI schemes
		// dns:[//authority/]host[:port] -- DNS (default)
		r.Scheme()+"://authority/"+service,
		// grpc.WithBalancerName() 废弃
		grpc.WithDefaultServiceConfig(
			fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`,
				roundrobin.Name)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return
	}
	grpcClient := api.NewGreeterClient(conn)
	response, err := grpcClient.SayHello(context.Background(), &api.HelloRequest{
		Name: "yao",
	})
	if err != nil {
		log.Println(err)
	}
	fmt.Println(response.Message)
}
