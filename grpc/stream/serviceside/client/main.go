package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	api "grpc/api/serviceside"
	"io"
	"log"
)

var host string
var port string

func init() {
	flag.StringVar(&host, "host", "127.0.0.1", "host")
	flag.StringVar(&port, "port", "8080", "port")
	flag.Parse()
}

func main() {
	addr := fmt.Sprintf("%s:%s", host, port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return
	}
	client := api.NewServiceSideStreamClient(conn)
	streamValue, err := client.GetStreamValueFromService(context.Background(), &api.ClientRequest{
		Data: "Hello",
	})
	if err != nil {
		return
	}
	for {
		reply, err := streamValue.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		log.Println(reply.StreamValue)
	}
	streamValue.CloseSend()
}
