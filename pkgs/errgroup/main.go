package main

import (
	"context"
	"errors"
	"github.com/neilotoole/errgroup"
	"log"
)

func main() {
	g, _ := errgroup.WithContextN(context.Background(), 4, 50)

	for i := 0; i < 60; i++ {
		j := i
		g.Go(func() error {
			if j%7 == 0 {
				return errors.New("error j%7==0")
			}
			log.Println(j)
			return nil
		})
	}
	e := g.Wait()
	log.Println(e)
}
