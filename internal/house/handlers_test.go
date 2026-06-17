package house

import (
	"avito-tech/internal/entity"
	"avito-tech/internal/house/dto"
	"avito-tech/internal/middleware"
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

type houseServiceMock struct {
	CreateFunc func(ctx context.Context, house *dto.CreateHouseInput) (*entity.House, error)
}

func (m *houseServiceMock) Create(ctx context.Context, house *dto.CreateHouseInput) (*entity.House, error) {
	return m.CreateFunc(ctx, house)
}

func TestHandlerCreateHouse(t *testing.T) {
	tests := []struct {
		name string
		inputBody string
		inputRole entity.UserRole
		setupMock func(m *houseServiceMock)
		wantErr bool
		checkResult func(t *testing.T,rec *httptest.ResponseRecorder)
	}{
		{
			name: "Success: Дом успешно создан на уровне Handler",
			inputBody: `{"address": "Test_Address", "year": 1999, "developer": "someone"}`,
			inputRole: entity.RoleModerator,
			setupMock: func(m *houseServiceMock) {
				m.CreateFunc = func(ctx context.Context, house *dto.CreateHouseInput) (*entity.House, error) {
					if house.Address != "Test_Address" {
						t.Errorf("В сервис ушел неправильный адресс: %s", house.Address)
					}
					if house.Year != 1999 {
						t.Errorf("В сервис ушел неправильный год: %d", house.Year)
					}
					if house.Developer != "someone" {
						t.Errorf("В сервис ушел неправильный застройщик: %s", house.Developer)
					}

					return &entity.House{ID: 77, Address: house.Address, Year: house.Year, Developer: house.Developer}, nil
				}
			},
			checkResult: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusCreated {
					t.Errorf("Ожидали статус создания 201, получили %d", rec.Code)
				}
			},
		},
		{
			name: "Failed: Недопустимый уровень доступа",
			inputBody: `{"address": "Test_Address", "year": 1999, "developer": "someone"}`,
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
			name: "Failed: Неверный формат JSON",
			inputBody: `{"address": "Test_Address", "year": 1999, "developer": "someone"`,
			inputRole: entity.RoleModerator,
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
			mockService := &houseServiceMock{}
			if tc.setupMock != nil { tc.setupMock(mockService)}
			
			handlers := NewHouseHandler(mockService)
			req := httptest.NewRequest(http.MethodPost, "/house/create", bytes.NewBuffer([]byte(tc.inputBody)))
			ctx := context.WithValue(req.Context(), middleware.RoleKey, tc.inputRole)
			req = req.WithContext(ctx)
			rec := httptest.NewRecorder()

			handlers.Create(rec, req)

			if tc.checkResult != nil { tc.checkResult(t, rec) }
			}) 
	}
}

