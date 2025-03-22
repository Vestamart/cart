package cart

import (
	"context"
	"errors"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/vestamart/cart/internal/app/cart/mock"
	"github.com/vestamart/cart/internal/domain"
	"github.com/vestamart/cart/internal/localErr"
	"github.com/vestamart/loms/pkg/api/loms/v1"
	"testing"
)

func TestCartService_AddToCart(t *testing.T) {
	mc := minimock.NewController(t)
	repoMock := mock.NewCartRepositoryMock(mc)
	productMock := mock.NewProductServiceMock(mc)
	lomsMock := mock.NewLomsClientMock(mc)

	service := NewCartService(repoMock, productMock, lomsMock)

	tests := []struct {
		name         string
		skuID        int64
		userID       uint64
		count        uint16
		prepareMocks func()
		expectedErr  error
	}{
		{
			name:   "Valid input - success",
			skuID:  1003,
			userID: 456,
			count:  2,
			prepareMocks: func() {
				productMock.ExistItemMock.Return(nil)
				lomsMock.StocksInfoMock.Return(&loms.StocksInfoResponse{Count: 5}, nil)
				repoMock.AddToCartMock.Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:   "Not enough stock - error",
			skuID:  1003,
			userID: 456,
			count:  10,
			prepareMocks: func() {
				productMock.ExistItemMock.Return(nil)
				lomsMock.StocksInfoMock.Return(&loms.StocksInfoResponse{Count: 5}, nil)
			},
			expectedErr: localErr.ItemNotEnoughErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareMocks()
			err := service.AddToCart(context.Background(), tt.skuID, tt.userID, tt.count)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestCartService_RemoveFromCart(t *testing.T) {
	mc := minimock.NewController(t)

	repoMock := mock.NewCartRepositoryMock(mc)
	productMock := mock.NewProductServiceMock(mc)
	lomsMock := mock.NewLomsClientMock(mc)

	service := NewCartService(repoMock, productMock, lomsMock)

	tests := []struct {
		name         string
		skuID        int64
		userID       uint64
		prepareMocks func()
		expectedErr  error
	}{
		{
			name:   "Success case",
			skuID:  123,
			userID: 456,
			prepareMocks: func() {
				repoMock.RemoveFromCartMock.Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:   "Repository error",
			skuID:  789,
			userID: 101,
			prepareMocks: func() {
				repoMock.RemoveFromCartMock.Return(errors.New("failed to remove item"))
			},
			expectedErr: errors.New("failed to remove item"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareMocks()

			err := service.RemoveFromCart(context.Background(), tt.skuID, tt.userID)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCartService_ClearCart(t *testing.T) {
	mc := minimock.NewController(t)

	repoMock := mock.NewCartRepositoryMock(mc)
	productMock := mock.NewProductServiceMock(mc)
	lomsMock := mock.NewLomsClientMock(mc)

	service := NewCartService(repoMock, productMock, lomsMock)

	tests := []struct {
		name         string
		userID       uint64
		prepareMocks func()
		expectedErr  error
	}{
		{
			name:   "Valid userID - success",
			userID: 456,
			prepareMocks: func() {
				repoMock.ClearCartMock.Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:   "Repository error - failure",
			userID: 456,
			prepareMocks: func() {
				repoMock.ClearCartMock.Return(errors.New("database error"))
			},
			expectedErr: errors.New("database error"),
		},
		{
			name:   "Zero userID - success (no validation)",
			userID: 0,
			prepareMocks: func() {
				repoMock.ClearCartMock.Return(nil)
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareMocks()

			err := service.ClearCart(context.Background(), tt.userID)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCartService_GetCart(t *testing.T) {
	mc := minimock.NewController(t)

	repoMock := mock.NewCartRepositoryMock(mc)
	productMock := mock.NewProductServiceMock(mc)
	lomsMock := mock.NewLomsClientMock(mc)

	service := NewCartService(repoMock, productMock, lomsMock)

	tests := []struct {
		name         string
		userID       uint64
		prepareMocks func()
		expectedCart *domain.UserCart
		expectedErr  error
	}{
		{
			name:   "Valid cart with items - success",
			userID: 456,
			prepareMocks: func() {
				repoMock.GetCartMock.Return(map[int64]uint16{123: 2}, nil)
				productMock.GetProductMock.When(context.Background(), int64(123)).Then(&domain.ProductServiceResponse{
					Name:  "Test Product",
					Price: 100,
				}, nil)
			},
			expectedCart: &domain.UserCart{
				Items: []domain.CartItem{
					{
						Sku:   123,
						Name:  "Test Product",
						Count: 2,
						Price: 100,
					},
				},
				TotalPrice: 200,
			},
			expectedErr: nil,
		},
		{
			name:   "Empty cart - success",
			userID: 456,
			prepareMocks: func() {
				repoMock.GetCartMock.Return(map[int64]uint16{}, nil)
			},
			expectedCart: &domain.UserCart{
				Items:      nil,
				TotalPrice: 0,
			},
			expectedErr: nil,
		},
		{
			name:   "Repository error - failure",
			userID: 456,
			prepareMocks: func() {
				repoMock.GetCartMock.Return(nil, errors.New("database error"))
			},
			expectedCart: nil,
			expectedErr:  errors.New("database error"),
		},
		{
			name:   "Product service error - failure",
			userID: 4567,
			prepareMocks: func() {
				repoMock.GetCartMock.Return(map[int64]uint16{133: 2}, nil)
				productMock.GetProductMock.When(minimock.AnyContext, int64(133)).Then(nil, errors.New("product not found"))
			},
			expectedCart: nil,
			expectedErr:  errors.New("product not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareMocks()

			cart, err := service.GetCart(context.Background(), tt.userID)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
				assert.Nil(t, cart)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCart, cart)
			}
		})
	}
}

func TestCartService_CheckoutCart(t *testing.T) {
	mc := minimock.NewController(t)
	repoMock := mock.NewCartRepositoryMock(mc)
	productMock := mock.NewProductServiceMock(mc)
	lomsMock := mock.NewLomsClientMock(mc)

	service := NewCartService(repoMock, productMock, lomsMock)

	tests := []struct {
		name         string
		userID       uint64
		prepareMocks func()
		expectedID   int64
		expectedErr  error
	}{
		{
			name:   "Successful checkout",
			userID: 456,
			prepareMocks: func() {
				repoMock.GetCartMock.Return(map[int64]uint16{123: 2}, nil)
				productMock.GetProductMock.Return(&domain.ProductServiceResponse{Name: "Test Product", Price: 100}, nil)
				lomsMock.OrderCreateMock.Return(&loms.OrderCreateResponse{OrderId: 1}, nil)
				repoMock.ClearCartMock.Return(nil)
			},
			expectedID:  1,
			expectedErr: nil,
		},
		{
			name:   "Order creation fails",
			userID: 456,
			prepareMocks: func() {
				repoMock.GetCartMock.Return(map[int64]uint16{123: 2}, nil)
				productMock.GetProductMock.Return(&domain.ProductServiceResponse{Name: "Test Product", Price: 100}, nil)
				lomsMock.OrderCreateMock.Return(nil, errors.New("order creation failed"))
			},
			expectedID:  0,
			expectedErr: errors.New("order creation failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareMocks()
			orderID, err := service.CheckoutCart(context.Background(), tt.userID)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedID, orderID)
		})
	}
}
