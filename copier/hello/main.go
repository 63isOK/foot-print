package main

import (
	"fmt"

	"github.com/jinzhu/copier"
)

type src struct {
	A int
	B bool
	C string
}

func (s src) D() int {
	if s.B {
		return 1
	}

	return 0
}

type dest struct {
	A int
	D int
	C int
}

func (d dest) B() bool {
	return d.D != 0
}

func main() {
	foo()
	bar()
}

func foo() {
	println("==============")
	s := src{}
	d := dest{
		A: 1,
		D: 1,
		C: 3,
	}
	copier.Copy(&s, &d)
	fmt.Printf("%+v", s)
}

func bar() {
	println("==============")
	s := src{
		A: 1,
		B: false,
		C: "2",
	}
	d := dest{}
	copier.Copy(&d, &s)
	fmt.Printf("%+v", d)
}
