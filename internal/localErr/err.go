package localErr

import "errors"

var ErrSkuNotExist = errors.New("sku not exist")
var ItemNotEnoughErr = errors.New("item not enough")
