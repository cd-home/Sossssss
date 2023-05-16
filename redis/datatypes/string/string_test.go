package string

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

// TestStringSet 设置 string
func TestStringSet(t *testing.T) {
	// cli.SetEx()
	err := cli.Set(ctx, "stock", 10, time.Second*10).Err()
	if err != nil {
		t.Log(err)
	}
}

// TestStringGet 获取 string
func TestStringGet(t *testing.T) {
	value, err := cli.Get(ctx, "stock").Result()
	if err != nil {
		t.Error(err)
	}
	t.Log(value)
}

// TestStringSetNX 不存在即是设置
func TestStringSetNX(t *testing.T) {
	for i := 10; i > 0; i-- {
		// 不存在即设置, 否则不做任何操作
		err := cli.SetNX(ctx, "stock", i, time.Second*10).Err()
		if err != nil {
			t.Error(err)
		}
	}
	// key存在时才能设置
	//cli.SetXX()
}
