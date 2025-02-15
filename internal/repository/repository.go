package repository

import (
	"context"
	"errors"
	"github.com/vestamart/homework/internal/domain"
)

type Storage = map[int64][]domain.UserCart

type Repository struct {
	storage Storage
}

func NewRepository(cap int) *Repository {
	return &Repository{storage: make(Storage, cap)}
}

func (r *Repository) AddToCart(_ context.Context, item domain.Item, userId int64) (domain.Item, error) {
	if item.Count < 1 {
		return item, errors.New("invalid cart count")
	}

	r.storage[userId] = append(r.storage[userId])

	return item, nil
}

//func (r *Repository) RemoveFromCart(_ context.Context, cart domain.Item) (domain.UserCart, error) {
//
//}
//
//func (r *Repository) ClearCart(_ context.Context) (domain.UserCart, error) {
//
//}
//
//func (r *Repository) GetCart(_ context.Context, cart domain.Item) ([]domain.UserCart, error) {
//
//}
