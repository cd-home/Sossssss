package main

import (
	"context"
	"flag"
	etcdClient "github.com/rpcxio/rpcx-etcd/client"
	"github.com/smallnest/rpcx/client"
	"log"
	"rpcx/etcd/api"
	"time"
)

var (
	etcdAddr = flag.String("etcdAddr", "10.211.55.18:2379", "etcd address")
	basePath = flag.String("base", "/rpcx", "prefix path")
)

func main() {
	flag.Parse()
	// 服务发现
	d, _ := etcdClient.NewEtcdV3Discovery(*basePath, "Arith", []string{*etcdAddr}, false, nil)

	// 轮询
	option := client.DefaultOption
	// 设置心跳
	option.Heartbeat = true
	option.HeartbeatInterval = time.Second
	xClient := client.NewXClient("Arith", client.Failover, client.RoundRobin, d, option)

	// Auth Token
	xClient.Auth("")

	defer xClient.Close()

	args := &api.Args{
		A: 1,
		B: 2,
	}

	reply := &api.Reply{}

	for {
		// 同步调用
		err := xClient.Call(context.Background(), "Mul", args, reply)
		if err != nil {
			log.Println(err)
		}
		log.Println(reply.C)
		log.Println(reply.Server)
		time.Sleep(time.Second * 2)
	}
}
