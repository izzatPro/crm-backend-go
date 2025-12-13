package utils

import (
	"testing"
)

func TestAuthorizeUser(t *testing.T) {
	tests := []struct {
		name        string
		userRole    string
		allowedRoles []string
		want        bool
		wantErr     bool
	}{
		{
			name:        "authorized user",
			userRole:    "admin",
			allowedRoles: []string{"admin", "user"},
			want:        true,
			wantErr:     false,
		},
		{
			name:        "unauthorized user",
			userRole:    "guest",
			allowedRoles: []string{"admin", "user"},
			want:        false,
			wantErr:     true,
		},
		{
			name:        "single allowed role match",
			userRole:    "admin",
			allowedRoles: []string{"admin"},
			want:        true,
			wantErr:     false,
		},
		{
			name:        "empty allowed roles",
			userRole:    "admin",
			allowedRoles: []string{},
			want:        false,
			wantErr:     true,
		},
		{
			name:        "empty user role",
			userRole:    "",
			allowedRoles: []string{"admin"},
			want:        false,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AuthorizeUser(tt.userRole, tt.allowedRoles...)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthorizeUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AuthorizeUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

