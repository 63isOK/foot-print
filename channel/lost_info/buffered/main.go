package main

import (
	"os"
	"os/signal"
)

func main() {
	pipe := make(chan os.Signal, 1)
	signal.Notify(pipe, os.Interrupt)
	for i := 0; i < 10; i++ {
		<-pipe
		println(i)
	}
}
