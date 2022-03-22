package hello

import "net/rpc"

// NOTE: step2: 基于rpc/grpc实现对接口的封装
// 客户端代码,封装了对接口的调用
// 这层封装会在Plugin.Client实现中使用
// 这里套多层的目的是为了简化插件的用法

// Client is plugin client
type Client struct {
	client *rpc.Client
}

func (c *Client) Greet() (resp string, err error) {
	err = c.client.Call("Plugin.Greet", new(interface{}), &resp)
	return
}

func (c *Client) GreetAgain() string {
	var resp string
	err := c.client.Call("Plugin.GreetAgain", new(interface{}), &resp)
	if err != nil {
		return err.Error()
	}
	return resp
}
