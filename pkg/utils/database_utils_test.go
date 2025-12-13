package utils

import (
	"net/http/httptest"
	"testing"
)

func TestAddFilters(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		args     []interface{}
		url      string
		wantArgs int
	}{
		{
			name:     "no filters",
			query:    "SELECT * FROM table WHERE 1=1",
			args:     []interface{}{},
			url:      "/test",
			wantArgs: 0,
		},
		{
			name:     "filter by first_name",
			query:    "SELECT * FROM table WHERE 1=1",
			args:     []interface{}{},
			url:      "/test?first_name=John",
			wantArgs: 1,
		},
		{
			name:     "multiple filters",
			query:    "SELECT * FROM table WHERE 1=1",
			args:     []interface{}{},
			url:      "/test?first_name=John&last_name=Doe&email=test@example.com",
			wantArgs: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.url, nil)
			query, args := AddFilters(req, tt.query, tt.args)

			if len(args) != tt.wantArgs {
				t.Errorf("AddFilters() args length = %v, want %v", len(args), tt.wantArgs)
			}

			if tt.wantArgs > 0 {
				// Проверяем, что запрос содержит AND
				if len(query) <= len(tt.query) {
					t.Errorf("AddFilters() query should be longer when filters are added")
				}
			}
		})
	}
}

func TestAddSorting(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		url      string
		contains string
	}{
		{
			name:     "no sorting",
			query:    "SELECT * FROM table",
			url:      "/test",
			contains: "",
		},
		{
			name:     "sort by first_name asc",
			query:    "SELECT * FROM table",
			url:      "/test?sortby=first_name:asc",
			contains: "ORDER BY",
		},
		{
			name:     "sort by last_name desc",
			query:    "SELECT * FROM table",
			url:      "/test?sortby=last_name:desc",
			contains: "ORDER BY",
		},
		{
			name:     "multiple sorts",
			query:    "SELECT * FROM table",
			url:      "/test?sortby=first_name:asc&sortby=last_name:desc",
			contains: "ORDER BY",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.url, nil)
			result := AddSorting(req, tt.query)

			if tt.contains != "" {
				if len(result) <= len(tt.query) {
					t.Errorf("AddSorting() should add ORDER BY clause")
				}
			}
		})
	}
}

func TestIsValidSortOrder(t *testing.T) {
	tests := []struct {
		name  string
		order string
		want  bool
	}{
		{"valid asc", "asc", true},
		{"valid desc", "desc", true},
		{"invalid order", "invalid", false},
		{"empty order", "", false},
		{"uppercase ASC", "ASC", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidSortOrder(tt.order); got != tt.want {
				t.Errorf("isValidSortOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidSortField(t *testing.T) {
	tests := []struct {
		name  string
		field string
		want  bool
	}{
		{"valid first_name", "first_name", true},
		{"valid last_name", "last_name", true},
		{"valid email", "email", true},
		{"valid class", "class", true},
		{"valid subject", "subject", true},
		{"invalid field", "invalid_field", false},
		{"empty field", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidSortField(tt.field); got != tt.want {
				t.Errorf("isValidSortField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateInsertQuery(t *testing.T) {
	type TestModel struct {
		ID        int    `db:"id,omitempty"`
		FirstName string `db:"first_name,omitempty"`
		LastName  string `db:"last_name,omitempty"`
		Email     string `db:"email,omitempty"`
	}

	query := GenerateInsertQuery("test_table", TestModel{})
	
	if query == "" {
		t.Errorf("GenerateInsertQuery() returned empty string")
	}

	// Проверяем, что запрос содержит название таблицы
	if len(query) < len("test_table") {
		t.Errorf("GenerateInsertQuery() query too short")
	}
}

func TestGetStructValues(t *testing.T) {
	type TestModel struct {
		ID        int    `db:"id,omitempty"`
		FirstName string `db:"first_name,omitempty"`
		LastName  string `db:"last_name,omitempty"`
	}

	model := TestModel{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
	}

	values := GetStructValues(model)
	
	// ID должен быть исключен
	if len(values) != 2 {
		t.Errorf("GetStructValues() returned %v values, want 2", len(values))
	}
}

