package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryRepository_AddToCart(t *testing.T) {
	tests := []struct {
		name         string
		skuID        int64
		userID       uint64
		count        uint16
		initialState CartStorage
		expectedCart map[int64]uint16
		expectedErr  error
	}{
		{
			name:         "Add to empty cart - success",
			skuID:        123,
			userID:       456,
			count:        2,
			initialState: CartStorage{},
			expectedCart: map[int64]uint16{123: 2},
			expectedErr:  nil,
		},
		{
			name:         "Add to existing cart - increase count",
			skuID:        123,
			userID:       456,
			count:        3,
			initialState: CartStorage{456: map[int64]uint16{123: 2}},
			expectedCart: map[int64]uint16{123: 5},
			expectedErr:  nil,
		},
		{
			name:         "Add new item to existing cart - success",
			skuID:        789,
			userID:       456,
			count:        1,
			initialState: CartStorage{456: map[int64]uint16{123: 2}},
			expectedCart: map[int64]uint16{123: 2, 789: 1},
			expectedErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewRepository(10)
			repo.cartStorage = tt.initialState

			err := repo.AddToCart(context.Background(), tt.skuID, tt.userID, tt.count)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCart, repo.cartStorage[tt.userID])
			}
		})
	}
}

func TestInMemoryRepository_RemoveFromCart(t *testing.T) {
	tests := []struct {
		name         string
		skuID        int64
		userID       uint64
		initialState CartStorage
		expectedCart map[int64]uint16
		expectedErr  error
	}{
		{
			name:         "Remove from existing cart - success",
			skuID:        123,
			userID:       456,
			initialState: CartStorage{456: map[int64]uint16{123: 2, 789: 1}},
			expectedCart: map[int64]uint16{789: 1},
			expectedErr:  nil,
		},
		{
			name:         "Remove from empty cart - success",
			skuID:        123,
			userID:       456,
			initialState: CartStorage{},
			expectedCart: map[int64]uint16{},
			expectedErr:  nil,
		},
		{
			name:         "Remove non-existent item - success",
			skuID:        999,
			userID:       456,
			initialState: CartStorage{456: map[int64]uint16{123: 2}},
			expectedCart: map[int64]uint16{123: 2},
			expectedErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewRepository(10)
			repo.cartStorage = tt.initialState

			err := repo.RemoveFromCart(context.Background(), tt.skuID, tt.userID)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				userCart, ok := repo.cartStorage[tt.userID]
				if len(tt.expectedCart) > 0 || ok {
					assert.Equal(t, tt.expectedCart, userCart)
				} else {
					assert.False(t, ok, "Cart should not exist for user")
				}
			}
		})
	}
}

func TestInMemoryRepository_ClearCart(t *testing.T) {
	tests := []struct {
		name         string
		userID       uint64
		initialState CartStorage
		expectedErr  error
	}{
		{
			name:         "Clear existing cart - success",
			userID:       456,
			initialState: CartStorage{456: map[int64]uint16{123: 2}},
			expectedErr:  nil,
		},
		{
			name:         "Clear non-existent cart - error",
			userID:       789,
			initialState: CartStorage{456: map[int64]uint16{123: 2}},
			expectedErr:  errors.New("user not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewRepository(10)
			repo.cartStorage = tt.initialState

			err := repo.ClearCart(context.Background(), tt.userID)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
				if _, ok := tt.initialState[tt.userID]; !ok {
					assert.Contains(t, repo.cartStorage, uint64(456), "Original cart should remain")
				}
			} else {
				assert.NoError(t, err)
				_, exists := repo.cartStorage[tt.userID]
				assert.False(t, exists, "Cart should be deleted")
			}
		})
	}
}

func TestInMemoryRepository_GetCart(t *testing.T) {
	tests := []struct {
		name         string
		userID       uint64
		initialState CartStorage
		expectedCart map[int64]uint16
		expectedErr  error
	}{
		{
			name:         "Get existing cart - success",
			userID:       456,
			initialState: CartStorage{456: map[int64]uint16{123: 2, 789: 1}},
			expectedCart: map[int64]uint16{123: 2, 789: 1},
			expectedErr:  nil,
		},
		{
			name:         "Get non-existent cart - nil",
			userID:       789,
			initialState: CartStorage{456: map[int64]uint16{123: 2}},
			expectedCart: nil,
			expectedErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewRepository(10)
			repo.cartStorage = tt.initialState

			cart, err := repo.GetCart(context.Background(), tt.userID)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCart, cart)
			}
		})
	}
}
