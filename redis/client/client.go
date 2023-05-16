package client

import (
	"context"
	"github.com/go-redis/redis/v9"
	"log"
	"time"
)

func NewSimpleClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Network:      "tcp",
		Addr:         "10.211.55.18:6379",
		Password:     "",
		DB:           0,
		PoolSize:     10,
		MinIdleConns: 5,
		MaxIdleConns: 10,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		PoolTimeout:  2 * time.Second,
	})
	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Redis Connection Successful: " + pong)
	return rdb
}
