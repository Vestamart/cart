package main

import (
	"github.com/vestamart/homework/internal/app/cart"
	"github.com/vestamart/homework/internal/client"
	"github.com/vestamart/homework/internal/config"
	"github.com/vestamart/homework/internal/delivery"
	"github.com/vestamart/homework/internal/repository"
	"github.com/vestamart/homework/pkg/api/loms/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
)

func main() {
	log.Println("App started")

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	clientProduct := client.NewClient(cfg.ProductClient.URL, cfg.ProductClient.Token)

	connLOMS, err := grpc.NewClient("localhost"+cfg.LOMSServer.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer connLOMS.Close()

	lomsClient := loms.NewLomsClient(connLOMS)

	repo := repository.NewRepository(100)
	service := cart.NewCartService(repo, clientProduct, lomsClient)
	server := delivery.NewServer(*service)

	router := delivery.NewRouter(server)
	mux := http.NewServeMux()
	router.SetupRoutes(mux)

	log.Print("Server running on port" + cfg.CartServer.Port)
	if err = http.ListenAndServe(cfg.CartServer.Port, mux); err != nil {
		log.Fatal(err)
	}
}
