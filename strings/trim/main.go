package main

import "strings"

func main() {
	println(strings.TrimPrefix("aababcd", "aab"))
	println(strings.TrimLeft("aababcd", "ab"))
}
