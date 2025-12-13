package middlewares

import (
	"net/http"
	"net/http/httptest"
	"os"
	"restapi/pkg/utils"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestJWTMiddleware(t *testing.T) {
	// Устанавливаем тестовый JWT_SECRET
	originalSecret := os.Getenv("JWT_SECRET")
	os.Setenv("JWT_SECRET", "test-secret-key-for-jwt-middleware")
	defer func() {
		if originalSecret != "" {
			os.Setenv("JWT_SECRET", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
	}()

	tests := []struct {
		name           string
		setCookie      bool
		tokenValid     bool
		expectedStatus int
	}{
		{
			name:           "valid token",
			setCookie:     true,
			tokenValid:     true,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing cookie",
			setCookie:      false,
			tokenValid:     false,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid token",
			setCookie:      true,
			tokenValid:     false,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)

			if tt.setCookie {
				var tokenString string
				if tt.tokenValid {
					// Создаем валидный токен
					token, err := utils.SignToken(1, "testuser", "admin")
					if err != nil {
						t.Fatalf("Failed to create token: %v", err)
					}
					tokenString = token
				} else {
					// Создаем невалидный токен
					tokenString = "invalid.token.here"
				}
				req.AddCookie(&http.Cookie{
					Name:  "Bearer",
					Value: tokenString,
				})
			}

			rr := httptest.NewRecorder()
			handler := JWTMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Проверяем, что контекст установлен
				role := r.Context().Value(utils.ContextKey("role"))
				if role == nil && tt.tokenValid {
					t.Error("Role not set in context")
				}
				w.WriteHeader(http.StatusOK)
			}))

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("JWTMiddleware() status = %v, want %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestJWTMiddlewareExpiredToken(t *testing.T) {
	originalSecret := os.Getenv("JWT_SECRET")
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer func() {
		if originalSecret != "" {
			os.Setenv("JWT_SECRET", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
	}()

	// Создаем истекший токен
	expiredTime := time.Now().Add(-24 * time.Hour)
	claims := jwt.MapClaims{
		"uid":  1,
		"user": "testuser",
		"role": "admin",
		"exp":  jwt.NewNumericDate(expiredTime),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("test-secret-key"))

	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  "Bearer",
		Value: tokenString,
	})

	rr := httptest.NewRecorder()
	handler := JWTMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("JWTMiddleware() should reject expired token, got status %v", status)
	}
}

