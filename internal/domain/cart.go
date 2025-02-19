package domain

import "errors"

type Sku = int64

type Item struct {
	SkuID Sku    `json:"sku_id"`
	Name  string `json:"name"`
	Count uint16 `json:"count"`
	Price uint32 `json:"price"`
}

type UserCart struct {
	Items      []Item `json:"items"`
	TotalPrice uint32 `json:"total_price"`
}

type ClientResponse struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

var ErrSkuNotExist = errors.New("sku not exist")
