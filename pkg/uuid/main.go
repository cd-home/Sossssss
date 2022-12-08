package main

import (
	"fmt"
	"github.com/google/uuid"
)

func main() {
	// v1 一般情况下使用 v1 v4 即可
	id, _ := uuid.NewUUID()
	fmt.Println(id, id.Version().String())

	// v4
	id = uuid.New()
	fmt.Println(id, id.Version().String())

	// v2 group
	id, _ = uuid.NewDCEGroup()
	fmt.Println(id, id.Version().String())

	// v2 person
	id, _ = uuid.NewDCEPerson()
	fmt.Println(id, id.Version().String())

	// v3
	pre, _ := uuid.NewDCEPerson()
	id = uuid.NewMD5(pre, []byte("secret 1"))
	fmt.Println(id, id.Version().String())

	// v5
	pre, _ = uuid.NewDCEPerson()
	id = uuid.NewSHA1(pre, []byte("secret 2"))
	fmt.Println(id, id.Version().String())
}
