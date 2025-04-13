package main

import (
	"context"
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
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Println("App started")

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	clientProduct := client.NewClient(cfg.ProductClient.URL, cfg.ProductClient.Token)

	connLOMS, err := grpc.Dial("loms-service:"+cfg.LOMSServer.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
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

	srv := &http.Server{
		Addr:    ":" + cfg.CartServer.Port,
		Handler: loggedMux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Server running on port: %s", cfg.CartServer.Port)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	<-stop
	log.Println("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Println("Server gracefully stopped")
}
