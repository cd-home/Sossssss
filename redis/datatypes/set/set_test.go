package set

import (
	"context"
	"github.com/go-redis/redis/v9"
	"redis/client"
	"testing"
)

var cli *redis.Client
var ctx context.Context

func init() {
	cli = client.NewSimpleClient()
	ctx = context.Background()
}

func TestSetAdd(t *testing.T) {
	err := cli.SAdd(ctx, "group", "a", "b", "c").Err()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetRem(t *testing.T) {
	err := cli.SRem(ctx, "group", "a").Err()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetMembers(t *testing.T) {
	data, err := cli.SMembers(ctx, "group").Result()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data)
}

func TestSetGet(t *testing.T) {
	err := cli.SPop(ctx, "group").Err()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSisMember(t *testing.T) {
	b, err := cli.SIsMember(ctx, "group", "a").Result()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(b)
}

func TestSetSinter(t *testing.T) {
	cli.SAdd(ctx, "s1", 1, 2, 3, 4)
	cli.SAdd(ctx, "s2", 2, 3, 4, 5)
	cli.SAdd(ctx, "s3", 3, 4, 5, 6)

	// 交集
	inter, err := cli.SInter(ctx, "s1", "s2", "s3").Result()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(inter)

	// 并集
	union, err := cli.SUnion(ctx, "s1", "s2", "s3").Result()
	if err != nil {
		t.Log(err)
	}
	t.Log(union)

	// 差集
	diff, err := cli.SDiff(ctx, "s1", "s2", "s3").Result()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(diff)
}
