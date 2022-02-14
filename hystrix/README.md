# hystrix-go

熔断,通过隔离不同的服务,来防止某一服务失败导致整个项目失败.
熔断是框架机制,不管是核心服务还是边缘服务,都需要配上.

熔断更多时候是解决流量过多的场景,也被称为"过载保护",和电路上的保险丝一样,
应用场景是"过载".

hystrix-go取法于网飞的hystrix,通过将远端系统/服务/第三方库的访问隔离起来,
停止级联错误,增加复杂分布式系统的弹性,容忍延时和错误的库.

## 使用

***依赖外部系统的逻辑,作为command,丢在hystrix.Go中***

```go
hystrix.Go("my_command", func() error {
  // 调用其他服务.当其他服务是健康时,会执行此处的代码
  return nil
}, nil)
```

***定义后备行为***

```go
hystrix.Go("my_command", func() error {
  // 调用其他行为
  return nil
}, func(err error) error {
  // 当服务不健康或超时时,执行预定义的行为
  return nil
})
```

***同步***

hystrix.Go的行为类似于协程,如果要实现同步效果,有两种方式:

1. 额外使用chan完成
2. 使用hystrix.Do代替hystrix.Go

***配置***

```go
type CommandConfig struct {
  Timeout                int // command执行超时时长
  MaxConcurrentRequests  int // 最大并发量
  RequestVolumeThreshold int // 10s统计窗口内请求数阈值,达到这个阈值后才会判断是否要开启熔断
  SleepWindow            int // 熔断多久后尝试服务是否健康
  ErrorPercentThreshold  int // 错误率阈值,请求数/错误率都达到阈值后,开启熔断
}
```

