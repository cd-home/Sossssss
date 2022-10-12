package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"sync"
	"time"
)

var ChatManager *Chat

type Chat struct {
	mu         sync.Mutex
	subscriber map[*subscriber]struct{}
}

type subscriber struct {
	ws    *websocket.Conn
	msg   chan []byte
	close func()
	stop  <-chan struct{}
}

func init() {
	ChatManager = &Chat{
		subscriber: make(map[*subscriber]struct{}),
	}
}

func main() {
	// static file service
	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/sub", Subscribe)
	http.HandleFunc("/pub", Publisher)

	_ = http.ListenAndServe(":8080", nil)
}

func Subscribe(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		log.Println(err)
		return
	}
	suber := &subscriber{
		ws:  ws,
		msg: make(chan []byte),
		close: func() {
			//ws.Close(websocket.StatusPolicyViolation, "conn slow")
		},
	}
	ChatManager.mu.Lock()
	ChatManager.subscriber[suber] = struct{}{}
	ChatManager.mu.Unlock()
	go subscribe(context.Background(), suber)
}

func subscribe(ctx context.Context, suber *subscriber) {
	heartBeat := func() (<-chan time.Time, <-chan struct{}) {
		ch := make(chan time.Time, 1)
		tk := time.NewTicker(2 * time.Second)
		stop := make(chan struct{})
		go func() {
			defer tk.Stop()
			defer close(stop)
			for {
				select {
				case <-stop:
					return
				case t := <-tk.C:
					ch <- t
				}
			}
		}()
		return ch, stop
	}
	beat, stop := heartBeat()
	suber.stop = stop
	for {
		select {
		case msg := <-suber.msg:
			//wsjson.Write(ctx, suber.ws, msg)
			_ = suber.ws.Write(ctx, websocket.MessageText, msg)
		case t := <-beat:
			_ = suber.ws.Write(ctx, websocket.MessageText, []byte(t.Format("20060102150405")))
		}
	}
}

func Publisher(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 200)
		return
	}
	// 广播
	for suber := range ChatManager.subscriber {
		select {
		case suber.msg <- body:
		}
	}
	w.WriteHeader(http.StatusAccepted)
}
