package auth

import (
	"avito-tech/internal/auth/dto"
	"avito-tech/internal/entity"
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	GenerateDummyToken(ctx context.Context, role entity.UserRole) (string, error)
	Register(ctx context.Context, newUser *dto.RegisterInput) (*entity.User, error)
	Login(ctx context.Context, input *dto.LoginInput) (*dto.LoginResponse, error)
}

type authService struct {
	jwtKey []byte
	repo   Repository
}

// Login implements [Service].
func (s *authService) Login(ctx context.Context, input *dto.LoginInput) (*dto.LoginResponse, error) {
	user, err := s.repo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, errors.New("invalid email or password") 
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		return nil, errors.New("invalid email or password") 
	}

	tokenString, err := s.generateToken(user.ID, user.Role)
	if err != nil {
		return nil, errors.New("invalid email or password") 
	}

	return &dto.LoginResponse{
		Token: tokenString,
		Role:  user.Role,
	}, nil
}

// Register implements [Service].
func (s *authService) Register(ctx context.Context, newUser *dto.RegisterInput) (*entity.User, error) {
	hashPassword, err := hashPassword(newUser.Password)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.Register(ctx, &entity.User{
		Email:        newUser.Email,
		PasswordHash: hashPassword,
		Role:         newUser.Role,
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GenerateDummyToken implements [Service].
func (s *authService) GenerateDummyToken(ctx context.Context, role entity.UserRole) (string, error) {
	var dummyID uint = 100
	if role == entity.RoleModerator {
		dummyID = 200
	}

	tokenString, err := s.generateToken(dummyID, role)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func NewAuthService(jwtKey []byte, repo Repository) Service {
	return &authService{
		jwtKey: jwtKey,
		repo:   repo,
	}
}

func (s *authService) generateToken(userID uint, role entity.UserRole) (string, error) {
	exprationTime := time.Now().Add(24 * time.Hour)

	claims := &entity.CustomClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exprationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.jwtKey)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
