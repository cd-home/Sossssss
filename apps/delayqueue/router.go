package main

import "github.com/gin-gonic/gin"

func initRouter() *gin.Engine {
	e := gin.Default()
	e.POST("/job", AddJob)
	e.GET("/job", PopJob)
	return e
}
