package server

import (
	"avito-tech/internal/auth"
	"avito-tech/internal/flat"
	"avito-tech/internal/house"
	"net/http"
)

func (s *Server) MapRoutes() http.Handler {
	jwtKey := s.cfg.JWTSecret

	mainMux := http.NewServeMux()
	
	houseRepo := house.NewHouseRepository(s.db)
	houseService := house.NewHouseService(houseRepo)
	houseHandlers := house.NewHouseHandler(houseService)

	flatRepo := flat.NewFlatRepository(s.db)
	flatService := flat.NewFlatService(flatRepo)
	flatHandlers := flat.NewFlatHandler(flatService)

	authRepository := auth.NewAuthRepository(s.db)
	authService := auth.NewAuthService([]byte(s.cfg.JWTSecret), authRepository)
	authHandlers := auth.NewAuthHandler(authService)

	house.RegisterRoutes(mainMux, houseHandlers, []byte(jwtKey))
	flat.RegisterRoutes(mainMux, flatHandlers, []byte(jwtKey))
	auth.RegisterRoutes(mainMux, authHandlers)


	return mainMux
}