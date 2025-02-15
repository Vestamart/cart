package app

import "context"

type CartService interface {
	AddToCartHandler(ctx context.Context)
}
