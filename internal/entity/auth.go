package entity

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserRole string

const (
	RoleClient    UserRole = "client"
	RoleModerator UserRole = "moderator"
)

type CustomClaims struct {
	UserID uint
	Role   UserRole
	jwt.RegisteredClaims
}

type User struct {
	ID           uint      `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"` 
	Role         UserRole  `json:"role" db:"role"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"` 
}