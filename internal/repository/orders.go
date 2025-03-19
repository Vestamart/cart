package repository

import (
	"context"
	"errors"
	"github.com/vestamart/homework/internal/domain"
)

// Errors
var (
	OrderNotExistErr = errors.New("order does not exist")
)

type OrderID = int64

type OrdersStorage = map[OrderID]domain.Order

type InMemoryOrderRepository struct {
	orderStorage OrdersStorage
	lastOrderID  OrderID
}

func NewInMemoryOrderRepository(cap int) *InMemoryOrderRepository {
	return &InMemoryOrderRepository{orderStorage: make(OrdersStorage, cap), lastOrderID: 0}
}

func (r *InMemoryOrderRepository) Create(_ context.Context, userID int64, items *[]domain.Item) (OrderID, error) {
	r.lastOrderID++
	orderID := r.lastOrderID

	r.orderStorage[orderID] = domain.Order{
		UserID: userID,
		Status: 0,
		Items:  *items,
	}

	return orderID, nil
}

func (r *InMemoryOrderRepository) SetStatus(_ context.Context, orderID int64, status domain.OrderStatus) error {
	v, ok := r.orderStorage[orderID]
	if !ok {
		return OrderNotExistErr
	}
	v.Status = status

	r.orderStorage[orderID] = v
	return nil
}

func (r *InMemoryOrderRepository) GetByID(_ context.Context, orderID int64) (*domain.Order, error) {
	v, ok := r.orderStorage[orderID]
	if !ok {
		return nil, OrderNotExistErr
	}

	return &v, nil
}
