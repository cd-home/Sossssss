package main

import (
	"context"
	"fmt"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c, _, err := websocket.Dial(ctx, "ws://127.0.0.1:9999/ws", nil)
	if err != nil {
		panic(err)
	}
	defer c.Close(websocket.StatusInternalError, "Error")

	err = wsjson.Write(ctx, c, "WebSocket Server")
	if err != nil {
		panic(err)
	}

	var v string
	err = wsjson.Read(ctx, c, &v)
	if err != nil {
		panic(err)
	}
	fmt.Printf("服务端响应：%v\n", v)

	c.Close(websocket.StatusNormalClosure, "")
}
