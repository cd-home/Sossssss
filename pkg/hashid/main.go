package main

import (
	"fmt"
	"github.com/speps/go-hashids/v2"
)

func main() {
	hd := hashids.NewData()
	hd.Salt = "this is salt"
	hd.MinLength = 30
	h, _ := hashids.NewWithData(hd)
	e, _ := h.Encode([]int{10})
	fmt.Println(e)

	d, _ := h.DecodeWithError(e)
	fmt.Println(d)
}
