# 尤达

尤达讲究的是非标准

```txt
// 很多语言都推荐这种写法,防止逻辑判断误写为赋值 v=1
if 1 == v {
}
```

但在go中,天然地就避免了这种错误,因为 if v =1 {} 编译就会报错,
所以go中不推荐使用尤达写法.