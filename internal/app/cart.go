package app

import (
	"context"
	"errors"
	"github.com/vestamart/homework/internal/client"
	"github.com/vestamart/homework/internal/domain"
)

type CartRepository interface {
	AddToCart(ctx context.Context, skuID int64, userID uint64, count uint16) (*domain.UserCart, error)
	RemoveFromCart(ctx context.Context, skuID int64, userID uint64) (*domain.UserCart, error)
	ClearCart(ctx context.Context, userID uint64) (*domain.UserCart, error)
	GetCart(ctx context.Context, userID uint64) (*domain.UserCart, error)
}

type ProductService interface {
	ExistItem(ctx context.Context, sku int64) error
	GetProductHandler(ctx context.Context, sku int64) (*client.Response, error)
}

type CartService struct {
	repository     CartRepository
	productService ProductService
}

type GetCartItemResponse struct {
	Sku   int64  `json:"sku"`
	Name  string `json:"name"`
	Count uint16 `json:"count"`
	Price uint32 `json:"price"`
}
type GetCartResponse struct {
	Items      []GetCartItemResponse `json:"items"`
	TotalPrice uint32                `json:"total_price"`
}

func NewCartService(repository CartRepository, client ProductService) *CartService {
	return &CartService{repository: repository, productService: client}
}

func (s *CartService) AddToCart(ctx context.Context, skuID int64, userID uint64, count uint16) (*domain.UserCart, error) {
	if skuID < 1 || userID < 1 {
		return nil, errors.New("skuID or userID must be greater than 0")
	}
	if err := s.productService.ExistItem(ctx, skuID); err != nil {
		return nil, err
	}

	return s.repository.AddToCart(ctx, skuID, userID, count)
}

func (s *CartService) RemoveFromCart(ctx context.Context, skuID int64, userID uint64) (*domain.UserCart, error) {
	return s.repository.RemoveFromCart(ctx, skuID, userID)
}

func (s *CartService) ClearCart(ctx context.Context, userID uint64) (*domain.UserCart, error) {
	return s.repository.ClearCart(ctx, userID)
}

func (s *CartService) GetCart(ctx context.Context, userID uint64) (*GetCartResponse, error) {
	userCart, err := s.repository.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	var totalPrice uint32
	var cart GetCartResponse
	for sku, count := range userCart.Items {
		resp, err := s.productService.GetProductHandler(ctx, sku)
		if err != nil {
			return nil, err
		}
		totalPrice += uint32(count) * resp.Price
		cart.Items = append(cart.Items, GetCartItemResponse{
			Sku:   sku,
			Name:  resp.Name,
			Count: count,
			Price: resp.Price,
		})
	}
	cart.TotalPrice = totalPrice
	return &cart, nil

}
