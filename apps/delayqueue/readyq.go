package main

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v9"
	"time"
)

const (
	timeout = time.Second * 2
)

type ReadyQ struct {
	Topic string
	JobId string
}

func (r *ReadyQ) Push(ctx context.Context) error {
	return rdb.RPush(ctx, r.Topic, r.JobId).Err()
}

func (r *ReadyQ) Pop(ctx context.Context) (string, error) {
	jobs, err := rdb.BLPop(ctx, timeout, r.Topic).Result()
	if errors.Is(err, redis.Nil) {
		return "", KeyNotExist
	}
	if err != nil {
		return "", err
	}
	return jobs[0], nil
}
