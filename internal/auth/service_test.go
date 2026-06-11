package auth

import (
	"context"
	"fmt"
	"testing"
	"time"

	"avito-tech/internal/auth/dto"
	"avito-tech/internal/entity"

	"golang.org/x/crypto/bcrypt"
)


type authRepoMock struct {
	RegisterFunc   func(ctx context.Context, user *entity.User) (*entity.User, error)
	GetByEmailFunc func(ctx context.Context, email string) (*entity.User, error)
}

func (m *authRepoMock) Register(ctx context.Context, u *entity.User) (*entity.User, error)   { return m.RegisterFunc(ctx, u) }
func (m *authRepoMock) GetByEmail(ctx context.Context, e string) (*entity.User, error) { return m.GetByEmailFunc(ctx, e) }

var testPasswordHash string

func init() {
	hash, _ := bcrypt.GenerateFromPassword([]byte("1234"), bcrypt.MinCost)
	testPasswordHash = string(hash)
}
func TestRegisterUser(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	tests := []struct {
		name        string
		input       dto.RegisterInput
		setupMock   func(m *authRepoMock)
		wantErr     bool
		expectedErr string
		checkResult func(t *testing.T, res *entity.User)
	}{
		{
			name:  "Success - Успешная регистрация",
			input: dto.RegisterInput{Email: "valid@mail.ru", Password: "1234", Role: "client"},
			setupMock: func(m *authRepoMock) {
				m.RegisterFunc = func(ctx context.Context, u *entity.User) (*entity.User, error) {
					u.ID = 7
					u.CreatedAt = now
					return u, nil
				}
			},
			checkResult: func(t *testing.T, res *entity.User) {
				if res.ID != 7 { t.Errorf("Ожидался ID 7, получили %d", res.ID) }
			},
		},
		{
			name:        "Failure - Email уже зарегистрирован",
			input:       dto.RegisterInput{Email: "existing@mail.ru", Password: "1234", Role: "client"},
			setupMock:   func(m *authRepoMock) { m.RegisterFunc = func(context.Context, *entity.User) (*entity.User, error) { return nil, fmt.Errorf("user with this email already exists") } },
			wantErr:     true,
			expectedErr: "user with this email already exists",
		},
		{
			name:        "Failure - Сбой базы данных",
			input:       dto.RegisterInput{Email: "db_error@mail.ru", Password: "1234", Role: "client"},
			setupMock:   func(m *authRepoMock) { m.RegisterFunc = func(context.Context, *entity.User) (*entity.User, error) { return nil, fmt.Errorf("connection refused") } },
			wantErr:     true,
			expectedErr: "connection refused",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &authRepoMock{}
			if tc.setupMock != nil { tc.setupMock(mockRepo) }

			service := NewAuthService([]byte("secret"), mockRepo)
			result, err := service.Register(context.Background(), &tc.input)

			if tc.wantErr {
				if err == nil || err.Error() != tc.expectedErr { t.Errorf("Ожидалась ошибка %q, получили %v", tc.expectedErr, err) }
				return
			}
			if err != nil { t.Fatalf("Не ожидали ошибку: %v", err) }
			if tc.checkResult != nil { tc.checkResult(t, result) }
		})
	}
}

func TestLoginUser(t *testing.T) {
	tests := []struct {
		name        string
		input       dto.LoginInput
		setupMock   func(m *authRepoMock)
		wantErr     bool
		expectedErr string
		checkResult func(t *testing.T, res *dto.LoginResponse)
	}{
		{
			name:  "Success - Успешная авторизация",
			input: dto.LoginInput{Email: "test@mail.ru", Password: "1234"},
			setupMock: func(m *authRepoMock) {
				m.GetByEmailFunc = func(ctx context.Context, email string) (*entity.User, error) {
					return &entity.User{ID: 1, Role: entity.RoleModerator, PasswordHash: testPasswordHash}, nil
				}
			},
			checkResult: func(t *testing.T, res *dto.LoginResponse) {
				if res.Role != entity.RoleModerator { t.Errorf("Ожидался Role MODERATOR, получили %s", res.Role) }
				if res.Token == "" { t.Error("Ожидался заполненный токен") }
			},
		},
		{
			name:        "Failure - Неправильный Email",
			input:       dto.LoginInput{Email: "test@mail.ru", Password: "1234"},
			setupMock:   func(m *authRepoMock) { m.GetByEmailFunc = func(context.Context, string) (*entity.User, error) { return nil, fmt.Errorf("invalid email or password") } },
			wantErr:     true,
			expectedErr: "invalid email or password",
		},
		{
			name:  "Failure - Неверный пароль",
			input: dto.LoginInput{Email: "test@mail.ru", Password: "wrong_password"},
			setupMock: func(m *authRepoMock) {
				m.GetByEmailFunc = func(ctx context.Context, email string) (*entity.User, error) {
					return &entity.User{ID: 1, Role: entity.RoleModerator, PasswordHash: testPasswordHash}, nil
				}
			},
			wantErr:     true,
			expectedErr: "invalid email or password",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &authRepoMock{}
			if tc.setupMock != nil { tc.setupMock(mockRepo) }

			service := NewAuthService([]byte("secret"), mockRepo)
			result, err := service.Login(context.Background(), &tc.input)

			if tc.wantErr {
				if err == nil || err.Error() != tc.expectedErr { t.Errorf("Ожидалась ошибка %q, получили %v", tc.expectedErr, err) }
				return
			}
			if err != nil { t.Fatalf("Не ожидали ошибку: %v", err) }
			if tc.checkResult != nil { tc.checkResult(t, result) }
		})
	}
}
