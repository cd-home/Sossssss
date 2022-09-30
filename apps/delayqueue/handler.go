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
	bucket := &Bucket{JobID: job.ID, TimeStamp: float64(job.Delay)}
	err = bucket.Add(context.Background(), <-bucketNameCh, float64(job.Delay), job.ID)
	if err != nil {
		ctx.JSON(http.StatusOK, err)
		return
	}
	ctx.JSON(http.StatusOK, "OK")
}
