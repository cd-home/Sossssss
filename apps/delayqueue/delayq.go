package main

import (
	"context"
	"errors"
	"log"
	"strconv"
	"time"
)

const (
	BucketSize = 5
	BucketBase = "Bucket_"
)

var (
	timers           []*time.Ticker
	bucketNameCh     <-chan string
	stopBucketNameCh chan struct{}
)

func init() {
	bucketNameCh, stopBucketNameCh = generateBucketName()
	StartTimers()
}

func generateBucketName() (<-chan string, chan struct{}) {
	bucketNameCh := make(chan string)
	stop := make(chan struct{})
	go func() {
		defer close(stop)
		i := 0
		for {
			select {
			case <-stop:
				return
			case bucketNameCh <- BucketBase + strconv.Itoa(i):
				if i >= BucketSize {
					i = 1
				} else {
					i++
				}
			}
		}
	}()
	return bucketNameCh, stop
}

func StartTimers() {
	timers = make([]*time.Ticker, BucketSize)
	for i := 0; i < BucketSize; i++ {
		timer := time.NewTicker(time.Second * 2)
		timers[i] = timer
		bucketName := <-bucketNameCh
		stop := make(chan struct{})
		go func() {
			defer close(stop)
			defer timer.Stop()
			for {
				select {
				case now := <-timer.C:
					ScanBucket(now, bucketName)
				}
			}
		}()
	}
}

func ScanBucket(now time.Time, bucketName string) {
	ctx := context.Background()
	bucket := &Bucket{}
	for {
		// Get Bucket
		err := bucket.Get(ctx, bucketName)
		if errors.Is(err, KeyNotExist) {
			return
		}
		if err != nil {
			log.Println(err)
			return
		}
		if bucket.TimeStamp > float64(now.Unix()) {
			return
		}
		// Get Job
		job := &Job{}
		err = job.Get(ctx, bucket.JobID)
		if err != nil {
			log.Printf("Get Job [%s] Failed", bucket.JobID)
			continue
		}
		// Push Ready Queue
		readyQ := &ReadyQ{
			Topic: job.Topic,
			JobId: job.ID,
		}
		err = readyQ.Push(ctx)
		if err != nil {
			log.Printf("Push Job [%s] To ReadyQ Failed", readyQ.JobId)
			continue
		}
		// Delete
		err = bucket.Remove(ctx, bucketName, job.ID)
		if err != nil {
			return
		}
	}
}
