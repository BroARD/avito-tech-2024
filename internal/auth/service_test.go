package auth

import (
	"context"
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

// 2. Исправленный тест
func TestRegisterUser_Success(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	mockAuthRepo := &authRepoMock{
		RegisterFunc: func(ctx context.Context, user *entity.User) (*entity.User, error) {
			if user.PasswordHash == "1234" {
				t.Error("Сервис не захэшировал пароль перед сохранением в репозиторий")
			}

			user.ID = 7
			user.CreatedAt = now
			return user, nil
		},
	}

	mockJWTKey := []byte("mock_jwt_key")
	service := NewAuthService(mockJWTKey, mockAuthRepo)

	input := dto.RegisterInput{
		Email:    "test@mail.ru",
		Password: "1234",
		Role:     "client",
	}

	result, err := service.Register(context.Background(), &input)

	if err != nil {
		t.Fatalf("Ожидался успешный результат, но получена ошибка: %v", err)
	}

	if result.ID != 7 {
		t.Errorf("Ожидался ID пользователя = 7, получили = %d", result.ID)
	}

	if result.Email != "test@mail.ru" {
		t.Errorf("Ожидался Email = test@mail.ru, получили = %s", result.Email)
	}

	if !result.CreatedAt.Equal(now) {
		t.Errorf("Ожидалось время CreatedAt = %v, получили = %v", now, result.CreatedAt)
	}
}
