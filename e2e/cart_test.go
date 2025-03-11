package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vestamart/homework/internal/delivery"
)

func makeAddToCartRequest(t *testing.T, userID, skuID string, count uint16) *http.Request {
	reqBody := delivery.AddToCartRequest{Count: count}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}
	req, err := http.NewRequest("POST", "/user/"+userID+"/cart/"+skuID, bytes.NewBuffer(bodyBytes))
	if err != nil {
		t.Fatalf("failed to create add request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req
}

func makeRequest(t *testing.T, method, path string) *http.Request {
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		t.Fatalf("failed to create %s request: %v", method, err)
	}
	return req
}

func TestRemoveFromCartHandler(t *testing.T) {
	tests := []struct {
		name      string
		userID    string
		skuID     string
		count     uint16
		wantEmpty bool
	}{
		{
			name:      "Remove item",
			userID:    "1",
			skuID:     "123",
			count:     2,
			wantEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := SetupTestServer()

			reqAdd := makeAddToCartRequest(t, tt.userID, tt.skuID, tt.count)
			rrAdd := httptest.NewRecorder()
			server.Handler.ServeHTTP(rrAdd, reqAdd)
			assert.Equal(t, http.StatusOK, rrAdd.Code, "expected POST to succeed")

			req := makeRequest(t, "DELETE", "/user/"+tt.userID+"/cart/"+tt.skuID)
			rr := httptest.NewRecorder()
			server.Handler.ServeHTTP(rr, req)
			assert.Equal(t, http.StatusOK, rr.Code, "expected DELETE to return 200 OK")

			reqGet := makeRequest(t, "GET", "/user/"+tt.userID+"/cart")
			rrGet := httptest.NewRecorder()
			server.Handler.ServeHTTP(rrGet, reqGet)
			assert.Equal(t, http.StatusOK, rrGet.Code, "expected GET to return 200 OK")

			var resp delivery.GetCartResponse
			err := json.NewDecoder(rrGet.Body).Decode(&resp)
			assert.NoError(t, err, "expected no error decoding response")
			if tt.wantEmpty {
				assert.Empty(t, resp.Items, "cart should be empty after deletion")
			}
		})
	}
}

func TestGetCartHandler(t *testing.T) {
	tests := []struct {
		name      string
		userID    string
		skuID     string
		count     uint16
		wantItems int
		wantSku   int64
		wantCount uint16
		wantPrice uint32
		wantTotal uint32
	}{
		{
			name:      "Name item",
			userID:    "2",
			skuID:     "456",
			count:     3,
			wantItems: 1,
			wantSku:   456,
			wantCount: 3,
			wantPrice: 100,
			wantTotal: 300,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := SetupTestServer()

			reqAdd := makeAddToCartRequest(t, tt.userID, tt.skuID, tt.count)
			rrAdd := httptest.NewRecorder()
			server.Handler.ServeHTTP(rrAdd, reqAdd)
			assert.Equal(t, http.StatusOK, rrAdd.Code, "expected POST to succeed")

			req := makeRequest(t, "GET", "/user/"+tt.userID+"/cart")
			rr := httptest.NewRecorder()
			server.Handler.ServeHTTP(rr, req)
			assert.Equal(t, http.StatusOK, rr.Code, "expected GET to return 200 OK")

			var resp delivery.GetCartResponse
			err := json.NewDecoder(rr.Body).Decode(&resp)
			assert.NoError(t, err, "expected no error decoding response")

			assert.Len(t, resp.Items, tt.wantItems, "expected correct number of items")
			if tt.wantItems > 0 {
				assert.Equal(t, tt.wantSku, resp.Items[0].Sku, "expected correct skuID")
				assert.Equal(t, tt.wantCount, resp.Items[0].Count, "expected correct count")
				assert.Equal(t, tt.wantPrice, resp.Items[0].Price, "expected correct price")
				assert.Equal(t, tt.wantTotal, resp.TotalPrice, "expected correct total price")
			}
		})
	}
}
