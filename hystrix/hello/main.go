package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/afex/hystrix-go/hystrix"

	"github.com/gin-gonic/gin"
	"gopkg.in/resty.v1"
)

var start = time.Now()

func main() {
	// 开启dashboard
	// hystrixStreamHandler := hystrix.NewStreamHandler()
	// hystrixStreamHandler.Start()
	// go http.ListenAndServe(net.JoinHostPort("", "8081"), hystrixStreamHandler)

	go server()

	hystrix.ConfigureCommand("test", hystrix.CommandConfig{
		Timeout:                10,
		MaxConcurrentRequests:  100,
		RequestVolumeThreshold: 10,
		SleepWindow:            500,
		ErrorPercentThreshold:  20,
	})

	for i := 0; i < 20; i++ {
		_ = hystrix.Do("test", func() error {
			resp, _ := resty.New().R().Get("http://localhost:8080/ping")
			if resp.IsError() {
				return fmt.Errorf("error code: %s", resp.Status())
			}
			return nil
		}, func(e error) error {
			fmt.Println("fallback err: ", e)
			return e
		})
		time.Sleep(100 * time.Millisecond)
	}
}

func server() {
	e := gin.Default()
	e.GET("/ping", func(ctx *gin.Context) {
		if time.Since(start) < 201*time.Millisecond {
			ctx.String(http.StatusInternalServerError, "pong")
			return
		}

		ctx.String(http.StatusOK, "pong")
	})

	e.Run(":8080")
}
