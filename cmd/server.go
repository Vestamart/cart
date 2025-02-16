package main

import (
	"github.com/vestamart/homework/internal/app"
	"github.com/vestamart/homework/internal/delivery"
	"github.com/vestamart/homework/internal/repository"
	"log"
	"net/http"
)

func main() {
	log.Println("App started")

	repo := repository.NewRepository(100)
	service := app.NewCartService(repo)
	server := delivery.NewServer(service)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /user/{user_id}/cart/{sku_id}", server.AddToCartHandler)
	mux.HandleFunc("DELETE /user/{user_id}/cart/{sku_id}", server.RemoveFromCartHandler)
	mux.HandleFunc("DELETE /user/{user_id}/cart", server.ClearCartHandler)
	mux.HandleFunc("GET /user/{user_id}/cart", server.GetCartHandler)

	log.Print("Server running on port 8080")
	if err := http.ListenAndServe("127.0.0.1:8082", mux); err != nil {
		panic(err)
	}
}
