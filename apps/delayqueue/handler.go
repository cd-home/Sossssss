package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddJob(ctx *gin.Context) {
	var job Job
	err := ctx.ShouldBindJSON(&job)
	if err != nil {
		ctx.JSON(http.StatusOK, err)
		return
	}
	err = job.Add(context.Background(), job.ID)
	if err != nil {
		ctx.JSON(http.StatusOK, err)
		return
	}
	bucket := &Bucket{JobID: job.ID, TimeStamp: job.Delay}
	err = bucket.Add(context.Background(), <-bucketNameCh, job.Delay, job.ID)
	if err != nil {
		ctx.JSON(http.StatusOK, err)
		return
	}
	ctx.JSON(http.StatusOK, "OK")
}

func PopJob(ctx *gin.Context) {
	topic := ctx.Query("topic")
	rq := &ReadyQ{Topic: topic}
	jobID, err := rq.Pop(context.Background())
	if err != nil {
		ctx.JSON(http.StatusOK, err)
		return
	}
	job := &Job{}
	err = job.Get(context.Background(), jobID)
	if err != nil {
		ctx.JSON(http.StatusOK, err)
		return
	}
	ctx.JSON(http.StatusOK, job)
}
