package flat

import (
	"avito-tech/internal/entity"
	"avito-tech/internal/flat/dto"
	"avito-tech/internal/middleware"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type flatServiceMock struct {
	CreateFunc func(ctx context.Context, createdFlat *dto.CreateFlatInput) (*entity.Flat, error)
	GetAllByHouseIDFunc func(ctx context.Context, houseID uint) ([]entity.Flat, error)
	GetApprovedByHouseIDFunc func(ctx context.Context, houseID uint) ([]entity.Flat, error)
	ChangeStatusFunc func(ctx context.Context, flatID uint, newStatus string) (*entity.Flat, error)
}

func (m *flatServiceMock) ChangeStatus(ctx context.Context, flatID uint, newStatus string) (*entity.Flat, error) {
	return m.ChangeStatusFunc(ctx, flatID, newStatus)
}

// GetAllByHouseID implements [Service].
func (m *flatServiceMock) GetAllByHouseID(ctx context.Context, houseID uint) ([]entity.Flat, error) {
	return m.GetAllByHouseIDFunc(ctx, houseID)
}

// GetApprovedByHouseID implements [Service].
func (m *flatServiceMock) GetApprovedByHouseID(ctx context.Context, houseID uint) ([]entity.Flat, error) {
	return m.GetApprovedByHouseIDFunc(ctx, houseID)
}

// Create implements [Service].
func (m *flatServiceMock) Create(ctx context.Context, createdFlat *dto.CreateFlatInput) (*entity.Flat, error) {
	return m.CreateFunc(ctx, createdFlat)
}

func TestHandlerCreateFlat(t *testing.T) {
	tests := []struct {
		name string
		inputBody string
		setupMock func(m *flatServiceMock)
		wantErr bool
		checkResult func(t *testing.T,rec *httptest.ResponseRecorder)
	}{
		{
			name: "Success: Квартира успешно создана на уровне Handler",
			inputBody: `{"house_id": 1, "number": 11, "price": 10000, "rooms": 1}`,
			setupMock: func(m *flatServiceMock) {
				m.CreateFunc = func(ctx context.Context, createdFlat *dto.CreateFlatInput) (*entity.Flat, error) {
					if createdFlat.HouseID != 1 {
						t.Errorf("В сервис ушел неправильный номер дома: %d", createdFlat.HouseID)
					}
					if createdFlat.Number != 11 {
						t.Errorf("В сервис ушел неправильный номер квартиры: %d", createdFlat.Number)
					}
					if createdFlat.Price != 10000 {
						t.Errorf("В сервис ушла неправильная цена: %d", createdFlat.Price)
					}
					if createdFlat.Rooms != 1 {
						t.Errorf("В сервис ушло неверное кол-во комнат: %d", createdFlat.Price)
					}

					return &entity.Flat{ID: 7, HouseID: createdFlat.HouseID, Number: createdFlat.Number, Price: createdFlat.Price, Rooms: createdFlat.Rooms}, nil
				}
			},
			checkResult: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusCreated {
					t.Errorf("Ожидали статус создания 201, получили %d", rec.Code)
				}
			},
		},
		{
			name: "Failed: Ошибка на уровне сервиса",
			inputBody: `{"house_id": 1, "number": 11, "price": 10000, "rooms": 1}`,
			setupMock: func(m *flatServiceMock) {
				m.CreateFunc = func(ctx context.Context, createdFlat *dto.CreateFlatInput) (*entity.Flat, error) {
					return nil, fmt.Errorf("Какая то ошибка на уровне сервиса")
				}
			},
			wantErr: true,
			checkResult: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusInternalServerError {
					t.Errorf("Ожидали статус ошибки 500, получили %d", rec.Code)
				}
			},
		},
		{
			name: "Failed: Неверный формат JSON",
			inputBody: `{"house_id": 1, "number": 11, "price": 10000, "rooms": 1`,
			setupMock: nil,
			wantErr: true,
			checkResult: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusBadRequest {
					t.Errorf("Ожидали статус ошибки доступа 400, получили %d", rec.Code)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService := &flatServiceMock{}
			if tc.setupMock != nil { tc.setupMock(mockService)}
			
			handlers := NewFlatHandler(mockService)
			req := httptest.NewRequest(http.MethodPost, "/flat/create", bytes.NewBuffer([]byte(tc.inputBody)))
			rec := httptest.NewRecorder()

			handlers.Create(rec, req)

			if tc.checkResult != nil { tc.checkResult(t, rec) }
			}) 
	}
}

