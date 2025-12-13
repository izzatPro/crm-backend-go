package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCors(t *testing.T) {
	tests := []struct {
		name           string
		origin         string
		expectedStatus int
	}{
		{
			name:           "allowed origin",
			origin:         "https://my-origin-url.com",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "another allowed origin",
			origin:         "https://www.myfrontend.com",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "disallowed origin",
			origin:         "https://evil.com",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "empty origin",
			origin:         "",
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}

			rr := httptest.NewRecorder()
			handler := Cors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK && tt.origin != "" {
				if rr.Header().Get("Access-Control-Allow-Origin") != tt.origin {
					t.Errorf("Access-Control-Allow-Origin header not set correctly")
				}
			}
		})
	}
}

func TestCorsPreflight(t *testing.T) {
	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "https://my-origin-url.com")

	rr := httptest.NewRecorder()
	handler := Cors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("preflight request returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Проверяем CORS заголовки
	if rr.Header().Get("Access-Control-Allow-Methods") == "" {
		t.Errorf("Access-Control-Allow-Methods header not set")
	}
}

func TestIsOriginAllowed(t *testing.T) {
	tests := []struct {
		name   string
		origin string
		want   bool
	}{
		{"allowed origin 1", "https://my-origin-url.com", true},
		{"allowed origin 2", "https://www.myfrontend.com", true},
		{"allowed origin 3", "https://localhost:3000", true},
		{"disallowed origin", "https://evil.com", false},
		{"empty origin", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isOriginAllowed(tt.origin); got != tt.want {
				t.Errorf("isOriginAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}

