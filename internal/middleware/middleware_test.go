package middleware

import (
	"avito-tech/internal/entity"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


func TestAuthMiddleware(t *testing.T) {
	jwtKey := []byte("secret_key")

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &entity.CustomClaims{
		UserID: 100,
		Role:   entity.RoleClient,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtKey)

	tests := []struct {
		name         string
		authHeader   string
		expectedCode int
		checkContext func(t *testing.T, ctx context.Context)
	}{
		{
			name:         "Success - Валидный токен клиента",
			authHeader:   "Bearer " + tokenString, // Подставляем динамический токен!
			expectedCode: http.StatusOK,
			checkContext: func(t *testing.T, ctx context.Context) {
				role, ok := ctx.Value(RoleKey).(entity.UserRole)
				if !ok || role != entity.RoleClient {
					t.Errorf("Ожидали роль client, получили: %v", role)
				}
			},
		},
		{
			name:         "Failed - Токен отсутствует",
			authHeader:   "",
			expectedCode: http.StatusUnauthorized,
		},
	}
   	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/house/create", nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}
			rec := httptest.NewRecorder()

			var capturedContext context.Context

			finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				capturedContext = r.Context()
				w.WriteHeader(http.StatusOK) 
			})

			testHandler := AuthMiddleware(jwtKey, finalHandler)

			testHandler.ServeHTTP(rec, req)

			if rec.Code != tc.expectedCode {
				t.Errorf("Ожидали статус-код %d, получили %d", tc.expectedCode, rec.Code)
			}

			if rec.Code == http.StatusOK && tc.checkContext != nil {
				tc.checkContext(t, capturedContext)
			}
		})
	}

}