func TestGetByHouseID(t *testing.T) {
	tests := []struct {
		name string
		inputHouseID string
		inputRole entity.UserRole
		setupMock func(m *flatServiceMock)
		wantErr bool
		checkResult func(t *testing.T,rec *httptest.ResponseRecorder)
	}{
		{
			name: "Success: Получение списка квартир от имени Модератора",
			inputHouseID: "1",
			inputRole: entity.RoleModerator,
			setupMock: func(m *flatServiceMock) {
				m.GetAllByHouseIDFunc = func(ctx context.Context, houseID uint) ([]entity.Flat, error) {
					return []entity.Flat{{HouseID: 1}, {HouseID: 1}}, nil
				}
				m.GetApprovedByHouseIDFunc = nil
			},
			checkResult: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusOK {
					t.Errorf("Ожидали статус 200, получили %d", rec.Code)
				}
			},
		},
		{
			name: "Success: Получение списка квартир от имени Клиента",
			inputHouseID: "1",
			inputRole: entity.RoleClient,
			setupMock: func(m *flatServiceMock) {
				m.GetAllByHouseIDFunc = nil
				m.GetApprovedByHouseIDFunc = func(ctx context.Context, houseID uint) ([]entity.Flat, error) {
					return []entity.Flat{{HouseID: 1}, {HouseID: 1}}, nil
				}
			},
			checkResult: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusOK {
					t.Errorf("Ожидали статус 200, получили %d", rec.Code)
				}
			},
		},
		{
			name: "Failed: Неправильно указан номер квартиры",
			inputHouseID: "abc",
			inputRole: entity.RoleClient,
			setupMock: nil,
			wantErr: true,
			checkResult: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusBadRequest {
					t.Errorf("Ожидали статус 400, получили %d", rec.Code)
				}
			},
		},
		{
			name: "Failed: Отказано в доступе",
			inputHouseID: "1",
			inputRole: entity.UserRole(""),
			setupMock: nil,
			wantErr: true,
			checkResult: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusForbidden {
					t.Errorf("Ожидали статус ошибки доступа 403, получили %d", rec.Code)
				}
			},
		},
	}
		for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService := &flatServiceMock{}
			if tc.setupMock != nil { tc.setupMock(mockService)}
			
			handlers := NewFlatHandler(mockService)
			targetURL := fmt.Sprintf("/house/%s", tc.inputHouseID)
			req := httptest.NewRequest(http.MethodGet, targetURL, nil)
			req.SetPathValue("id", tc.inputHouseID)

			ctx := context.WithValue(req.Context(), middleware.RoleKey, tc.inputRole)
			req = req.WithContext(ctx)
			rec := httptest.NewRecorder()

			handlers.GetByHouseID(rec, req)

			if tc.checkResult != nil { tc.checkResult(t, rec) }
			}) 
	}
}

func TestChangeStatus(t *testing.T) {
	tests := []struct {
		name string
		inputRole entity.UserRole
		inputBody string
		setupMock func(m *flatServiceMock)
		wantErr bool
		checkResult func(t *testing.T,rec *httptest.ResponseRecorder)
	}{
		{
			name: "Success: Обновление статуса",
			inputBody: `{"id": 1, "status": "in moderate"}`,
			inputRole: entity.RoleModerator,
			setupMock: func(m *flatServiceMock) {
				m.ChangeStatusFunc = func(ctx context.Context, flatID uint, newStatus string) (*entity.Flat, error) {
					return &entity.Flat{Status: newStatus}, nil
				}
			},
			checkResult: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusOK {
					t.Errorf("Ожидали статус 200, получили %d", rec.Code)
				}
			},
		},
		{
			name: "Failed: Отказано в доступе",
			inputBody: `{"id": 1, "status": "in moderate"}`,
			inputRole: entity.RoleClient,
			setupMock: nil,
			wantErr: true,
			checkResult: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusForbidden {
					t.Errorf("Ожидали статус ошибки доступа 403, получили %d", rec.Code)
				}
			},
		},
		{
			name: "Failed: Неправильный формат JSON",
			inputBody: `{"id": 1, "status": "in moderate"`,
			inputRole: entity.RoleModerator,
			setupMock: nil,
			wantErr: true,
			checkResult: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusBadRequest {
					t.Errorf("Ожидали статус ошибки доступа 400, получили %d", rec.Code)
				}
			},
		},
		{
			name: "Failed: Ошибка на уровне сервиса",
			inputBody: `{"id": 1, "status": "in moderate"}`,
			inputRole: entity.RoleModerator,
			setupMock: func(m *flatServiceMock) {
				m.ChangeStatusFunc = func(ctx context.Context, flatID uint, newStatus string) (*entity.Flat, error) {
					return nil, fmt.Errorf("Ошибка на уровне сервиса")
				}
			},
			wantErr: true,
			checkResult: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusInternalServerError {
					t.Errorf("Ожидали статус ошибки доступа 500, получили %d", rec.Code)
				}
			},
		},
		
	
	}
		for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService := &flatServiceMock{}
			if tc.setupMock != nil { tc.setupMock(mockService)}
			
			handlers := NewFlatHandler(mockService)
			req := httptest.NewRequest(http.MethodPost, "/flat/update", bytes.NewBuffer([]byte(tc.inputBody)))

			ctx := context.WithValue(req.Context(), middleware.RoleKey, tc.inputRole)
			req = req.WithContext(ctx)
			rec := httptest.NewRecorder()

			handlers.ChangeStatus(rec, req)

			if tc.checkResult != nil { tc.checkResult(t, rec) }
			}) 
	}
}