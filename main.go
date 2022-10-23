package main

import (
	"log"
	"net/http"

	"github.com/uptrace/bunrouter"
	"goka-example/config"
	"goka-example/handler"
)

func main() {
	// routes
	r := bunrouter.New()
	r.GET("/", func(w http.ResponseWriter, req bunrouter.Request) error {
		bunrouter.JSON(w, bunrouter.H{
			"message": "pong",
		})
		return nil
	})

	// goka processor/consumer
	config.InitGokaProcessor()

	// goka view
	atView := config.InitAboveThresholdView()
	bView := config.InitBalanceView()

	// routes with handlers
	h := handler.NewHandler(atView, bView)
	r.GET("/wallet", h.GetWalletHandler)
	r.POST("/wallet", h.PostWalletHandler)

	port := ":5000"
	log.Printf("running on port %v", port)
	log.Println(http.ListenAndServe(port, r))
}
