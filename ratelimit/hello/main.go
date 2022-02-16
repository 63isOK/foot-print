package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"gopkg.in/resty.v1"
)

func main() {
	go client()

	g := gin.Default()
	g.Use(RateLimit(time.Second, 2, 1))
	g.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	g.Run(":8080")
}

// RateLimit 限流中间件
func RateLimit(fillInterval time.Duration, cap, quantum int64) gin.HandlerFunc {
	bucket := ratelimit.NewBucketWithQuantum(fillInterval, cap, quantum)
	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) == 0 {
			c.String(http.StatusTooManyRequests, "rate limit ...")
			c.Abort()
		} else {
			c.Next()
		}
	}
}

func client() {
	time.Sleep(time.Second * 3)

	for i := 0; i < 20; i++ {
		resp, _ := resty.New().R().Get("http://localhost:8080/ping")
		fmt.Println(string(resp.Body()))
		time.Sleep(500 * time.Millisecond)
	}
}
