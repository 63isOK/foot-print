package main

import (
	"reflect"
	"testing"

	. "github.com/agiledragon/gomonkey/v2"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMethod(t *testing.T) {
	Convey("测试方法", t, func() {
		// Convey("private.Method", func() {
		//   var p *privateObject
		//   patches := ApplyMethod(reflect.TypeOf(p), "GetName", func(*privateObject) string {
		//     return "private"
		//   })
		//   defer patches.Reset()
		//
		//   obj := new(privateObject)
		//   ret := obj.GetName()
		//   So(ret, ShouldEqual, "private")
		// })

		Convey("private.method 私有方法", func() {
			var p *privateObject
			patches := ApplyPrivateMethod(reflect.TypeOf(p), "setName", func(obj *privateObject) {
				obj.name = "123"
			})
			defer patches.Reset()

			obj := new(privateObject)
			obj.setName("123")
			So(obj.name, ShouldEqual, "123")
		})
	})
}
