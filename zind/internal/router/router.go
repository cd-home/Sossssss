package router

import (
	"github.com/gin-gonic/gin"
	"redis/client"
	"zind/internal/handler"
)

func SetUpRouter() *gin.Engine {
	e := gin.Default()

	// setUp redis
	rds := client.NewSimpleClient()

	// inject
	wsh := handler.Conn{}
	user := handler.UserHandler{Rds: rds}

	e.GET("/ws", wsh.UpgradeWebsocket)
	e.POST("/login", user.Login)

	return e
}
