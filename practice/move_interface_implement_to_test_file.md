# 在测试中验证接口的实现

```go
func TestInterface_Impl(t *testing.T) {
  var _ interface = new(impl)
}
```

检测代码放在普通源码中,会在编译期触发检查.
`还是放在test中比较合适,毕竟test就是用来干这个事的.`
