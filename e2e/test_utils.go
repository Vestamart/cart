package e2e

import (
	"context"
	"net/http"

	"github.com/vestamart/homework/internal/app"
	"github.com/vestamart/homework/internal/delivery"
	"github.com/vestamart/homework/internal/domain"
	"github.com/vestamart/homework/internal/repository"
)

type mockProductService struct{}

func (m *mockProductService) ExistItem(_ context.Context, _ int64) error {
	return nil
}

func (m *mockProductService) GetProduct(_ context.Context, _ int64) (*domain.ProductServiceResponse, error) {
	return &domain.ProductServiceResponse{Name: "Test Product", Price: 100}, nil
}

func SetupTestServer() *http.Server {
	repo := repository.NewRepository(10)
	productService := &mockProductService{}
	cartService := app.NewCartService(repo, productService)
	server := delivery.NewServer(*cartService)
	router := delivery.NewRouter(server)
	mux := http.NewServeMux()
	router.SetupRoutes(mux)
	return &http.Server{Handler: mux}
}
