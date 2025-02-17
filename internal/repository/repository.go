package repository

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/vestamart/homework/internal/delivery"
	"github.com/vestamart/homework/internal/domain"
)

type CartStorage = map[uint64]domain.UserCart

type Repository struct {
	cartStorage CartStorage
}

type GetCartRequest struct {
	Items      []domain.Item `json:"items"`
	TotalPrice uint32        `json:"total_price"`
}

func NewRepository(cap int) *Repository {
	return &Repository{cartStorage: make(CartStorage, cap)}
}

func (r *Repository) AddToCart(_ context.Context, skuID int64, userID uint64, count uint16) (*domain.UserCart, error) {
	cart, ok := r.cartStorage[userID]
	if !ok {
		cart = domain.UserCart{}
	}
	var itemFound bool
	clientRequest, err := delivery.GetProductHandler(skuID)
	if err != nil {
		return nil, err
	}

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
			Name:  clientRequest.Name,
			Count: count,
			Price: clientRequest.Price,
		})
	}
	r.cartStorage[userID] = cart

	return &cart, nil
}

func (r *Repository) RemoveFromCart(_ context.Context, skuID int64, userID uint64) (*domain.UserCart, error) {
	cart, ok := r.cartStorage[userID]
	if !ok {
		cart = domain.UserCart{}
	}

	var newItem []domain.Item
	//var deletedItem uint32
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

func (r *Repository) ClearCart(_ context.Context, userID uint64) (*domain.UserCart, error) {
	_, ok := r.cartStorage[userID]
	if !ok {
		return nil, errors.New("user not found")
	} else {
		delete(r.cartStorage, userID)
	}

	return nil, nil
}

func (r *Repository) GetCart(_ context.Context, userId uint64) ([]byte, error) {
	if _, ok := r.cartStorage[userId]; !ok {
		return nil, nil
	}
	var totalPrice uint32
	for _, item := range r.cartStorage[userId].Items {
		totalPrice += item.Price * uint32(item.Count)
	}
	getCartRequest := GetCartRequest{
		Items:      r.cartStorage[userId].Items,
		TotalPrice: totalPrice,
	}
	if jsonData, err := json.Marshal(getCartRequest); err != nil {
		return nil, errors.New("failed to marshal cart data")
	} else {
		return jsonData, nil
	}
}
