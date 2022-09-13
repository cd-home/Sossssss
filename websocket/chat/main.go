package main

import (
	"context"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"strconv"
	"time"
)

func main() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, req *http.Request) {
		conn, err := websocket.Accept(w, req, &websocket.AcceptOptions{
			InsecureSkipVerify: true,
		})
		if err != nil {
			log.Println(err)
			return
		}
		ctx := context.Background()
		go func() {
			for {
				var v interface{}
				err = wsjson.Read(ctx, conn, &v)
				if err != nil {
					log.Println(err)
					return
				}
				log.Printf("From Client：%v\n", v)
			}
		}()
		for i := 0; i < 10; i++ {
			err = wsjson.Write(ctx, conn, "To WebSocket Client "+strconv.Itoa(i))
			time.Sleep(time.Second)
			if err != nil {
				log.Println(err)
				return
			}
		}
		//conn.Close(websocket.StatusNormalClosure, "")
	})

	log.Fatal(http.ListenAndServe(":9999", nil))
}
