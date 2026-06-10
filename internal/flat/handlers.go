package flat

import (
	"avito-tech/internal/entity"
	"avito-tech/internal/flat/dto"
	"avito-tech/internal/middleware"
	"encoding/json"
	"net/http"
	"strconv"
)

type Handlers interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetByHouseID(w http.ResponseWriter, r *http.Request)
	ChangeStatus(w http.ResponseWriter, r *http.Request)
}

type flatHandler struct {
	serv Service
}

// ChangeStatus implements [Handlers].
func (h *flatHandler) ChangeStatus(w http.ResponseWriter, r *http.Request) {
	var input dto.UpdateStatusuInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Неправильные поля для смены статуса"})
		return
	}

	userRole, ok := r.Context().Value(middleware.RoleKey).(entity.UserRole)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": "АВТОРИЗУЙСЯ ДЕБИК"})
		return
	}
	if userRole != "moderator" {
		http.Error(w, `{"error":"Это не твой уровень дорогой"}`, http.StatusBadRequest)
		return
	}

	flat, err := h.serv.ChangeStatus(r.Context(), input.FlatID, input.NewStatus)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(flat)

}

// GetByHouseID implements [Handlers].
func (h *flatHandler) GetByHouseID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	houseID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error":"Неправильный ID дома"}`, http.StatusBadRequest)
	}

	userRole, ok := r.Context().Value(middleware.RoleKey).(entity.UserRole)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": "АВТОРИЗУЙСЯ ДЕБИК"})
		return
	}

	if userRole == "moderator" {
		flats, err := h.serv.GetAllByHouseID(r.Context(), uint(houseID))
		if err != nil {
			http.Error(w, `{"error": "Ошибка при получение квартир из этого дома"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(flats)
	} else {
		flats, err := h.serv.GetApprovedByHouseID(r.Context(), uint(houseID))
		if err != nil {
			http.Error(w, `{"error": "Ошибка при получение квартир из этого дома"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(flats)
	}
}

// Create implements [Handlers].
func (h *flatHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateFlatInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Неправильные поля для создания"})
		return
	}

	flat, err := h.serv.Create(r.Context(), &input)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(flat)
}

func NewFlatHandler(serv Service) Handlers {
	return &flatHandler{
		serv: serv,
	}
}
