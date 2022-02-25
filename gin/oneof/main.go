package main

import (
	"github.com/gin-gonic/gin"
)

// User 用户
type User struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Country string `json:"country" binding:"required,oneof=china taiwan"`
}

func main() {
	g := gin.Default()
	g.POST("/ping", func(c *gin.Context) {
		var user User
		err := c.ShouldBind(&user)
		if err != nil {
			c.JSON(200, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"user": user,
		})
	})
	g.Run(":8888")
}
