package main

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
)

func main() {
	id := uuid.NewV4()
	fmt.Println(id, id.Version())
}
