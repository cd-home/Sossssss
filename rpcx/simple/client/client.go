package main

import (
	"context"
	"github.com/smallnest/rpcx/client"
	"log"
	"rpcx/simple/api"
)

func main() {
	d, _ := client.NewPeer2PeerDiscovery("tcp@127.0.0.1:8972", "")
	xClient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xClient.Close()

	args := &api.Args{
		A: 1,
		B: 2,
	}
	reply := &api.Reply{}

	// 同步调用
	_ = xClient.Call(context.Background(), "Mul", args, reply)
	log.Println(reply.C)

	// 异步调用
	ch, _ := xClient.Go(context.Background(), "Mul", args, reply, nil)
	<-ch.Done
	r := ch.Reply
	res := r.(*api.Reply)
	log.Println(res.C)
}
