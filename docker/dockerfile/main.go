package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		log.Println(r.URL)
		_, _ = io.Copy(rw, bytes.NewReader([]byte("Hello K8s")))
	})
	log.Println("Listen: 8999")
	_ = http.ListenAndServe(":8999", nil)
}
