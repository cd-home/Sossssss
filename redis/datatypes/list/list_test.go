package list

import (
	"context"
	"github.com/go-redis/redis/v9"
	"redis/cluster"
	"testing"
	"time"
)

var cli *redis.ClusterClient
var ctx = context.Background()

func init() {
	cli = cluster.NewClusterClient()
}

// TestLPush 列表 LPush
func TestLPush(t *testing.T) {
	err := cli.LPush(ctx, "webapp_names", "yao").Err()
	if err != nil {
		t.Error(err)
	}
	err = cli.LPush(ctx, "webapp_names", "mike").Err()
	if err != nil {
		t.Error(err)
	}
	// 设置过期
	cli.Expire(ctx, "webapp_names", time.Second*10)
}

// TestRPop 列表 RPop
func TestRPop(t *testing.T) {
	t.Log(cli.RPop(ctx, "names").String())
}

// TestBRPop 阻塞读取
func TestBRPop(t *testing.T) {
	t.Log(cli.BRPop(ctx, time.Second*10, "names").String())
}
