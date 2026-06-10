package flat

import (
	"avito-tech/internal/middleware"
	"net/http"
)


func RegisterRoutes(mux *http.ServeMux, h Handlers, jwtKey []byte) *http.ServeMux{
	mux.Handle("POST /flat/create", http.HandlerFunc(h.Create))
	mux.Handle("GET /house/{id}", middleware.AuthMiddleware(jwtKey, http.HandlerFunc(h.GetByHouseID)))
	mux.Handle("POST /flat/update", middleware.AuthMiddleware(jwtKey, http.HandlerFunc(h.ChangeStatus)))
	return mux
}