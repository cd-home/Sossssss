package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	api "grpc/api/clientside"
	"log"
	"strconv"
	"time"
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
	client := api.NewClientSideStreamClient(conn)
	stream, err := client.SendStreamToService(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		stream.Send(&api.ClientRequest{StreamValue: "From Client: " + strconv.Itoa(i)})
		log.Println("send: ", strconv.Itoa(i))
		time.Sleep(time.Second)
	}
	log.Println("send ok")
	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply)
}
