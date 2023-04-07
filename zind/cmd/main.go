package main

import (
	"log"
	"zind/internal/router"
)

func main() {
	// rpc service

	// api service
	log.Fatal(router.SetUpRouter().Run(":8080"))
}
