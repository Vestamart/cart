package repository

import (
	"context"
	"errors"
	"github.com/vestamart/homework/internal/app/cart"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryRepository_AddToCart(t *testing.T) {
	tests := []struct {
		name         string
		skuID        int64
		userID       uint64
		count        uint16
		prepareCart  func(ctx context.Context, repo cart.Repository)
		expectedCart map[int64]uint16
		expectedErr  error
	}{
		{
			name:         "Add to empty cart - success",
			skuID:        123,
			userID:       456,
			count:        2,
			prepareCart:  nil,
			expectedCart: map[int64]uint16{123: 2},
			expectedErr:  nil,
		},
		{
			name:   "Add to existing cart - increase count",
			skuID:  123,
			userID: 456,
			count:  3,
			prepareCart: func(ctx context.Context, repo cart.Repository) {
				_ = repo.AddToCart(ctx, 123, 456, 2)
			},
			expectedCart: map[int64]uint16{123: 5},
			expectedErr:  nil,
		},
		{
			name:   "Add new item to existing cart - success",
			skuID:  789,
			userID: 456,
			count:  1,
			prepareCart: func(ctx context.Context, repo cart.Repository) {
				_ = repo.AddToCart(ctx, 123, 456, 2)
			},
			expectedCart: map[int64]uint16{123: 2, 789: 1},
			expectedErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewRepository(10)
			ctx := context.Background()

			if tt.prepareCart != nil {
				tt.prepareCart(ctx, repo)
			}

			err := repo.AddToCart(ctx, tt.skuID, tt.userID, tt.count)
			assert.NoError(t, err)

			cart, err := repo.GetCart(ctx, tt.userID)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCart, cart)
		})
	}
}

func TestInMemoryRepository_RemoveFromCart(t *testing.T) {
	tests := []struct {
		name         string
		skuID        int64
		userID       uint64
		prepareCart  func(ctx context.Context, repo cart.Repository)
		expectedCart map[int64]uint16
		expectedErr  error
	}{
		{
			name:   "Remove from existing cart - success",
			skuID:  123,
			userID: 456,
			prepareCart: func(ctx context.Context, repo cart.Repository) {
				_ = repo.AddToCart(ctx, 123, 456, 2)
				_ = repo.AddToCart(ctx, 789, 456, 1)
			},
			expectedCart: map[int64]uint16{789: 1},
			expectedErr:  nil,
		},
		{
			name:         "Remove from empty cart - success",
			skuID:        123,
			userID:       456,
			prepareCart:  nil,
			expectedCart: nil,
			expectedErr:  nil,
		},
		{
			name:   "Remove non-existent item - success",
			skuID:  999,
			userID: 456,
			prepareCart: func(ctx context.Context, repo cart.Repository) {
				_ = repo.AddToCart(ctx, 123, 456, 2)
			},
			expectedCart: map[int64]uint16{123: 2},
			expectedErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewRepository(10)
			ctx := context.Background()

			if tt.prepareCart != nil {
				tt.prepareCart(ctx, repo)
			}

			err := repo.RemoveFromCart(ctx, tt.skuID, tt.userID)
			assert.NoError(t, err)

			cart, err := repo.GetCart(ctx, tt.userID)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCart, cart)
		})
	}
}

func TestInMemoryRepository_ClearCart(t *testing.T) {
	tests := []struct {
		name        string
		userID      uint64
		prepareCart func(ctx context.Context, repo cart.Repository)
		expectedErr error
	}{
		{
			name:   "Clear existing cart - success",
			userID: 456,
			prepareCart: func(ctx context.Context, repo cart.Repository) {
				_ = repo.AddToCart(ctx, 123, 456, 2)
			},
			expectedErr: nil,
		},
		{
			name:        "Clear non-existent cart - error",
			userID:      789,
			prepareCart: nil,
			expectedErr: errors.New("user not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewRepository(10)
			ctx := context.Background()

			if tt.prepareCart != nil {
				tt.prepareCart(ctx, repo)
			}

			err := repo.ClearCart(ctx, tt.userID)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				cart, err := repo.GetCart(ctx, tt.userID)
				assert.NoError(t, err)
				assert.Nil(t, cart)
			}
		})
	}
}

func TestInMemoryRepository_GetCart(t *testing.T) {
	tests := []struct {
		name         string
		userID       uint64
		prepareCart  func(ctx context.Context, repo cart.Repository)
		expectedCart map[int64]uint16
		expectedErr  error
	}{
		{
			name:   "Get existing cart - success",
			userID: 456,
			prepareCart: func(ctx context.Context, repo cart.Repository) {
				_ = repo.AddToCart(ctx, 123, 456, 2)
				_ = repo.AddToCart(ctx, 789, 456, 1)
			},
			expectedCart: map[int64]uint16{123: 2, 789: 1},
			expectedErr:  nil,
		},
		{
			name:         "Get non-existent cart - nil",
			userID:       789,
			prepareCart:  nil,
			expectedCart: nil,
			expectedErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewRepository(10)
			ctx := context.Background()

			if tt.prepareCart != nil {
				tt.prepareCart(ctx, repo)
			}

			cart, err := repo.GetCart(ctx, tt.userID)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCart, cart)
		})
	}
}

// Бенчмарки
func BenchmarkHandler_AddToCart(b *testing.B) {
	repo := NewRepository(100)
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		_ = repo.AddToCart(ctx, int64(123), uint64(1), uint16(1))
	}
}

func BenchmarkHandler_RemoveFromCart(b *testing.B) {
	repo := NewRepository(100)
	ctx := context.Background()
	userID := uint64(1)
	skuID := int64(123)

	_ = repo.AddToCart(ctx, skuID, userID, 1)

	for i := 0; i < b.N; i++ {
		_ = repo.RemoveFromCart(ctx, skuID, userID)
		_ = repo.AddToCart(ctx, skuID, userID, 1)
	}
}
