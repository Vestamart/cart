package domain

type Sku = int64

type UserCart struct {
	Items      []Item
	TotalPrice uint32
}

type Item struct {
	SkuID Sku    `json:"sku_id"`
	Name  string `json:"name"`
	Count uint16 `json:"count"`
	Price uint32 `json:"price"`
}

type AddToCartRequest struct {
	Count uint16 `json:"count"`
}

type ClientRequest struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}
