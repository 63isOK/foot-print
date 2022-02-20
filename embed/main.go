package main

import "fmt"

type User struct {
	Name string
	Age  int
}

type Family struct {
	Address string
	User
}

func main() {
	f := Family{
		Address: "123 Main St",
		User: User{
			Name: "John",
			Age:  30,
		},
	}

	fmt.Println(f)
}
