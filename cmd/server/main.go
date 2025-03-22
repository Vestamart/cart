package main

import (
	"github.com/vestamart/cart/internal/app/cart"
	"github.com/vestamart/cart/internal/client"
	"github.com/vestamart/cart/internal/config"
	"github.com/vestamart/cart/internal/delivery"
	"github.com/vestamart/cart/internal/mw"
	"github.com/vestamart/cart/internal/repository"
	"github.com/vestamart/loms/pkg/api/loms/v1"
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

	connLOMS, err := grpc.NewClient("localhost:"+cfg.LOMSServer.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	loggedMux := mw.LoggerHTTP(mux)

	log.Print("Server running on port: " + cfg.CartServer.Port)
	if err = http.ListenAndServe(":"+cfg.CartServer.Port, loggedMux); err != nil {
		log.Fatal(err)
	}
}
