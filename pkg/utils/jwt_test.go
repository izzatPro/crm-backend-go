package utils

import (
	"os"
	"testing"
)

func TestSignToken(t *testing.T) {
	// Устанавливаем тестовые переменные окружения
	originalSecret := os.Getenv("JWT_SECRET")
	originalExpires := os.Getenv("JWT_EXPIRES_IN")
	
	os.Setenv("JWT_SECRET", "test-secret-key-for-jwt-signing")
	os.Setenv("JWT_EXPIRES_IN", "1h")
	defer func() {
		if originalSecret != "" {
			os.Setenv("JWT_SECRET", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
		if originalExpires != "" {
			os.Setenv("JWT_EXPIRES_IN", originalExpires)
		} else {
			os.Unsetenv("JWT_EXPIRES_IN")
		}
	}()

	tests := []struct {
		name     string
		userId   int
		username string
		role     string
		wantErr  bool
	}{
		{
			name:     "valid token",
			userId:   1,
			username: "testuser",
			role:     "admin",
			wantErr:  false,
		},
		{
			name:     "empty username",
			userId:   2,
			username: "",
			role:     "user",
			wantErr:  false,
		},
		{
			name:     "empty role",
			userId:   3,
			username: "user",
			role:     "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := SignToken(tt.userId, tt.username, tt.role)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && token == "" {
				t.Errorf("SignToken() returned empty token")
			}
		})
	}
}

func TestSignTokenDefaultExpiration(t *testing.T) {
	originalSecret := os.Getenv("JWT_SECRET")
	originalExpires := os.Getenv("JWT_EXPIRES_IN")
	
	os.Setenv("JWT_SECRET", "test-secret-key")
	os.Unsetenv("JWT_EXPIRES_IN")
	defer func() {
		if originalSecret != "" {
			os.Setenv("JWT_SECRET", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
		if originalExpires != "" {
			os.Setenv("JWT_EXPIRES_IN", originalExpires)
		}
	}()

	token, err := SignToken(1, "testuser", "admin")
	if err != nil {
		t.Fatalf("SignToken() failed: %v", err)
	}

	if token == "" {
		t.Errorf("SignToken() returned empty token")
	}
}

func TestSignTokenDifferentUsers(t *testing.T) {
	originalSecret := os.Getenv("JWT_SECRET")
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer func() {
		if originalSecret != "" {
			os.Setenv("JWT_SECRET", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
	}()

	token1, err1 := SignToken(1, "user1", "admin")
	if err1 != nil {
		t.Fatalf("SignToken() failed: %v", err1)
	}

	token2, err2 := SignToken(2, "user2", "user")
	if err2 != nil {
		t.Fatalf("SignToken() failed: %v", err2)
	}

	if token1 == token2 {
		t.Errorf("SignToken() returned same token for different users")
	}
}

