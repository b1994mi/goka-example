package main

import (
	"log"
	"net/http"

	"github.com/uptrace/bunrouter"
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

	h := handler.NewHandler()

	r.GET("/wallet", h.GetWalletHandler)

	port := ":5000"
	log.Printf("running on port %v", port)
	log.Println(http.ListenAndServe(port, r))
}
