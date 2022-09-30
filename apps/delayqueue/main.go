package main

import "log"

func main() {
	log.Fatal(initRouter().Run(":8080"))
}
