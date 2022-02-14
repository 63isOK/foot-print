# golang 静态检测

[golangci-lint](https://golangci-lint.run/)是一个lint聚合器.

特点:

1. 快,并行+基于go build缓存
2. 基于yaml格式,向cncf配置看齐
3. 集成度好,vscode/vim/github action都支持
4. 聚合,无需单独安装lint
5. 误报少

在项目目录运行`golangci-lint run`即可.

```shell
# 查看聚合了多少lint
golangci-lint help lint

# enable/E, disable/D, --enable-all, --disable-all 手动控制lint的开启或关闭
golangci-lint --enable-all -D errcheck
```
