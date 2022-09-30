package main

import (
	"errors"
	"github.com/go-redis/redis/v9"
	"redis/client"
)

var rdb *redis.Client

var (
	KeyNotExist = errors.New("KeyNotExist")
)

func init() {
	rdb = client.NewSimpleClient()
}
