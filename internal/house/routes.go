package house

import (
	"avito-tech/internal/middleware"
	"net/http"
)


func RegisterRoutes(mux *http.ServeMux, h Handlers, jwtKey []byte) *http.ServeMux{
	mux.Handle("POST /house/create", middleware.AuthMiddleware(jwtKey, http.HandlerFunc(h.Create)))

	return mux
}