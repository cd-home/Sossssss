package main

import (
	"log"
	"zzchat/internal/router"
)

func main() {
	log.Fatal(router.SetUpRouter().Run(":8080"))
}
