package auth

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, h Handlers) *http.ServeMux{
	mux.Handle("GET /dummyLogin", http.HandlerFunc(h.DummyLoginHandler))
	mux.Handle("POST /register", http.HandlerFunc(h.Register))
	mux.Handle("POST /login", http.HandlerFunc(h.Login))
	return mux
}