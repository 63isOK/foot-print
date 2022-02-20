# juju/ratelimit源码分析

令牌桶,抛开测试外,全是利用标准库实现的.

## 令牌桶简单介绍

令牌桶是一种限流算法,除了具有限流功能外,还对突增流量有一定弹性.

### 限流的介绍

在这个万物互联的时代,很多服务都依赖了很多外部服务,也被很多外部服务依赖,
当被外部服务依赖时,当前服务就是服务提供商.

> 将`消费者`定义为服务消费者,将`提供者`定义为服务提供者.

在分布式/微服务兴起的时代,服务的高可用被一次次的强调,抛开程序的鲁棒性,
最近几年流量突激的问题作为高可用问题被一次次重视,其中流量突然增大和异常流量
被重点研究,应对方法分别是限流和熔断.

- `限流`:作为提供者,对超出负载能力的流量进行限制处理.
- `熔断`:作为消费者,对异常流量的服务进行断路保护的处理.

这次只介绍限流,下篇文章介绍熔断.

### 限流的常用算法

限流常用的算法有计数器,滑窗,漏桶,以及本次的主角:令牌桶.

每种算法都有各自的优缺点,令牌桶相比其他算法,优点是面对突增的流量有一定弹性,
当桶内的令牌耗光后,会退化为漏桶.

### 令牌桶的实现原理

令牌桶里存放的是令牌,当有请求时,从桶内取出一个令牌,当桶内没有令牌时,对当前
请求做限流处理,具体的处理方式依据不同的需求而不同,最常见的方式是拒绝服务,
这种最暴力的方式也是最有效的,在某些核心场景(如支付)就不能直接拒绝.

向桶内添加令牌是有固定速率和数量的,从长时间来看,令牌添加的速度就是最大负载,
平常时,桶内剩余令牌就是可应付的突增流量.

从上面来看,令牌桶的核心就是计算桶内剩余令牌数.影响剩余令牌数的因数如下:

1. 桶容量,令牌消费速度小于生产速度时,容量就是剩余令牌数
2. 令牌的添加间隔和添加数量,影响了令牌的生产速度

## 源码分析

### 结构分析

供3个文件,ratelimit.go实现了令牌桶;reader.go封装了令牌桶,最后一个是测试.

### 令牌桶实现的分析

> 任何涉及时间流逝的库,为了测试方便,都应该添加接口,启用打桩,方便测试

```go
type Clock interface {
  Now() time.Time
  Sleep(d time.Duration)
}

// realClock 表示真实时钟,非暴露
type realClock struct{}
func (realClock) Now() time.Time {
  return time.Now()
}
func (realClock) Sleep(d time.Duration) {
  time.Sleep(d)
}
```

> 固定时间间隔添加令牌,也不关心何时取令牌,那么时钟滴答tick就是固定时间间隔

将所有涉及到的时间全部转换为tick,方便计算.

```go
func (tb *Bucket) currentTick(now time.Time) int64 {
  return int64(now.Sub(tb.startTime) / tb.fillInterval)
}
```

桶中存储了初始时间和最后更新的tick

```go
type Bucket struct {
  clock Clock
  startTime time.Time
  capacity int64
  quantum int64
  fillInterval time.Duration
  mu sync.Mutex
  availableTokens int64
  latestTick int64
}
```

每次调用令牌桶暴露的接口,都会更新latestTick

```go
// 更新剩余令牌数,也是令牌桶唯一需要计算的地方
// 入参为currentTick(now()): 当前时间对应的tick
func (tb *Bucket) adjustavailableTokens(tick int64) {
  lastTick := tb.latestTick
  tb.latestTick = tick
  // 如果桶满了,就不用更新令牌剩余数了
  if tb.availableTokens >= tb.capacity {
    return
  }
  // 更新剩余令牌数,如果溢出了,则设为桶容量
  tb.availableTokens += (tick - lastTick) * tb.quantum
  if tb.availableTokens > tb.capacity {
    tb.availableTokens = tb.capacity
  }
  return
}
```

任何问题,只要转换成数学问题后,性能就提升了一大截,
这里并没有在每个间隔时间到了去更新桶,而是在每次接口调用时去更新,
通过tick差乘以quantum得到令牌增量.

