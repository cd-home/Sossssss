package main

import (
	"log"
	"zzchat/internal/router"
)

func main() {
	// rpc service

	// api service
	log.Fatal(router.SetUpRouter().Run(":8080"))
}
