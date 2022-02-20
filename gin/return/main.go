package main

import "github.com/gin-gonic/gin"

func main() {
	g := gin.Default()
	g.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	g.Run(":9999")
}
