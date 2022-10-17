package hash

import (
	"context"
	"github.com/go-redis/redis/v9"
	"redis/client"
	"testing"
)

var cli *redis.Client
var ctx = context.Background()

func init() {
	cli = client.NewSimpleClient()
}

func TestHSet(t *testing.T) {
	err := cli.HSet(ctx, "yao_", "age", "20", "len", "18").Err()
	if err != nil {
		t.Error(err)
	}
	// 存在即是不设置
	//cli.HSetNX()
}

func TestHGet(t *testing.T) {
	t.Log(cli.HGet(ctx, "yao_", "age").String())
	t.Log(cli.HGet(ctx, "yao_", "len").String())
}

func TestHGetAll(t *testing.T) {
	t.Log(cli.HGetAll(ctx, "yao_").Result())
}
