package app

import (
	"context"
	"errors"
	"github.com/vestamart/homework/internal/client"
	"github.com/vestamart/homework/internal/domain"
)

type CartRepository interface {
	AddToCart(_ context.Context, skuID int64, userID uint64, count uint16) error
	RemoveFromCart(_ context.Context, skuID int64, userID uint64) error
	ClearCart(_ context.Context, userID uint64) error
	GetCart(_ context.Context, userID uint64) (map[int64]uint16, error)
}

type ProductService interface {
	ExistItem(ctx context.Context, sku int64) error
	GetProduct(ctx context.Context, sku int64) (*client.Response, error)
}

type CartService struct {
	repository     CartRepository
	productService ProductService
}

func NewCartService(repository CartRepository, client ProductService) *CartService {
	return &CartService{repository: repository, productService: client}
}

func (s *CartService) AddToCart(ctx context.Context, skuID int64, userID uint64, count uint16) error {
	if skuID < 1 || userID < 1 {
		return errors.New("skuID or userID must be greater than 0")
	}
	if err := s.productService.ExistItem(ctx, skuID); err != nil {
		return err
	}
	return s.repository.AddToCart(ctx, skuID, userID, count)
}

func (s *CartService) RemoveFromCart(ctx context.Context, skuID int64, userID uint64) error {
	return s.repository.RemoveFromCart(ctx, skuID, userID)
}

func (s *CartService) ClearCart(ctx context.Context, userID uint64) error {
	return s.repository.ClearCart(ctx, userID)
}

func (s *CartService) GetCart(ctx context.Context, userID uint64) (*domain.UserCart, error) {
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
