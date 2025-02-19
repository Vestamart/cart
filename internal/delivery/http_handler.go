package delivery

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vestamart/homework/internal/app"
	"github.com/vestamart/homework/internal/domain"
	"net/http"
	"strconv"
)

type Server struct {
	cartService app.CartRepository
}

func NewServer(cartService app.CartRepository) *Server {
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

	defer r.Body.Close()

	var addToCartRequest AddToCartRequest
	if err = json.NewDecoder(r.Body).Decode(&addToCartRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if addToCartRequest.Count < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = s.cartService.AddToCart(r.Context(), skuID, userID, addToCartRequest.Count)
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

	_, err = s.cartService.RemoveFromCart(r.Context(), skuID, userID)
}

func (s Server) ClearCartHandler(w http.ResponseWriter, r *http.Request) {
	rawUserID := r.PathValue("user_id")
	userID, err := strconv.ParseUint(rawUserID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = s.cartService.ClearCart(r.Context(), userID)

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

	if jsonData, err := json.Marshal(cart); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		fmt.Fprint(w, string(jsonData))
	}
}
