package main

import (
	"context"
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
}

func generateBucketName() (<-chan string, chan struct{}) {
	bucketName := make(chan string)
	stop := make(chan struct{})
	go func() {
		defer close(stop)
		i := 0
		for {
			select {
			case <-stop:
				return
			case bucketName <- BucketBase + strconv.Itoa(i):
				if i >= BucketSize {
					i = 1
				} else {
					i++
				}
			}
		}
	}()
	return bucketName, stop
}

func StartTimers() {
	timers = make([]*time.Ticker, BucketSize)
	for i := 0; i < BucketSize; i++ {
		timer := time.NewTicker(time.Second)
		timers[i] = timer
		bucketName := <-bucketNameCh
		stop := make(chan struct{})
		go func() {
			defer close(stop)
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
		err := bucket.Get(ctx, bucketName)
		if err != nil {
			log.Println(err)
			return
		}
		if bucket.TimeStamp > float64(now.Unix()) {
			return
		}
		job := &Job{}
		err = job.Get(ctx, bucket.JobID)
		if err != nil {
			return
		}

	}
}
