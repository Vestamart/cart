package repository

import (
	"context"
	"errors"
	"github.com/vestamart/homework/internal/domain"
)

type CartStorage = map[uint64]domain.UserCart

type InMemoryRepository struct {
	cartStorage CartStorage
}

func NewRepository(cap int) *InMemoryRepository {
	return &InMemoryRepository{cartStorage: make(CartStorage, cap)}
}

func (r *InMemoryRepository) AddToCart(_ context.Context, skuID int64, userID uint64, count uint16) (*domain.UserCart, error) {
	userCart, ok := r.cartStorage[userID]
	if !ok {
		userCart = domain.UserCart{
			Items: make(map[int64]uint16),
		}
	}

	if _, ok := userCart.Items[skuID]; ok {
		userCart.Items[skuID] += count
	} else {
		userCart.Items[skuID] = count
	}

	r.cartStorage[userID] = userCart
	return &userCart, nil
}

func (r *InMemoryRepository) RemoveFromCart(_ context.Context, skuID int64, userID uint64) (*domain.UserCart, error) {
	userCart, ok := r.cartStorage[userID]
	if !ok {
		userCart = domain.UserCart{}
	}

	delete(userCart.Items, skuID)

	return &userCart, nil
}

func (r *InMemoryRepository) ClearCart(_ context.Context, userID uint64) (*domain.UserCart, error) {
	_, ok := r.cartStorage[userID]
	if !ok {
		return nil, errors.New("user not found")
	} else {
		delete(r.cartStorage, userID)
	}

	return &domain.UserCart{}, nil
}

func (r *InMemoryRepository) GetCart(_ context.Context, userID uint64) (*domain.UserCart, error) {
	userCart, ok := r.cartStorage[userID]
	if !ok {
		return &domain.UserCart{}, nil
	}

	return &userCart, nil
}
