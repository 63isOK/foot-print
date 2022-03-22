package main

import (
	"fmt"
	"time"
)

func main() {
	var a, b string
	a = "abc"
	b = ""
	println(a > b)
	println(a < b)
	println(a == b)
	Hello()
}

func Hello() {
	fmt.Println("b")
	fmt.Println("b")
	println("c")
}

func a() {
	i := 0
	for {
		i++
		println("aa")
		select {
		default:
		case <-time.After(time.Second):
			if i > 300 {
				for i < 400 {
					println("aa")
				}
			}
		}
	}
}
