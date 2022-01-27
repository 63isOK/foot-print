package main

import (
	"testing"

	. "github.com/agiledragon/gomonkey/v2"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPrivateFunc(t *testing.T) {
	Convey("测试私有函数", t, func() {
		Convey("测试无参无返回值的私有函数", func() {
			defer func() {
				err := recover()
				So(err, ShouldBeNil)
			}()

			patches := ApplyFunc(privateFunc, func() {
				println("privateFunc - change")
			})
			defer patches.Reset()

			privateFunc()
		})

		Convey("测试带参的私有函数", func() {
			patches := ApplyFunc(privateFuncSum, func(int, int) int {
				return 5
			})
			defer patches.Reset()

			result := privateFuncSum(1, 2)
			So(result, ShouldEqual, 5)
			result = privateFuncSum(5, 2)
			So(result, ShouldEqual, 5)
			result = privateFuncSum(5, 2)
			So(result, ShouldEqual, 5)
		})

		Convey("测试带参的私有函数(参数类型不同)", func() {
			patches := ApplyFunc(privateFuncSwap, func(int, float64) (float64, int) {
				return 3.2, 5
			})
			defer patches.Reset()

			i, j := privateFuncSwap(1, 1.2)
			So(i, ShouldEqual, 3.2)
			So(j, ShouldEqual, 5)

			i, j = privateFuncSwap(100, 13.4)
			So(i, ShouldEqual, 3.2)
			So(j, ShouldEqual, 5)
		})
	})
}
