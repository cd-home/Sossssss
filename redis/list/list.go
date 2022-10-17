package list

import (
	"context"
	"redis/client"
)

func LPush(ctx context.Context, key string, val string) error {
	cli := client.NewSimpleClient()
	return cli.LPush(ctx, key, val).Err()
}

func RPop(ctx context.Context, key string) (string, error) {
	cli := client.NewSimpleClient()
	return cli.RPop(ctx, key).Result()
}
