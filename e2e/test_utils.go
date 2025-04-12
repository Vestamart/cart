package e2e

//
//import (
//	"context"
//	"github.com/vestamart/homework/internal/app/cart"
//	"github.com/vestamart/homework/pkg/api/loms/v1"
//	"net/http"
//
//	"github.com/vestamart/homework/internal/delivery"
//	"github.com/vestamart/homework/internal/domain"
//	"github.com/vestamart/homework/internal/repository"
//)
//
//type mockProductService struct{}
//
//func (m *mockProductService) ExistItem(_ context.Context, _ int64) error {
//	return nil
//}
//
//func (m *mockProductService) GetProduct(_ context.Context, _ int64) (*domain.ProductServiceResponse, error) {
//	return &domain.ProductServiceResponse{Name: "Test Product", Price: 100}, nil
//}
//
//func SetupTestServer() *http.Server {
//	repo := repository.NewRepository(10)
//	productService := &mockProductService{}
//	var lomsClient loms.LomsClient
//	cartService := cart.NewCartService(repo, productService, lomsClient)
//	server := delivery.NewServer(*cartService)
//	router := delivery.NewRouter(server)
//	mux := http.NewServeMux()
//	router.SetupRoutes(mux)
//	return &http.Server{Handler: mux}
//}
