package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecurityHeaders(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	handler := SecurityHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	// Проверяем наличие security headers
	expectedHeaders := map[string]string{
		"X-Frame-Options":                "DENY",
		"X-XSS-Protection":               "1; mode=block",
		"X-Content-Type-Options":         "nosniff",
		"Strict-Transport-Security":     "max-age=63072000; includeSubDomains; preload",
		"Content-Security-Policy":         "default-src 'self'",
		"Referrer-Policy":                "no-referrer",
		"X-Permitted-Cross-Domain-Policies": "none",
		"Cache-Control":                  "no-store, no-cache, must-revalidate, max-age=0",
		"Cross-Origin-Resource-Policy":   "same-origin",
		"Cross-Origin-Opener-Policy":     "same-origin",
		"Cross-Origin-Embedder-Policy":  "require-corp",
	}

	for header, expectedValue := range expectedHeaders {
		actualValue := rr.Header().Get(header)
		if actualValue != expectedValue {
			t.Errorf("SecurityHeaders() %s = %v, want %v", header, actualValue, expectedValue)
		}
	}

	// Проверяем, что Server header пустой
	if rr.Header().Get("Server") != "" {
		t.Errorf("Server header should be empty, got %s", rr.Header().Get("Server"))
	}
}

