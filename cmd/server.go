package main

import (
	"encoding/json"
	"fmt"
	"github.com/vestamart/homework/internal/delivery"
	"io"
	"log"
	"net/http"
	"strconv"
)

type Cart struct {
	SkuID int64  `json:"sku_id"`
	Name  string `json:"name"`
	Count uint16 `json:"count"`
	Price uint32 `json:"price"`
}

type UserCart struct {
	Items      []Cart `json:"items"`
	TotalPrice uint32 `json:"total_price"`
}

func main() {
	log.Println("App started")

	cartMap := make(map[uint64]UserCart)
	//repository := repository.NewRepository(100)
	//
	//mux := http.NewServeMux()
	//mux.HandleFunc("POST /user/{us		er_id}/cart/{sku_id}", delivery.AddToCartHandler)
	//mux.HandleFunc("DELETE /user/{user_id}/cart/{sku_id}", delivery.RemoveFromCartHandler)
	//mux.HandleFunc("DELETE /user/{user_id}/cart", delivery.ClearCartHandler)
	//mux.HandleFunc("GET /user/{user_id}/cart", delivery.GetCartHandler)

	http.HandleFunc("POST /user/{user_id}/cart/{sku_id}", func(w http.ResponseWriter, r *http.Request) {
		rawUserID := r.PathValue("user_id")
		userId, err := strconv.ParseUint(rawUserID, 10, 64)
		RawSkuID := r.PathValue("sku_id")
		skuId, err := strconv.ParseInt(RawSkuID, 10, 64)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			_, errOut := fmt.Fprintf(w, "{\"error\":\"%s\"}", err)
			if errOut != nil {
				log.Printf("POST /user/{user_id}/cart/{sku_id} failed: %s", errOut)
				return
			}
			return
		}
		cart, exists := cartMap[userId]
		if !exists {
			cart = UserCart{}
		}
		if userId < 1 && skuId < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			_, errOut := fmt.Fprintf(w, "{\"error\":\"%s\"}", "userId or skuId must be greater than 1")
			if errOut != nil {
				log.Printf("POST /user/{user_id}/cart/{sku_id} failed: %s", errOut)
				return
			}
			return
		}

		body, err := io.ReadAll(r.Body)

		type AddToCartRequest struct {
			Count uint16 `json:"count"`
		}

		var addRequest AddToCartRequest

		if err := json.Unmarshal(body, &addRequest); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			_, errOut := fmt.Fprintf(w, "{\"error\":\"%s\"}", "JSON parse error")
			if errOut != nil {
				log.Printf("POST /user/{user_id}/cart/{sku_id} failed: %s", errOut)
				return
			}
			return
		}

		if addRequest.Count < 1 {
			w.WriteHeader(http.StatusPreconditionFailed)
			w.Header().Set("Content-Type", "application/json")
			_, errOut := fmt.Fprintf(w, "{\"error\":\"%s\"}", "count must be greater than 1")
			if errOut != nil {
				log.Printf("POST /user/{user_id}/cart/{sku_id} failed: %s", errOut)
				return
			}
			return
		}

		type ClientRequest struct {
			Name  string `json:"name"`
			Price uint32 `json:"price"`
		}

		var clientRequest ClientRequest

		if json.Unmarshal(delivery.GetProductHandler(skuId), &clientRequest); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			_, errOut := fmt.Fprintf(w, "{\"error\":\"%s\"}", "JSON parse error")
			if errOut != nil {
				log.Printf("POST http://route256.pavl.uk:8080/get_product failed: %s", errOut)
				return
			}
			return
		}

		cart.Items = append(cart.Items, Cart{
			SkuID: skuId,
			Name:  clientRequest.Name,
			Count: addRequest.Count,
			Price: clientRequest.Price,
		})

		cart.TotalPrice += uint32(addRequest.Count) * clientRequest.Price
		cartMap[userId] = cart
		fmt.Printf("ADD ITEM\tUserCart: %v\n", cartMap)
		return
	})

	http.HandleFunc("DELETE /user/{user_id}/cart/{sku_id}", func(w http.ResponseWriter, r *http.Request) {
		rawUserID := r.PathValue("user_id")
		userId, err := strconv.ParseUint(rawUserID, 10, 64)
		RawSkuID := r.PathValue("sku_id")
		skuId, err := strconv.ParseInt(RawSkuID, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			_, errOut := fmt.Fprintf(w, "{\"error\":\"%s\"}", err)
			if errOut != nil {
				log.Printf("DELETE /user/{user_id}/cart/{sku_id} failed: %s", errOut)
				return
			}
			return
		}
		var deletedItemPrice uint32
		cart, exists := cartMap[userId]
		if !exists {
			fmt.Fprintf(w, "{\"error\":\"%s\"}", "cart not found")
			return
		}
		// Новый список
		newItem := []Cart{}

		for _, item := range cart.Items {
			if item.SkuID == skuId {
				deletedItemPrice = uint32(item.Count) * item.Price
				continue
			}
			newItem = append(newItem, item)
		}
		cart.Items = newItem

		cart.TotalPrice -= deletedItemPrice

		if len(cart.Items) == 0 {
			delete(cartMap, userId)
		} else {
			cartMap[userId] = cart
		}

		fmt.Printf("DELETE ITEM\tUserCart: %v\n", cartMap)

		return
	})

	http.HandleFunc("DELETE /user/{user_id}/cart", func(w http.ResponseWriter, r *http.Request) {
		rawUserID := r.PathValue("user_id")
		userId, err := strconv.ParseUint(rawUserID, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			_, errOut := fmt.Fprintf(w, "{\"error\":\"%s\"}", err)
			if errOut != nil {
				log.Printf("DELETE /user/{user_id}/cart/ failed: %s", errOut)
				return
			}
			return
		}
		_, exists := cartMap[userId]
		if !exists {
			fmt.Fprintf(w, "{\"error\":\"%s\"}", "cart not found")
			return
		} else {
			delete(cartMap, userId)
			fmt.Printf("DELETE ALL\t UserCart: %v\n", cartMap)
		}
		return
	})

	http.HandleFunc("GET /user/{user_id}/cart", func(w http.ResponseWriter, r *http.Request) {
		rawUserID := r.PathValue("user_id")
		userId, err := strconv.ParseUint(rawUserID, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			_, errOut := fmt.Fprintf(w, "{\"error\":\"%s\"}", err)
			if errOut != nil {
				log.Printf("GET /user/{user_id}/cart/ failed: %s", errOut)
				return
			}
			return
		}
		if _, exists := cartMap[userId]; !exists {
			fmt.Fprintf(w, "{\"error\":\"%s\"}", "cart not found")
		}

		if jsonData, err := json.Marshal(cartMap[userId]); err != nil {
			log.Printf("GET /user/{user_id}/cart/ failed json parse: %s", err)
		} else {
			fmt.Fprintf(w, "%s", string(jsonData))
		}
		return
	})

	log.Print("Server running on port 8080")
	defer log.Print("Server shut down")
	if err := http.ListenAndServe("127.0.0.1:8082", nil); err != nil {
		panic(err)
	}
}