> 令牌桶的核心实现后,就是丰富对外提供的功能

```go
// 桶容量
func (tb *Bucket) Capacity() int64 {
  return tb.capacity
}
// 每秒添加的令牌数,可以映射为业务中的QPS
func (tb *Bucket) Rate() float64 {
  return 1e9 * float64(tb.quantum) / float64(tb.fillInterval)
}
```

1e9/time.Duration,源自:

```go
type Duration int64
const (
  Nanosecond  Duration = 1
  Microsecond          = 1000 * Nanosecond
  Millisecond          = 1000 * Microsecond
  Second               = 1000 * Millisecond
  Minute               = 60 * Second
  Hour                 = 60 * Minute
)
```

从上面可以看出,1秒就是1e9,Duration是基于int64的新类型定义.

查询剩余令牌数:

```go
func (tb *Bucket) Available() int64 {
  return tb.available(tb.clock.Now())
}
func (tb *Bucket) available(now time.Time) int64 {
  tb.mu.Lock()
  defer tb.mu.Unlock()
  // currentTick是计算当前时间对应的tick
  // adjustavailableTokens是计算剩余令牌数
  tb.adjustavailableTokens(tb.currentTick(now))
  return tb.availableTokens
}
```

消耗指定数量的令牌,如果令牌不足,就只消耗掉所有剩余令牌

```go
func (tb *Bucket) TakeAvailable(count int64) int64 {
  tb.mu.Lock()
  defer tb.mu.Unlock()
  return tb.takeAvailable(tb.clock.Now(), count)
}
func (tb *Bucket) takeAvailable(now time.Time, count int64) int64 {
  if count <= 0 {
    return 0
  }
  tb.adjustavailableTokens(tb.currentTick(now))
  if tb.availableTokens <= 0 {
    return 0
  }
  if count > tb.availableTokens {
    count = tb.availableTokens
  }
  tb.availableTokens -= count
  return count
}
```

还有另一种消耗令牌的方式:消耗指定数量的令牌,外加一个等待时长,
如果在等待时长内还不能获取指定令牌数,就直接返回失败,并不消耗令牌,
供业务做进一步的业务逻辑处理;`如果在等待时长内能获取到指定令牌数,
则直接将剩余令牌数刷为负`.

```go
func (tb *Bucket) TakeMaxDuration(count int64, 
  maxWait time.Duration) (time.Duration, bool) {

  tb.mu.Lock()
  defer tb.mu.Unlock()
  return tb.take(tb.clock.Now(), count, maxWait)
}
func (tb *Bucket) take(now time.Time, count int64,
  maxWait time.Duration) (time.Duration, bool) {

  // 参数和剩余令牌数检查
  if count <= 0 {
    return 0, true
  }
  tick := tb.currentTick(now)
  tb.adjustavailableTokens(tick)
  avail := tb.availableTokens - count
  if avail >= 0 {
    // 如果剩余令牌数够,直接消耗掉
    tb.availableTokens = avail
    return 0, true
  }

  // 令牌数不足,先计算满足令牌数要多上时间
  endTick := tick + (-avail+tb.quantum-1)/tb.quantum
  endTime := tb.startTime.Add(time.Duration(endTick) * tb.fillInterval)
  waitTime := endTime.Sub(now)
  if waitTime > maxWait {
    // 如果超过了最大等待时间,直接返回失败
    return 0, false
  }
  // 在等待期限内,可以拿够令牌数
  // 剩余令牌数为负了,这是消耗了未来令牌数,必须用在特定场景
  tb.availableTokens = avail
  return waitTime, true
}
```

包中还提供了一个更加常用的函数Take,其最大等待时间为超大值,如果令牌不足则必
定消耗未来令牌数:

```go
const infinityDuration time.Duration = 0x7fffffffffffffff
```

Take算是TakeMaxDuration的一种特殊情况.

最后一个功能,基于TakeMaxDuration,如果消耗了未来的令牌数,则阻塞消耗掉tick.

```go
func (tb *Bucket) Wait(count int64) {
  if d := tb.Take(count); d > 0 {
    tb.clock.Sleep(d)
  }
}
```

剩余代码就是不同形式的构造函数,就不具体分析了.

### 测试代码分析

## ratelimit中令人眼前一亮的实现
