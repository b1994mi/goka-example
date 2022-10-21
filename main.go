package main

import (
	"github.com/gin-gonic/gin"
	"goka-example/handler"
)

func main() {
	// routes
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	h := handler.NewHandler()

	r.GET("/", h.GetWalletHandler)

	r.Run(":5000")
}
