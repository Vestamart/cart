package delivery

import (
	"encoding/json"
	"errors"
	"github.com/vestamart/homework/internal/app"
	"github.com/vestamart/homework/internal/domain"
	"io"
	"log"
	"net/http"
	"strconv"
)

type GetCartResponse struct {
	Items      []GetCartItemResponse `json:"items"`
	TotalPrice uint32                `json:"total_price"`
}

type GetCartItemResponse struct {
	Sku   int64  `json:"sku_id"`
	Name  string `json:"name"`
	Count uint16 `json:"count"`
	Price uint32 `json:"price"`
}

type Server struct {
	cartService app.CartService
}

func NewServer(cartService app.CartService) *Server {
	return &Server{cartService: cartService}
}

// AddToCartRequest Request form
type AddToCartRequest struct {
	Count uint16 `json:"count"`
}

// Server Handlers

func (s Server) AddToCartHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rawUserID := r.PathValue("user_id")
	userID, err := strconv.ParseUint(rawUserID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	RawSkuID := r.PathValue("sku_id")
	skuID, err := strconv.ParseInt(RawSkuID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("DROP ON HTTP:", err)
		}
	}(r.Body)

	var addToCartRequest AddToCartRequest
	if err = json.NewDecoder(r.Body).Decode(&addToCartRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if addToCartRequest.Count < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.cartService.AddToCart(r.Context(), skuID, userID, addToCartRequest.Count)
	if err != nil {
		if errors.Is(err, domain.ErrSkuNotExist) {
			w.WriteHeader(http.StatusPreconditionFailed)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	w.WriteHeader(http.StatusOK)
}

func (s Server) RemoveFromCartHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rawUserID := r.PathValue("user_id")
	userID, err := strconv.ParseUint(rawUserID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	RawSkuID := r.PathValue("sku_id")
	skuID, err := strconv.ParseInt(RawSkuID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.cartService.RemoveFromCart(r.Context(), skuID, userID)
}

func (s Server) ClearCartHandler(w http.ResponseWriter, r *http.Request) {
	rawUserID := r.PathValue("user_id")
	userID, err := strconv.ParseUint(rawUserID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.cartService.ClearCart(r.Context(), userID)

}

func (s Server) GetCartHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rawUserID := r.PathValue("user_id")
	userID, err := strconv.ParseUint(rawUserID, 10, 64)
	if err != nil || userID < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cart, err := s.cartService.GetCart(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp := GetCartResponse{
		Items:      make([]GetCartItemResponse, 0, len(cart.Items)),
		TotalPrice: 0,
	}

	for _, item := range cart.Items {
		resp.Items = append(resp.Items, GetCartItemResponse{
			Sku:   item.Sku,
			Name:  item.Name,
			Count: item.Count,
			Price: item.Price,
		})
	}
	resp.TotalPrice = cart.TotalPrice

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(cart)
}
