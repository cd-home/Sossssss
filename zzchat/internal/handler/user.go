package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
)

type UserHandler struct {
	Rds *redis.Client
}

func (u *UserHandler) Login(ctx *gin.Context) {

}

func (u *UserHandler) Register(ctx *gin.Context) {

}

func (u *UserHandler) Logout(ctx *gin.Context) {

}

func (u *UserHandler) LogOff(ctx *gin.Context) {

}
