package hello

// NOTE: step1: 定义接口

// Greeter is a hi
type Greeter interface {
	Greet() (string, error)
	GreetAgain() string
}
