package main

import (
	"apps/ftp/config"
	"net/http"
)

func main() {

	_ = http.ListenAndServe(config.PORT, nil)
}
