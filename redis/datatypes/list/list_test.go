package list

import (
	"context"
	"github.com/go-redis/redis/v9"
	"redis/client"
	"testing"
	"time"
)

var cli *redis.Client
var ctx = context.Background()

func init() {
	cli = client.NewSimpleClient()
}

func TestLPush(t *testing.T) {
	err := cli.LPush(ctx, "names", "yao").Err()
	if err != nil {
		t.Error(err)
	}
	err = cli.LPush(ctx, "names", "mike").Err()
	if err != nil {
		t.Error(err)
	}
	// 设置过期
	//cli.Expire(ctx, "names", time.Second*5)
}

func TestRPop(t *testing.T) {
	t.Log(cli.RPop(ctx, "names").String())
}

func TestBRPop(t *testing.T) {
	t.Log(cli.BRPop(ctx, time.Second*10, "names").String())
}
