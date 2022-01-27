package main

type demo int

func New() demo {
	return demo(0)
}

func (d demo) Hello() demo {
	println("Hello")
	return d
}

func (d demo) World() demo {
	println("World")
	return d
}

func (d demo) End() demo {
	println("End")
	return d
}

func main() {
	defer New().Hello().World().End()

	println("end...")
}
