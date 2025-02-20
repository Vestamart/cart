package main

import (
	"github.com/vestamart/homework/internal/app"
	"github.com/vestamart/homework/internal/client"
	"github.com/vestamart/homework/internal/delivery"
	"github.com/vestamart/homework/internal/repository"
	"log"
	"net/http"
)

func main() {
	log.Println("App started")

	cfg, err := client.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	clientProduct := client.NewClient(cfg.Client.URL, cfg.Client.Token)

	repo := repository.NewRepository(100)
	service := app.NewCartService(repo, clientProduct)
	server := delivery.NewServer(*service)

	router := delivery.NewRouter(server)
	mux := http.NewServeMux()
	router.SetupRoutes(mux)

	log.Print("Server running on port" + cfg.Server.Port)
	if err = http.ListenAndServe(cfg.Server.Port, mux); err != nil {
		log.Fatal(err)
	}
}
