package main

func main() {
	f1 := foo(1)
	f1()
	f1()
	f1()

	f2 := foo(2)
	f2()
	f2()
	f2()

	f1()
	f1()
	f1()

	f2()
	f2()
	f2()

	println("===========")
	b1 := boo(1)
	b1()
	b1()
	b1()

	b2 := boo(2)
	b2()
	b2()
	b2()

	b1()
	b1()
	b1()
	b2()
	b2()
	b2()
}

func foo(i int) func() {
	return func() {
		println(i)
	}
}

func boo(i int) func() {
	callback := func() func() {
		return func() {
			println(i)
		}
	}

	return callback()
}
