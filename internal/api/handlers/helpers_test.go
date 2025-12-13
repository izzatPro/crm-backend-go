package handlers

import (
	"restapi/internal/models"
	"testing"
)

func TestCheckBlankFields(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{
			name: "all required fields filled",
			value: models.Student{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
				Class:     "10A",
			},
			wantErr: false,
		},
		{
			name: "empty first name",
			value: models.Exec{
				FirstName: "",
				LastName:  "Doe",
				Email:     "john@example.com",
			},
			wantErr: true,
		},
		{
			name: "empty last name",
			value: models.Exec{
				FirstName: "John",
				LastName:  "",
				Email:     "john@example.com",
			},
			wantErr: true,
		},
		{
			name: "empty email",
			value: models.Exec{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckBlankFields(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckBlankFields() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetFieldNames(t *testing.T) {
	model := models.Exec{}
	fields := GetFieldNames(model)

	if len(fields) == 0 {
		t.Errorf("GetFieldNames() returned empty slice")
	}

	// Проверяем, что поля содержат ожидаемые значения
	expectedFields := map[string]bool{
		"id":                    true,
		"first_name":            true,
		"last_name":             true,
		"email":                 true,
		"username":              true,
		"password":              true,
		"password_changed_at":   true,
		"user_created_at":       true,
		"password_reset_token":  true,
		"password_token_expires": true,
		"inactive_status":       true,
		"role":                  true,
	}

	for _, field := range fields {
		if !expectedFields[field] {
			t.Errorf("GetFieldNames() returned unexpected field: %s", field)
		}
	}
}

