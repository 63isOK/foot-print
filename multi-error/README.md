# 多错误库

本意是将多个错误作为一个错误返回.

## hashicorp/go-multierror

兼容标准的errors库,支持As/Is/Unwrap等方法,特点是调用方如果知道函数返回多个错误,
则可以通过Unwrap来访问具体的错误.

