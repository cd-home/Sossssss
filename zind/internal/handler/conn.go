package handler

import (
	"github.com/gin-gonic/gin"
	"nhooyr.io/websocket"
)

type Conn struct {
}

func (ws *Conn) UpgradeWebsocket(ctx *gin.Context) {
	conn, err := websocket.Accept(ctx.Writer, ctx.Request, &websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		ctx.JSON(0, "Not Support WebSocket")
		return
	}
	// 链接添加到管理中心
	conn.Subprotocol()
	// 开启读写 goroutine

	ctx.JSON(200, "OK")
}
