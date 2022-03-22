package main

import (
	"errors"
	"log"
	"net/rpc"
	"os/exec"
	"time"

	"github.com/63isOK/foot-print/plugin/rpc/hello"
	"github.com/hashicorp/go-plugin"
)

// NOTE: step5: 插件使用者调用插件

func main() {
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "BASIC_PLUGIN",
			MagicCookieValue: "hello",
		},
		Plugins: map[string]plugin.Plugin{
			"hello": &hello.Plugin{},
			"hi":    &hello.Plugin{},
		},
		Cmd: exec.Command("./plugin/hello"),
	})
	defer client.Kill()

	rpcClient, err := client.Client()
	if err != nil {
		log.Fatal(err)
	}
	raw, err := rpcClient.Dispense("hello")
	if err != nil {
		log.Fatal(err)
	}
	helloObject := raw.(hello.Greeter)

	hi, err := rpcClient.Dispense("hi")
	if err != nil {
		log.Fatal(err)
	}
	hiObject, ok := hi.(hello.Greeter)
	if !ok {
		log.Fatal("unexpected type from Dispense")
	}

	for {
		time.Sleep(time.Second)
		{
			println(helloObject.Greet())
			println(helloObject.GreetAgain())
		}

		{
			str, err := hiObject.Greet()
			if errors.Is(err, rpc.ErrShutdown) {
				println("plugin shutdown")
				rpcClient, err = client.Restart()
				if err != nil {
					println(err.Error())
				} else {
					hi, err = rpcClient.Dispense("hi")
					if err != nil {
						println(err.Error())
					}
					if hiObject, ok = hi.(hello.Greeter); !ok {
						println("unexpected type from Dispense")
						continue
					}
					str, err = hiObject.Greet()
					if err != nil {
						println("restart failed", err.Error())
					}
				}
			}
			println(str)
			println(hiObject.GreetAgain())
		}
	}
}
