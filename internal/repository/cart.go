package repository

import (
	"context"
	"errors"
)

// map[userID]map[skuID]count
type CartStorage = map[uint64]map[int64]uint16

type InMemoryCartRepository struct {
	cartStorage CartStorage
}

func NewRepository(cap int) *InMemoryCartRepository {
	return &InMemoryCartRepository{cartStorage: make(CartStorage, cap)}
}

func (r *InMemoryCartRepository) AddToCart(_ context.Context, skuID int64, userID uint64, count uint16) error {
	userCart, ok := r.cartStorage[userID]
	if !ok {
		userCart = make(map[int64]uint16)
	}

	if _, ok := userCart[skuID]; ok {
		userCart[skuID] += count
	} else {
		userCart[skuID] = count
	}

	r.cartStorage[userID] = userCart
	return nil
}

func (r *InMemoryCartRepository) RemoveFromCart(_ context.Context, skuID int64, userID uint64) error {
	userCart, ok := r.cartStorage[userID]
	if !ok {
		userCart = make(map[int64]uint16)
	}

	delete(userCart, skuID)

	return nil
}

func (r *InMemoryCartRepository) ClearCart(_ context.Context, userID uint64) error {
	_, ok := r.cartStorage[userID]
	if !ok {
		return errors.New("user not found")
	} else {
		delete(r.cartStorage, userID)
	}

	return nil
}

func (r *InMemoryCartRepository) GetCart(_ context.Context, userID uint64) (map[int64]uint16, error) {
	_, ok := r.cartStorage[userID]
	if !ok {
		return nil, nil
	}

	return r.cartStorage[userID], nil
}
