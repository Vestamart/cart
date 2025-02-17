package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vestamart/homework/internal/domain"
	"io"
	"net/http"
	"strconv"
)

var ErrorEmptyClientBody = errors.New("body is empty")

// Server
type CartServer interface {
	AddToCart(_ context.Context, skuID int64, userID uint64, count uint16) (*domain.UserCart, error)
	RemoveFromCart(_ context.Context, skiID int64, userID uint64) (*domain.UserCart, error)
	ClearCart(_ context.Context, userID uint64) (*domain.UserCart, error)
	GetCart(_ context.Context, userId uint64) ([]byte, error)
}
type Server struct {
	cartService CartServer
}

func NewServer(cartService CartServer) *Server {
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
	userId, err := strconv.ParseUint(rawUserID, 10, 64)

	RawSkuID := r.PathValue("sku_id")
	skuId, err := strconv.ParseInt(RawSkuID, 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	body, _ := io.ReadAll(r.Body)

	var addToCartRequest AddToCartRequest

	if err := json.Unmarshal(body, &addToCartRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	if addToCartRequest.Count < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = s.cartService.AddToCart(r.Context(), skuId, userId, addToCartRequest.Count)
	if err != nil {
		if errors.Is(err, ErrorEmptyClientBody) {
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
	userId, err := strconv.ParseUint(rawUserID, 10, 64)
	RawSkuID := r.PathValue("sku_id")
	skuId, err := strconv.ParseInt(RawSkuID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = s.cartService.RemoveFromCart(r.Context(), skuId, userId)
}

func (s Server) ClearCartHandler(w http.ResponseWriter, r *http.Request) {
	rawUserID := r.PathValue("user_id")
	userId, err := strconv.ParseUint(rawUserID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = s.cartService.ClearCart(r.Context(), userId)

}

func (s Server) GetCartHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rawUserID := r.PathValue("user_id")
	userId, err := strconv.ParseUint(rawUserID, 10, 64)
	if err != nil || userId < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = s.cartService.GetCart(r.Context(), userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// Client Handlers
func GetProductHandler(sku int64) (*domain.ClientRequest, error) {
	url := "http://route256.pavl.uk:8080/get_product"
	text := fmt.Sprintf("{\n  \"token\": \"testtoken\",\n  \"sku\": %d\n}", sku)
	jsonBody := []byte(text)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, ErrorEmptyClientBody
	}
	defer resp.Body.Close()
	buffer, err := io.ReadAll(resp.Body)

	var clientRequest domain.ClientRequest
	if json.Unmarshal(buffer, &clientRequest) != nil {
		fmt.Println("failed parsing request body")
	}
	return &clientRequest, nil
}
