package delivery

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/vestamart/cart/internal/app/cart"
	"github.com/vestamart/cart/internal/localErr"
	"io"
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
	cartService cart.Service
	validator   *validator.Validate
}

func NewServer(cartService cart.Service) *Server {
	return &Server{
		cartService: cartService,
		validator:   validator.New(),
	}
}

// AddToCartRequest Request form
type AddToCartRequest struct {
	Count uint16 `json:"count" validate:"required,min=1"`
}

// GetCartByUserID
type GetCartByUserIDRequest struct {
	UserID uint64 `json:"user" validate:"required,min=1"`
}

// Server Handlers

func (s Server) AddToCartHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rawUserID := r.PathValue("user_id")
	userID, err := strconv.ParseUint(rawUserID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid user_id format"})
		return
	}

	RawSkuID := r.PathValue("sku_id")
	skuID, err := strconv.ParseInt(RawSkuID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid sku_id format"})
		return
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
	}(r.Body)

	var addToCartRequest AddToCartRequest
	if err = json.NewDecoder(r.Body).Decode(&addToCartRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	if err = s.validator.Struct(addToCartRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	if userID < 1 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "user_id must be positive"})
		return
	}

	if skuID < 1 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "sku_id must be positive"})
		return
	}

	err = s.cartService.AddToCart(r.Context(), skuID, userID, addToCartRequest.Count)
	if err != nil {
		if errors.Is(err, localErr.ErrSkuNotExist) {
			w.WriteHeader(http.StatusPreconditionFailed)
			json.NewEncoder(w).Encode(map[string]string{"error": "SKU does not exist"})
			return
		}
		if errors.Is(err, localErr.ItemNotEnoughErr) {
			w.WriteHeader(http.StatusPreconditionFailed)
			json.NewEncoder(w).Encode(map[string]string{"error": "Not enough items in stock"})
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to add item to cart"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Item added to cart successfully"})
}

func (s Server) RemoveFromCartHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rawUserID := r.PathValue("user_id")
	userID, err := strconv.ParseUint(rawUserID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid user_id format"})
		return
	}

	RawSkuID := r.PathValue("sku_id")
	skuID, err := strconv.ParseInt(RawSkuID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid sku_id format"})
		return
	}

	// Валидация параметров пути
	if userID < 1 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "user_id must be positive"})
		return
	}

	if skuID < 1 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "sku_id must be positive"})
		return
	}

	err = s.cartService.RemoveFromCart(r.Context(), skuID, userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to remove item from cart"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Item removed from cart successfully"})
}

func (s Server) ClearCartHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rawUserID := r.PathValue("user_id")
	userID, err := strconv.ParseUint(rawUserID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid user_id format"})
		return
	}

	// Валидация параметров пути
	if userID < 1 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "user_id must be positive"})
		return
	}

	err = s.cartService.ClearCart(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to clear cart"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Cart cleared successfully"})
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
	_ = json.NewEncoder(w).Encode(resp)
}

func (s Server) GetCartByUserIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var getCartByUserID GetCartByUserIDRequest
	if err := json.NewDecoder(r.Body).Decode(&getCartByUserID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	if err := s.validator.Struct(getCartByUserID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	orderID, err := s.cartService.CheckoutCart(r.Context(), getCartByUserID.UserID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to checkout cart"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Cart checked out successfully",
		"order_id": orderID,
	})
}
