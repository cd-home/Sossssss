package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v9"
)

const JobPrefix = "Job_"

type Job struct {
	Topic string  `json:"topic"`
	ID    string  `json:"id"`
	Delay float64 `json:"delay"`
	TTR   float64 `json:"ttr"`
	Body  string  `json:"body"`
}

func (j *Job) Get(ctx context.Context, key string) error {
	jobBytes, err := rdb.Get(ctx, JobPrefix+key).Bytes()
	if errors.Is(err, redis.Nil) {
		return KeyNotExist
	}
	err = json.Unmarshal(jobBytes, j)
	if err != nil {
		return err
	}
	return nil
}

func (j *Job) Add(ctx context.Context, key string) error {
	jobBytes, err := json.Marshal(j)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, JobPrefix+key, jobBytes, redis.KeepTTL).Err()
}

func (j *Job) Remove(ctx context.Context, key string) error {
	return rdb.Del(ctx, JobPrefix+key).Err()
}
