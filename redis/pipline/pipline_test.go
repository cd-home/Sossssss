package pipline

import (
	"context"
	"github.com/go-redis/redis/v9"
	"redis/client"
	"testing"
	"time"
)

var cli *redis.Client
var ctx context.Context

func init() {
	cli = client.NewSimpleClient()
	ctx = context.Background()
}

func TestPipLine(t *testing.T) {
	pipe := cli.Pipeline()

	pipe.SAdd(ctx, "s4", 1, 2, 3)
	res := pipe.SMembers(ctx, "s4")

	_, err := pipe.Exec(ctx)

	if err != nil {
		t.Fatal(err)
	}

	// 获取结果
	t.Log(res.Val())

	// better mode
	var incr *redis.IntCmd
	_, err = cli.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		incr = pipe.Incr(ctx, "counter")
		pipe.Expire(ctx, "counter", time.Second*10)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(incr.Val())
}

// TODO 事务
func TestTransaction(t *testing.T) {
	//var incr *redis.IntCmd
	//cli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
	//	incr = cli.Incr(ctx, "counter")
	//	return nil
	//})
}
