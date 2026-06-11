package auth

import (
	"context"
	"fmt"
	"testing"
	"time"

	"avito-tech/internal/auth/dto"
	"avito-tech/internal/entity"
)

// 1. Исправленный мок-репозиторий
type authRepoMock struct {
	RegisterFunc   func(ctx context.Context, user *entity.User) (*entity.User, error)
	GetByEmailFunc func(ctx context.Context, email string) (*entity.User, error)
}

func (m *authRepoMock) Register(ctx context.Context, user *entity.User) (*entity.User, error) {
	return m.RegisterFunc(ctx, user)
}

// ИСПРАВЛЕНО: вызываем замыкание Func, а не сам метод рекурсивно
func (m *authRepoMock) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	return m.GetByEmailFunc(ctx, email)
}

func TestRegisterUser(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	// Описываем структуру одного тест-кейса
	type testCase struct {
		name         string
		input        dto.RegisterInput
		setupMock    func(m *authRepoMock)
		wantErr      bool
		expectedErr  string
		checkResult  func(t *testing.T, res *entity.User) 
	}

	// Создаем "таблицу" со всеми сценариями
	tests := []testCase{
		{
			name: "Success - Успешная регистрация",
			input: dto.RegisterInput{Email: "valid@mail.ru", Password: "1234", Role: "client"},
			setupMock: func(m *authRepoMock) {
				m.RegisterFunc = func(ctx context.Context, u *entity.User) (*entity.User, error) {
					u.ID = 7
					u.CreatedAt = now
					return u, nil
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, res *entity.User) {
				if res.ID != 7 { 
					t.Errorf("Ожидался ID 7, получили %d", res.ID) 
				}
				if res.Email != "valid@mail.ru" {
					t.Errorf("Ожидался Email 'valid@mail.ru', получили %s", res.Email)
				}
			},
		},
		{
			name: "Failure - Email уже зарегистрирован",
			input: dto.RegisterInput{Email: "existing@mail.ru", Password: "1234", Role: "client"},
			setupMock: func(m *authRepoMock) {
				m.RegisterFunc = func(ctx context.Context, u *entity.User) (*entity.User, error) {
					return nil, fmt.Errorf("user with this email already exists")
				}
			},
			wantErr:     true,
			expectedErr: "user with this email already exists",
		},
		{
			name: "Failure - Сбой базы данных",
			input: dto.RegisterInput{Email: "db_error@mail.ru", Password: "1234", Role: "client"},
			setupMock: func(m *authRepoMock) {
				m.RegisterFunc = func(ctx context.Context, u *entity.User) (*entity.User, error) {
					return nil, fmt.Errorf("connection refused")
				}
			},
			wantErr:     true,
			expectedErr: "connection refused",
		},
	}

	// Запускаем цикл по всей таблице тестов
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &authRepoMock{}
			if tc.setupMock != nil {
				tc.setupMock(mockRepo)
			}

			service := NewAuthService([]byte("secret"), mockRepo)
			
			result, err := service.Register(context.Background(), &tc.input)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("Ожидалась ошибка, но метод выполнился успешно")
				}
				if err.Error() != tc.expectedErr {
					t.Errorf("Ожидалась ошибка '%s', но получена '%s'", tc.expectedErr, err.Error())
				}
				return 
			}

			if err != nil {
				t.Fatalf("Не ожидали ошибку, но получили: %v", err)
			}

			if tc.checkResult != nil {
				tc.checkResult(t, result)
			}
		})
	}
}

