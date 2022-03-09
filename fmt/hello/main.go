package main

import "fmt"

func main() {
	s := "1" + fmt.Sprint(123) + "2"
	fmt.Println(s)

	s += fmt.Sprint(struct{ A, B int }{1, 2})
	s += "hello"
	fmt.Println(s)

	s += fmt.Sprint("world")
	s += "!"
	fmt.Println(s)

	fmt.Println(fmt.Sprint(123, "hello"))
	fmt.Println(fmt.Sprint("hi", "hello"))
	fmt.Println(fmt.Sprint("hi", "hello", "world"))
}
