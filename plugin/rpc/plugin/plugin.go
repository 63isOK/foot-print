package main

import (
	"github.com/63isOK/foot-print/plugin/rpc/hello"

	"github.com/hashicorp/go-plugin"
)

// NOTE: step4: 插件作者编写插件

func main() {
	plugins := map[string]plugin.Plugin{
		"hello": &hello.Plugin{Impl: &helloImpl{}},
		"hi":    &hello.Plugin{Impl: &hiImpl{}},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "BASIC_PLUGIN",
			MagicCookieValue: "hello",
		},
		Plugins: plugins,
	})
}

type helloImpl struct{}

func (i *helloImpl) Greet() (string, error) {
	return "hello", nil
}

func (i *helloImpl) GreetAgain() string {
	return "hello, again"
}

type hiImpl struct{}

func (i *hiImpl) Greet() (string, error) {
	return "hi", nil
}

func (i *hiImpl) GreetAgain() string {
	return "hi, again"
}
