package repository

import (
	"context"
	"errors"
)

type CartStorage = map[uint64]map[int64]uint16

type InMemoryRepository struct {
	cartStorage CartStorage
}

func NewRepository(cap int) *InMemoryRepository {
	return &InMemoryRepository{cartStorage: make(CartStorage, cap)}
}

func (r *InMemoryRepository) AddToCart(_ context.Context, skuID int64, userID uint64, count uint16) error {
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

func (r *InMemoryRepository) RemoveFromCart(_ context.Context, skuID int64, userID uint64) error {
	userCart, ok := r.cartStorage[userID]
	if !ok {
		userCart = make(map[int64]uint16)
	}

	delete(userCart, skuID)

	return nil
}

func (r *InMemoryRepository) ClearCart(_ context.Context, userID uint64) error {
	_, ok := r.cartStorage[userID]
	if !ok {
		return errors.New("user not found")
	} else {
		delete(r.cartStorage, userID)
	}

	return nil
}

func (r *InMemoryRepository) GetCart(_ context.Context, userID uint64) (map[int64]uint16, error) {
	_, ok := r.cartStorage[userID]
	if !ok {
		return nil, nil
	}

	return r.cartStorage[userID], nil
}
