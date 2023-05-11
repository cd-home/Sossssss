package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/protocol"
	"github.com/smallnest/rpcx/server"
	"log"
	"rpcx/etcd/api"
	"time"
)

var (
	addr     = flag.String("addr", "127.0.0.1", "server address")
	port     = flag.String("port", "8972", "server port")
	etcdAddr = flag.String("etcdAddr", "10.211.55.18:2379", "etcd address")
	basePath = flag.String("base", "/rpcx", "prefix path")
)

var address string

// Arith 服务
type Arith struct {
}

func (t *Arith) Mul(ctx context.Context, args *api.Args, reply *api.Reply) error {
	reply.C = args.A + args.B
	reply.Server = fmt.Sprintf("%s", address)
	return nil
}

func main() {
	flag.Parse()

	address = fmt.Sprintf("%s:%s", *addr, *port)

	// 设置参数
	s := server.NewServer(
		server.WithReadTimeout(time.Second),
		server.WithWriteTimeout(time.Second),
	)

	addRegistryPlugin(s)

	// metadata一般为空, 可以设置 例如state=inactive 禁用服务  group=my_group
	_ = s.Register(new(Arith), "")

	s.AuthFunc = auth

	_ = s.Serve("tcp", address)
}

// addRegistryPlugin 服务注册
func addRegistryPlugin(s *server.Server) {
	log.Println(address)
	r := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: "tcp@" + address,
		EtcdServers:    []string{*etcdAddr},
		BasePath:       *basePath,
		UpdateInterval: time.Minute,
	}
	err := r.Start()
	if err != nil {
		log.Fatal(err)
	}
	s.Plugins.Add(r)
}

// auth 认证
func auth(ctx context.Context, req *protocol.Message, token string) error {
	return nil
}
