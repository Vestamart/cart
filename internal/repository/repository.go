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
	cart, ok := r.cartStorage[userID]
	if !ok {
		cart = domain.UserCart{}
	}
	var itemFound bool
	for i, item := range cart.Items {
		if item.SkuID == skuID {
			cart.Items[i].Count += count
			itemFound = true
			break
		}
	}

	if !itemFound {
		cart.Items = append(cart.Items, domain.Item{
			SkuID: skuID,
			Count: count,
		})
	}

	r.cartStorage[userID] = cart
	return &cart, nil
}

func (r *InMemoryRepository) RemoveFromCart(_ context.Context, skuID int64, userID uint64) (*domain.UserCart, error) {
	cart, ok := r.cartStorage[userID]
	if !ok {
		cart = domain.UserCart{}
	}

	var newItem []domain.Item
	for _, item := range cart.Items {
		if item.SkuID == skuID {
			continue
		}
		newItem = append(newItem, item)
	}
	cart.Items = newItem

	if len(cart.Items) == 0 {
		delete(r.cartStorage, userID)
	} else {
		r.cartStorage[userID] = cart
	}
	return &cart, nil
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
