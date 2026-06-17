package auth

import (
	"avito-tech/internal/auth/dto"
	"avito-tech/internal/entity"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)


type authServiceMock struct {
	GenerateDummyTokenFunc func(ctx context.Context, role entity.UserRole) (string, error)
	RegisterFunc func(ctx context.Context, newUser *dto.RegisterInput) (*entity.User, error)
	LoginFunc func(ctx context.Context, input *dto.LoginInput) (*dto.LoginResponse, error)
}

func (m *authServiceMock) Login(ctx context.Context, input *dto.LoginInput) (*dto.LoginResponse, error) {
	return m.LoginFunc(ctx, input)
}

// Register implements [Service].
func (m *authServiceMock) Register(ctx context.Context, newUser *dto.RegisterInput) (*entity.User, error) {
	return m.RegisterFunc(ctx, newUser)
}

// GenerateDummyToken implements [Service].
func (m *authServiceMock) GenerateDummyToken(ctx context.Context, role entity.UserRole) (string, error) {
	return m.GenerateDummyTokenFunc(ctx, role)
}

func TestHandlerRegister(t *testing.T) {
	tests := []struct {
		name string
		inputBody string
		setupMock func(m *authServiceMock)
		wantErr bool
		checkResult func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "Success: Регистрация пользователя",
			inputBody: `{"email": "test@mail.ru", "password": "1234", "role": "client"}`,
			setupMock: func(m *authServiceMock) {
				m.RegisterFunc = func(ctx context.Context, newUser *dto.RegisterInput) (*entity.User, error) {
					return &entity.User{ID: 1, Email: "test@mail.ru", Role: "client"}, nil
				}
			},
			checkResult: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusCreated {
					t.Errorf("Ожидался статус 201, получили %d", rec.Code)
				}
			},
		},
		{
			name: "Failed: Неверный формат JSON",
			inputBody: `{"email": "test@mail.ru", "password": "1234", "role": "client"`,
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
			inputBody: `{"email": "test@mail.ru", "password": "1234", "role": "client"}`,
			setupMock: func(m *authServiceMock) {
				m.RegisterFunc = func(ctx context.Context, newUser *dto.RegisterInput) (*entity.User, error) {
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
	}

		for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService := &authServiceMock{}
			if tc.setupMock != nil { tc.setupMock(mockService)}
			
			handlers := NewAuthHandler(mockService)
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer([]byte(tc.inputBody)))
			rec := httptest.NewRecorder()

			handlers.Register(rec, req)

			if tc.checkResult != nil { tc.checkResult(t, rec) }
			}) 
	}
}

func TestHandlerLogin(t *testing.T) {
	tests := []struct {
		name string
		inputBody string
		setupMock func(m *authServiceMock)
		wantErr bool
		checkResult func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "Success: Авторизация пользователя",
			inputBody: `{"email": "test@mail.ru", "password": "1234"}`,
			setupMock: func(m *authServiceMock) {
				m.LoginFunc = func(ctx context.Context, input *dto.LoginInput) (*dto.LoginResponse, error) {
					return &dto.LoginResponse{}, nil
				}
			},
			checkResult: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusOK {
					t.Errorf("Ожидался статус 200, получили %d", rec.Code)
				}
			},
		},
		{
			name: "Failed: Неверный формат JSON",
			inputBody: `{"email": "test@mail.ru", "password": "1234"`,
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
			inputBody: `{"email": "test@mail.ru", "password": "1234"}`,
			setupMock: func(m *authServiceMock) {
				m.LoginFunc = func(ctx context.Context, input *dto.LoginInput) (*dto.LoginResponse, error) {
					return nil, fmt.Errorf("Ошибка на уровне сервиса")
				}
			},
			wantErr: true,
			checkResult: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusInternalServerError {
					t.Errorf("Ожидали статус ошибки 500, получили %d", rec.Code)
				}
			},
		},
	}

		for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService := &authServiceMock{}
			if tc.setupMock != nil { tc.setupMock(mockService)}
			
			handlers := NewAuthHandler(mockService)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte(tc.inputBody)))
			rec := httptest.NewRecorder()

			handlers.Login(rec, req)

			if tc.checkResult != nil { tc.checkResult(t, rec) }
			}) 
	}
}

func TestHandlerDummyLogin(t *testing.T) {
	tests := []struct {
		name string
		inputRole entity.UserRole
		setupMock func(m *authServiceMock)
		wantErr bool
		checkResult func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "Success: Получение токена для Dummy",
			inputRole: entity.RoleClient,
			setupMock: func(m *authServiceMock) {
				m.GenerateDummyTokenFunc = func(ctx context.Context, role entity.UserRole) (string, error) {
					return "test_token", nil
				}
			},
			checkResult: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusOK {
					t.Errorf("Ожидался статус 200, получили %d", rec.Code)
				}
			},
		},
				{
			name: "Failed: Ошибка на уровне сервиса",
			inputRole: entity.RoleClient,
			setupMock: func(m *authServiceMock) {
				m.GenerateDummyTokenFunc = func(ctx context.Context, role entity.UserRole) (string, error) {
					return "", fmt.Errorf("Ошибка на уровне сервиса")
				}
			},
			wantErr: true,
			checkResult: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusInternalServerError {
					t.Errorf("Ожидали статус ошибки 500, получили %d", rec.Code)
				}
			},
		},
	}

		for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService := &authServiceMock{}
			if tc.setupMock != nil { tc.setupMock(mockService)}

			handlers := NewAuthHandler(mockService)
			
			targetURL := fmt.Sprintf("/dummyLogin?role=%s", tc.inputRole)
		
			req := httptest.NewRequest(http.MethodGet, targetURL, nil)
			rec := httptest.NewRecorder()

			handlers.DummyLoginHandler(rec, req)

			if tc.checkResult != nil { tc.checkResult(t, rec) }
			}) 
	}
}