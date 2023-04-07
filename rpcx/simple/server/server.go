package main

import (
	"context"
	"github.com/smallnest/rpcx/server"
	"rpcx/simple/api"
)

type Arith struct {
}

func (t *Arith) Mul(ctx context.Context, args *api.Args, reply *api.Reply) error {
	reply.C = args.A + args.B
	return nil
}

func main() {
	s := server.NewServer()
	//s.RegisterName("Arith", &Arith{}, "")
	err := s.Register(&Arith{}, "")
	if err != nil {
		return
	}
	err = s.Serve("tcp", ":8972")
	if err != nil {
		return
	}
}
