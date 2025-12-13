package utils

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name    string
		password string
		wantErr bool
	}{
		{
			name:     "valid password",
			password: "testpassword123",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  true,
		},
		{
			name:     "long password",
			password: "verylongpasswordthatisoveronehundredcharacterslongandshouldstillworkcorrectlywithoutanyissues",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HashPassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == "" {
				t.Errorf("HashPassword() returned empty string")
			}
			if !tt.wantErr {
				// Проверяем формат: salt.hash
				parts := got
				if len(parts) < 10 {
					t.Errorf("HashPassword() returned hash that is too short")
				}
			}
		})
	}
}

func TestVerifyPassword(t *testing.T) {
	password := "testpassword123"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() failed: %v", err)
	}

	tests := []struct {
		name          string
		password      string
		encodedHash   string
		wantErr       bool
	}{
		{
			name:        "correct password",
			password:    password,
			encodedHash: hashedPassword,
			wantErr:     false,
		},
		{
			name:        "incorrect password",
			password:    "wrongpassword",
			encodedHash: hashedPassword,
			wantErr:     true,
		},
		{
			name:        "invalid hash format",
			password:    password,
			encodedHash: "invalidhash",
			wantErr:     true,
		},
		{
			name:        "empty hash",
			password:    password,
			encodedHash: "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := VerifyPassword(tt.password, tt.encodedHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHashPasswordConsistency(t *testing.T) {
	password := "testpassword123"
	
	// Хешируем один и тот же пароль дважды
	hash1, err1 := HashPassword(password)
	if err1 != nil {
		t.Fatalf("HashPassword() failed: %v", err1)
	}

	hash2, err2 := HashPassword(password)
	if err2 != nil {
		t.Fatalf("HashPassword() failed: %v", err2)
	}

	// Хеши должны быть разными (из-за случайной соли)
	if hash1 == hash2 {
		t.Errorf("HashPassword() returned same hash for same password (should have different salts)")
	}

	// Но оба должны верифицироваться правильно
	if err := VerifyPassword(password, hash1); err != nil {
		t.Errorf("VerifyPassword() failed for first hash: %v", err)
	}

	if err := VerifyPassword(password, hash2); err != nil {
		t.Errorf("VerifyPassword() failed for second hash: %v", err)
	}
}

