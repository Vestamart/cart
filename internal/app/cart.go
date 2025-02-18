package app

import (
	"context"
	"errors"
	"github.com/vestamart/homework/internal/domain"
)

type CartRepository interface {
	AddToCart(_ context.Context, skuID int64, userID uint64, count uint16) (*domain.UserCart, error)
	RemoveFromCart(_ context.Context, skuID int64, userID uint64) (*domain.UserCart, error)
	ClearCart(_ context.Context, userID uint64) (*domain.UserCart, error)
	GetCart(_ context.Context, userID uint64) ([]byte, error)
}

type CartService struct {
	repository CartRepository
}

func NewCartService(repository CartRepository) *CartService {
	return &CartService{repository: repository}
}

func (s *CartService) AddToCart(ctx context.Context, skiID int64, userID uint64, count uint16) (*domain.UserCart, error) {
	if skiID < 1 || userID < 1 {
		return nil, errors.New("skiID or userID must be greater than 0")
	}

	return s.repository.AddToCart(ctx, skiID, userID, count)
}

func (s *CartService) RemoveFromCart(ctx context.Context, skiID int64, userID uint64) (*domain.UserCart, error) {
	return s.repository.RemoveFromCart(ctx, skiID, userID)
}

func (s *CartService) ClearCart(ctx context.Context, userID uint64) (*domain.UserCart, error) {
	return s.repository.ClearCart(ctx, userID)
}

func (s *CartService) GetCart(ctx context.Context, userID uint64) ([]byte, error) {
	return s.repository.GetCart(ctx, userID)
}
