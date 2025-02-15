package delivery

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// Server Handlers
func AddToCartHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("AddToCartHandler")
}

func RemoveFromCartHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Remove from Cart")
}

func ClearCartHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Clear Cart")
}

func GetCartHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Get Cart")
}

// Client Handlers
func GetProductHandler(sku int64) []byte {
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
	defer resp.Body.Close()
	buffer, err := io.ReadAll(resp.Body)
	return buffer
}
