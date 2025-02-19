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

type CartService struct {
	repository CartRepository
	client     *client.Client
}

func NewCartService(repository CartRepository, client *client.Client) *CartService {
	return &CartService{repository: repository, client: client}
}

func (s *CartService) AddToCart(ctx context.Context, skuID int64, userID uint64, count uint16) (*domain.UserCart, error) {
	if skuID < 1 || userID < 1 {
		return nil, errors.New("skuID or userID must be greater than 0")
	}
	if ok, err := s.client.ExistItem(skuID); err != nil && ok != true {
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

func (s *CartService) GetCart(ctx context.Context, userID uint64) (*domain.UserCart, error) {
	cart, err := s.repository.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}
	var totalPrice uint32

	for i, item := range cart.Items {
		resp, err := s.client.GetProductHandler(item.SkuID)
		if err != nil {
			return nil, err
		}
		cart.Items[i].Price = resp.Price
		cart.Items[i].Name = resp.Name
		totalPrice += resp.Price * uint32(cart.Items[i].Count)
	}

	getCartResponse := domain.UserCart{
		Items:      cart.Items,
		TotalPrice: totalPrice,
	}

	return &getCartResponse, nil

}
