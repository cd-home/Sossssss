package main

import (
	"context"
	"github.com/go-redis/redis/v9"
)

func NewClusterClient() *redis.ClusterClient {
	// if readOnly == true   slave
	// if readOnly == true and RouteByLatency 选择对应slot 延迟最低的master 或者 slave
	// if readOnly == true and RouteRandomly  随机选择对应slot master 或者 slave
	cli := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:          []string{":7000", ":7001"},
		ReadOnly:       true,
		RouteByLatency: false,
		RouteRandomly:  false,
	})
	if err := cli.Ping(context.Background()).Err(); err != nil {
		return cli
	}
	return nil
}
