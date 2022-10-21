package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// routes
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run(":5000")
}
