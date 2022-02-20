package main

import (
	"context"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	pipe := make(chan int)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case number := <-pipe:
				println(number)
				time.Sleep(time.Second * 5)
			}
		}
	}()

	for i := 0; i < 10; i++ {
		pipe <- i
	}
	cancel()
}
