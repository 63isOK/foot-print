package main

import (
	"fmt"
	"unsafe"
)

func main() {
	var i uint64 = 12

	fmt.Println(float2int(i))
}

func float2int(i uint64) float64 {
	return *(*float64)(unsafe.Pointer(&i))
}
