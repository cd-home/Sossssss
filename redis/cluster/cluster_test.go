package main

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestNewClusterClient(t *testing.T) {
	cli := NewClusterClient()
	set := cli.Set(context.Background(), "webapp_num", "10", time.Second*60)
	if set.Err() != nil {
		log.Println(set.Err())
	}

	res := cli.Get(context.Background(), "webapp_num")
	log.Println(res.String())
}
