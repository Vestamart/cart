package domain

import "errors"

type OrderStatus int

const (
	NEW OrderStatus = iota
	AWWAITING_PAYMENT
	FAILED
	PAYED
	CANCELLED
)

type Order struct {
	UserID int64
	Status OrderStatus
	Items  []Item
}

type Item struct {
	Sku   uint32
	Count uint32
}

type StocksItem struct {
	TotalCount uint32 `json:"total_count"`
	Reserved   uint32 `json:"reserved"`
}

var ItemNotEnoughErr = errors.New("item not enough")
