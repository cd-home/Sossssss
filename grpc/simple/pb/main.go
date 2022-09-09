package main

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"grpc/api"
)

func main() {
	r := &api.HelloRequest{Name: "yao"}
	bytes, err := proto.Marshal(r)
	if err != nil {
		return
	}
	fmt.Println(bytes)

	var req api.HelloRequest
	proto.Unmarshal(bytes, &req)
	fmt.Println(req.Name)
}
