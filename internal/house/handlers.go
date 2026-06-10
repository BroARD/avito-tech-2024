package house

import (
	"avito-tech/internal/entity"
	"avito-tech/internal/house/dto"
	"avito-tech/internal/middleware"
	"encoding/json"
	"net/http"
)

type Handlers interface {
	Create(w http.ResponseWriter, r *http.Request)
}

type houseHandler struct {
	serv Service
}

// Create implements [Handlers].
func (h *houseHandler) Create(w http.ResponseWriter, r *http.Request) {
	userRole, ok := r.Context().Value(middleware.RoleKey).(entity.UserRole)
	
	if !ok || userRole != entity.RoleModerator { 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": "Создать дом может только модерация!"})
		return
	}

	var input dto.CreateHouseInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Неправильные поля для создания"})
		return
	}

	house, err := h.serv.Create(r.Context(), &input)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка при создании дома"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(house)
}

func NewHouseHandler(serv Service) Handlers {
	return &houseHandler{
		serv: serv,
	}
}
