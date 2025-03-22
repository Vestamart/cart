package cart

import (
	"context"
	"errors"
	"github.com/vestamart/cart/internal/domain"
	"github.com/vestamart/cart/internal/localErr"
	"github.com/vestamart/loms/pkg/api/loms/v1"
)

//go:generate minimock -i github.com/vestamart/cart/internal/app/cart.Repository -o ./mock/repository_mock.go -n CartRepositoryMock -p mock
type Repository interface {
	AddToCart(_ context.Context, skuID int64, userID uint64, count uint16) error
	RemoveFromCart(_ context.Context, skuID int64, userID uint64) error
	ClearCart(_ context.Context, userID uint64) error
	GetCart(_ context.Context, userID uint64) (map[int64]uint16, error)
}

//go:generate minimock -i github.com/vestamart/cart/internal/app/cart.ProductService -o ./mock/product_service_mock.go -n ProductServiceMock -p mock
type ProductService interface {
	ExistItem(ctx context.Context, sku int64) error
	GetProduct(ctx context.Context, sku int64) (*domain.ProductServiceResponse, error)
}

//go:generate minimock -i github.com/vestamart/loms/pkg/api/loms/v1.LomsClient -o ./mock/loms_client_mock.go -n LomsClientMock -p mock
type Service struct {
	repository     Repository
	productService ProductService
	lomsService    loms.LomsClient
}

func NewCartService(repository Repository, client ProductService, loms loms.LomsClient) *Service {
	return &Service{repository: repository, productService: client, lomsService: loms}
}

func (s *Service) AddToCart(ctx context.Context, skuID int64, userID uint64, count uint16) error {
	if skuID < 1 || userID < 1 {
		return errors.New("skuID or userID must be greater than 0")
	}
	if err := s.productService.ExistItem(ctx, skuID); err != nil {
		return err
	}

	v, err := s.lomsService.StocksInfo(ctx, &loms.StocksInfoRequest{Sku: uint32(skuID)})
	if err != nil {
		return err
	}
	if uint16(v.Count) <= count {
		return localErr.ItemNotEnoughErr
	}

	return s.repository.AddToCart(ctx, skuID, userID, count)
}

func (s *Service) RemoveFromCart(ctx context.Context, skuID int64, userID uint64) error {
	return s.repository.RemoveFromCart(ctx, skuID, userID)
}

func (s *Service) ClearCart(ctx context.Context, userID uint64) error {
	return s.repository.ClearCart(ctx, userID)
}

func (s *Service) GetCart(ctx context.Context, userID uint64) (*domain.UserCart, error) {
	userCart, err := s.repository.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	var totalPrice uint32
	var cart domain.UserCart

	for sku, count := range userCart {
		resp, err := s.productService.GetProduct(ctx, sku)
		if err != nil {
			return nil, err
		}
		totalPrice += uint32(count) * resp.Price
		cart.Items = append(cart.Items, domain.CartItem{
			Sku:   sku,
			Name:  resp.Name,
			Count: count,
			Price: resp.Price,
		})
	}
	cart.TotalPrice = totalPrice
	return &cart, nil
}

func (s *Service) CheckoutCart(ctx context.Context, userID uint64) (int64, error) {
	cart, err := s.GetCart(ctx, userID)
	if err != nil {
		return 0, err
	}

	var items []*loms.Item
	for _, item := range cart.Items {
		items = append(items, &loms.Item{
			Sku:   uint32(item.Sku),
			Count: uint32(item.Count),
		})
	}

	lomsRequest := &loms.OrderCreateRequest{
		User:  int64(userID),
		Items: items,
	}

	orderID, err := s.lomsService.OrderCreate(ctx, lomsRequest)
	if err != nil {
		return 0, err
	}

	err = s.ClearCart(ctx, userID)
	if err != nil {
		return 0, err
	}

	return orderID.GetOrderId(), nil
}
