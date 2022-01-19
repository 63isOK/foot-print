package main

import "fmt"

// funcs 测试函数入口
var funcs = make([]func(), 0)

func main() {
	for _, function := range funcs {
		fmt.Println("============================")
		function()
	}
}
