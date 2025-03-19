package delivery

import (
	"net/http"
)

type Router struct {
	server *Server
}

func NewRouter(server *Server) *Router {
	return &Router{server: server}
}

func (r *Router) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /user/{user_id}/cart/{sku_id}", r.server.AddToCartHandler)
	mux.HandleFunc("DELETE /user/{user_id}/cart/{sku_id}", r.server.RemoveFromCartHandler)
	mux.HandleFunc("DELETE /user/{user_id}/cart", r.server.ClearCartHandler)
	mux.HandleFunc("GET /user/{user_id}/cart", r.server.GetCartHandler)
	mux.HandleFunc("POST /cart/checkout", r.server.GetCartByUserIDHandler)
}
