package domain

type UserCart struct {
	Items      []CartItem `json:"items"`
	TotalPrice uint32     `json:"total_price"`
}

type CartItem struct {
	Sku   int64  `json:"sku"`
	Name  string `json:"name"`
	Count uint16 `json:"count"`
	Price uint32 `json:"price"`
}

type ProductServiceResponse struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}
