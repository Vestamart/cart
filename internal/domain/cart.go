package domain

import (
	"errors"
	"fmt"
)

type UserCart struct {
	Items      []CartItem `json:"items"`
	TotalPrice uint32     `json:"total_price"`
}

type CartItem struct {
	Sku   int64  `json:"sku"`
	Name  string `json:"name"`
	Count uint16 `json:"count"`
	Price uint32 `json:"price"`
}

type ProductServiceResponse struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

// Validate проверяет корректность CartItem
func (ci *CartItem) Validate() error {
	if ci.Sku <= 0 {
		return errors.New("SKU must be positive")
	}
	if ci.Name == "" {
		return errors.New("name cannot be empty")
	}
	if ci.Count == 0 {
		return errors.New("count must be greater than 0")
	}
	if ci.Price == 0 {
		return errors.New("price must be greater than 0")
	}
	return nil
}

// Validate проверяет корректность UserCart
func (uc *UserCart) Validate() error {
	if len(uc.Items) == 0 {
		return errors.New("cart cannot be empty")
	}

	for i, item := range uc.Items {
		if err := item.Validate(); err != nil {
			return fmt.Errorf("item %d: %w", i, err)
		}
	}

	// Проверяем, что TotalPrice соответствует сумме цен товаров
	var calculatedTotal uint32
	for _, item := range uc.Items {
		calculatedTotal += item.Price * uint32(item.Count)
	}

	if calculatedTotal != uc.TotalPrice {
		return fmt.Errorf("total price mismatch: calculated %d, got %d", calculatedTotal, uc.TotalPrice)
	}

	return nil
}

// Validate проверяет корректность ProductServiceResponse
func (psr *ProductServiceResponse) Validate() error {
	if psr.Name == "" {
		return errors.New("product name cannot be empty")
	}
	if psr.Price == 0 {
		return errors.New("product price must be greater than 0")
	}
	return nil
}
