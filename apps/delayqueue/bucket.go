package main

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v9"
)

type Bucket struct {
	TimeStamp float64
	JobID     string
}

func (b *Bucket) Add(ctx context.Context, key string, score float64, jobID string) error {
	return rdb.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: jobID,
	}).Err()
}

func (b *Bucket) Get(ctx context.Context, key string) error {
	z, err := rdb.ZRangeWithScores(ctx, key, 0, 0).Result()
	if errors.Is(err, redis.Nil) || len(z) == 0 {
		return KeyNotExist
	}
	if err != nil {
		return err
	}

	b.TimeStamp = z[0].Score
	b.JobID = z[0].Member.(string)

	return nil
}

func (b *Bucket) Remove(ctx context.Context, key string, jobID string) error {
	return rdb.ZRem(ctx, key, jobID).Err()
}
