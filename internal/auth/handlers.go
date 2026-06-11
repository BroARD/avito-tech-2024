package auth

import (
	"avito-tech/internal/auth/dto"
	"avito-tech/internal/entity"
	"encoding/json"
	"net/http"
)

var JwtKey = []byte("custom_key")

type Handlers interface {
	DummyLoginHandler(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}

type authHandler struct {
	serv Service
}

// Login implements [Handlers].
func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input dto.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Неправильные поля для авторизации"})
		return
	}

	response, err := h.serv.Login(r.Context(), &input)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Неправильный логин или пароль"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Register implements [Handlers].
func (h *authHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input dto.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Неправильные поля для регистрации"})
		return
	}

	user, err := h.serv.Register(r.Context(), &input)

	response := *&dto.RegisterResponse{
		Email: user.Email,
		Role: user.Role,
	}
	
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка при регистрации"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DummyLoginHandler implements [Handlers].
func (h *authHandler) DummyLoginHandler(w http.ResponseWriter, r *http.Request) {
	role := entity.UserRole(r.URL.Query().Get("role"))
	if role != entity.RoleClient && role != entity.RoleModerator {
		role = entity.RoleClient
	}

	token, err := h.serv.GenerateDummyToken(r.Context(), role)
	if err != nil {
		http.Error(w, "Ошибка при создании токена", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
		"role":  string(role),
	})
}

func NewAuthHandler(serv Service) Handlers {
	return &authHandler{
		serv: serv,
	}
}
