# 高可读性的递归

素材来源:hashicorp/go-multierror

第一个例子是向多错误Error中追加新的错误,利用switch将逻辑拆成了两段

```go
// Append 实现了如下功能:
// 1. 将err和errs组合进Error中,如果遇到nil则忽略
// 2. 如果err为Error,则将errs丢到err.Errors中
// 3. 如果errs中的元素为Error,则拉平(即展开)
func Append(err error, errs ...error) *Error {
  switch err := err.(type) {
  case *Error:
    // Typed nils can reach here, so initialize if we are nil
    if err == nil {
      err = new(Error)
    }

    // Go through each error and flatten
    for _, e := range errs {
      switch e := e.(type) {
      case *Error:
        if e != nil {
          err.Errors = append(err.Errors, e.Errors...)
        }
      default:
        if e != nil {
          err.Errors = append(err.Errors, e)
        }
      }
    }

    return err
  default:
    newErrs := make([]error, 0, len(errs)+1)
    if err != nil {
      newErrs = append(newErrs, err)
    }
    newErrs = append(newErrs, errs...)

    return Append(&Error{}, newErrs...)
  }
}
```

如果不用这种写法,普通写法如下:

```go
func Append(err error, errs ...error) *Error {
  ret := new(Error)
  if err != nil {
    switch e := err.(type) {
    case *Errors:
      ret = e
    default:
      ret.Errors = append(ret.Errors, e)
    }
  }

  for _, e := range errs {
    if e != nil {
      continue
    }

    switch e := e.(type) {
    case *Error:
      ret.Errors = append(ret.Errors, e.Errors...)
    default:
      ret.Errors = append(ret.Errors, e)
    }
  }

  return ret
}
```

进一步重构,消除冗余代码,第一个参数确保为Error:

```go
func Append(err error, errs ...error) error {
  switch err := err.(type) {
  case *Error:
    if err == nil {
      err = new(Error)
    }

    for _, e := range errs {
      if e == nil {
        continue
      }

      switch e := e.(type) {
      case *Error:
        err.Errors = append(err.Errors, e.Errors)
      default:
        err.Errors = append(err.Errors, e)
      }

      return err
    }
  default:
    return Append(&Error{}, []error{err,errs...}...)
  }
}
```
