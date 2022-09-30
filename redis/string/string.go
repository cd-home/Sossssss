package string

import (
	"context"
	"github.com/go-redis/redis/v9"
	"redis/client"
)

var cli *redis.Client

func init() {
	cli = client.NewSimpleClient()
}

func Get(key string) string {
	ctx := context.Background()
	value, err := cli.Get(ctx, key).Result()
	if err != nil {
		return ""
	}
	return value
}
