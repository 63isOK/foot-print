# mock

go官方的mock库,和go版本兼容,针对接口进行mock.

使用步骤:

1. 定义接口
2. mockgen生成接口的mock对象
3. test中使用mock对象测试

对于一个方法,可mock多个实现,调用顺序有两种指定方式:

```go
// After
firstCall := mockObj.EXPECT().SomeMethod(1, "first")
secondCall := mockObj.EXPECT().SomeMethod(2, "second").After(firstCall)
mockObj.EXPECT().SomeMethod(3, "third").After(secondCall)

// InOrder
gomock.InOrder(
    mockObj.EXPECT().SomeMethod(1, "first"),
    mockObj.EXPECT().SomeMethod(2, "second"),
    mockObj.EXPECT().SomeMethod(3, "third"),
)
```

常用选项:

```txt
type Call
  // 可指定mock方法调用的次数,可结合某些业务场景来做测试:
  // 次数不对就报错
  func (c *Call) Times(n int) *Call
  func (c *Call) AnyTimes() *Call
  func (c *Call) MaxTimes(n int) *Call
  func (c *Call) MinTimes(n int) *Call

  // 返回 - mock
  func (c *Call) Return(rets ...interface{}) *Call
  // 返回 - stub
  func (c *Call) Do(f interface{}) *Call
  func (c *Call) DoAndReturn(f interface{}) *Call

  // 替换入参,前提是入参可修改(指针/切片/接口)
  func (c *Call) SetArg(n int, value interface{}) *Call

  // mock的条件
type Matcher
  func All(ms ...Matcher) Matcher
  func Any() Matcher
  func Eq(x interface{}) Matcher
  func Not(x interface{}) Matcher
  func Nil() Matcher
  func Len(i int) Matcher
```

mock主要针对接口提供了如下功能:

1. mock和stub(体现在Do/Return/DoAndReturn)
2. 次数指定
3. 条件

总的来说,功能齐全,暴露的方法也不复杂(就3类).
使用上,需要添加go generate,也就几行的事.
