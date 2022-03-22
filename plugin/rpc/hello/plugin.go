package hello

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// NOTE: step3: 屏蔽客户端对接口的使用,屏蔽服务端接口实现的写法
// 这步的写法更多的是统一对外的呈现,不管是客户端还是服务端

// Plugin 实现了插件的封装
type Plugin struct {
	Impl Greeter
}

func (p *Plugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &Server{Impl: p.Impl}, nil
}

func (p *Plugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &Client{client: c}, nil
}
