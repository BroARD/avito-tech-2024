package dto

import "avito-tech/internal/entity"

type RegisterInput struct {
	Email    string          `json:"email"`
	Password string          `json:"password"`
	Role     entity.UserRole `json:"role"`
}

type RegisterResponse struct {
	Email string          `json:"email"`
	Role  entity.UserRole `json:"role"`
}

type LoginResponse struct {
	Token string          `json:"token"`
	Role  entity.UserRole `json:"role"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
