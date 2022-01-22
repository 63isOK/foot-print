# 多错误库

本意是将多个错误作为一个错误返回.

## hashicorp/go-multierror

兼容标准的errors库,支持As/Is/Unwrap等方法,特点是调用方如果知道函数返回多个错误,
则可以通过Unwrap来访问具体的错误.

这个库在结构上做到了优雅,"扁平化/格式化/分组/前缀/排序"都是单独的文件实现,
每个文件都是实现独立的功能.

其次go-multierror的递归写法是这样的:

```golang
func foo(){
  switch{
  case x:
    foo()
  default:
  }
}
```
