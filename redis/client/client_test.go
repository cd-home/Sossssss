package client

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v9"
	"testing"
	"time"
)

func TestRedisSimpleClient(t *testing.T) {
	rdb := NewSimpleClient()
	err := rdb.Set(context.Background(), "name", "yao", time.Second*20).Err()
	if err != nil {
		t.Error(err)
	}
	value, err := rdb.Get(context.Background(), "name").Result()
	if errors.Is(err, redis.Nil) {
		t.Log("Key does not exist")
	} else if err != nil {
		t.Error(err)
	} else {
		t.Logf("value: %s", value)
	}
}

func TestZSET(t *testing.T) {
	rdb := NewSimpleClient()
	members := []redis.Z{
		{
			Score:  10,
			Member: "yao",
		},
		{
			Score:  20,
			Member: "mike",
		},
	}

	err := rdb.ZAdd(context.Background(), "bk", members...).Err()
	if err != nil {
		t.Error(err)
	}

	zs, err := rdb.ZRangeWithScores(context.Background(), "bk", 0, 0).Result()
	if err != nil {
		t.Error(err)
	}
	score := zs[0].Score
	member := zs[0].Member.(string)

	t.Log(score)
	t.Log(member)

	err = rdb.ZRem(context.Background(), "bk", "mike").Err()
	if err != nil {
		t.Error(err)
	}
}
