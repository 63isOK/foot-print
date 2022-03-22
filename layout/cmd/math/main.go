package main

import (
	"github.com/63isOK/foot-print/layout/internal/pkg/count"
	// "github.com/63isOK/foot-print/layout/internal/pkg/internal/sum"
)

func main() {
	println(count.Count(1, 2))
	// 不能引用internal/*/internal/包
	// println(sum.Sum(1, 2))
}
