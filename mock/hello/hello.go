package hello

//go:generate mockgen -destination mock_hello_test.go -package hello_test github.com/63isOK/foot-print/mock/hello Hi

// UserInfo user info.
type UserInfo struct {
	Name string
	Age  int
}

// Hi returns a friendly greeting.
type Hi interface {
	Hello(name string) string
	HelloAgain(name, tinyName string) UserInfo
}

// SayHi say hi to someone.
func SayHi(h Hi, name string) string {
	return h.Hello(name)
}

// SayHiAgain say hi to someone.
func SayHiAgain(h Hi, name, tinyName string) UserInfo {
	return h.HelloAgain(name, tinyName)
}
