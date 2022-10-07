package main

import (
	"fmt"
	"sqlx/mysql/db"
)

func main() {
	stats := db.DB.Stats()
	fmt.Println(stats)
}
