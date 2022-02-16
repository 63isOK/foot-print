# juju限流库

令牌桶的算法唯一需要计算的地方是从桶bucket中移除令牌.

一定时间间隔内,向桶中添加一定量的令牌

```go
type Bucket struct {
  clock Clock                 // 时钟,接口类型,可以进行
  startTime time.Time         // 初始时间
  capacity int64              // 容量
  quantum int64               // 每次向桶中添加的令牌数
  fillInterval time.Duration  // 添加令牌的间隔
  mu sync.Mutex               // 互斥量,守护下面两个变量
  availableTokens int64       // 剩余令牌数
  latestTick int64            // 最后一次tick
}
```

时钟接口,定义了获取当前时间和休眠,定义为接口是为了打桩,方便测试

```go
type Clock interface {
  Now() time.Time
  Sleep(d time.Duration)
}
```

构造令牌桶

```go
func NewBucketWithQuantumAndClock(fillInterval time.Duration,
  capacity, quantum int64, clock Clock) *Bucket {

  // 设置时钟
  if clock == nil {
    clock = realClock{}
  }

  // 入参校验
  if fillInterval <= 0 {
    panic("token bucket fill interval is not > 0")
  }
  if capacity <= 0 {
    panic("token bucket capacity is not > 0")
  }
  if quantum <= 0 {
    panic("token bucket quantum is not > 0")
  }

  // 构造
  return &Bucket{
    clock:           clock,
    startTime:       clock.Now(),
    latestTick:      0,
    fillInterval:    fillInterval,
    capacity:        capacity,
    quantum:         quantum,
    availableTokens: capacity,
  }
}

func NewBucketWithRateAndClock(rate float64, capacity int64, clock Clock)
  *Bucket {

  // Use the same bucket each time through the loop
  // to save allocations.
  tb := NewBucketWithQuantumAndClock(1, capacity, 1, clock)
  for quantum := int64(1); quantum < 1<<50; quantum = nextQuantum(quantum) {
    fillInterval := time.Duration(1e9 * float64(quantum) / rate)
    if fillInterval <= 0 {
      continue
    }
    tb.fillInterval = fillInterval
    tb.quantum = quantum
    if diff := math.Abs(tb.Rate() - rate); diff/rate <= rateMargin {
      return tb
    }
  }
  panic("cannot find suitable quantum for " +
    strconv.FormatFloat(rate, 'g', -1, 64))
}
```

Bucket.Take: 计算获取指定令牌所需时间,
Bucket.Wait: sleep,等待,直到能获取指定的令牌数,
Bucket.TakeMaxDuration: 计算指定时间内,获取指定令牌所需事件,
Bucket.WaitMaxDuration: 等待获取指定数量的令牌数,带时长,
Bucket.TakeAvailable: 获取桶内剩余令牌数,
Bucket.Available: 查看桶内剩余令牌数,

ratelimit.go源文件主要定义了桶,以及桶暴露的方法.
