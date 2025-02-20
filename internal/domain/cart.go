package domain

import "errors"

type UserCart struct {
	Items map[int64]uint16 `json:"items"`
}

var ErrSkuNotExist = errors.New("sku not exist")
