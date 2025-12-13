package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	tests := []struct {
		name           string
		requests       int
		expectedStatus int
	}{
		{
			name:           "within limit",
			requests:       3,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "at limit",
			requests:       5,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "over limit",
			requests:       6,
			expectedStatus: http.StatusTooManyRequests,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем новый лимитер для каждого теста
			testLimiter := NewRateLimiter(5, 1*time.Minute)
			
			var lastStatus int
			for i := 0; i < tt.requests; i++ {
				req := httptest.NewRequest("GET", "/test", nil)
				req.RemoteAddr = "127.0.0.1:12345"
				rr := httptest.NewRecorder()

				handler := testLimiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))

				handler.ServeHTTP(rr, req)
				lastStatus = rr.Code
			}

			if lastStatus != tt.expectedStatus {
				t.Errorf("RateLimiter() after %d requests status = %v, want %v", 
					tt.requests, lastStatus, tt.expectedStatus)
			}
		})
	}
}

func TestRateLimiterDifferentIPs(t *testing.T) {
	limiter := NewRateLimiter(2, 1*time.Minute)

	// Первый IP делает запросы
	req1 := httptest.NewRequest("GET", "/test", nil)
	req1.RemoteAddr = "127.0.0.1:12345"
	rr1 := httptest.NewRecorder()

	handler := limiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Делаем 2 запроса с первого IP (в пределах лимита)
	handler.ServeHTTP(rr1, req1)
	if rr1.Code != http.StatusOK {
		t.Errorf("First request should succeed, got %v", rr1.Code)
	}

	// Второй IP делает запрос
	req2 := httptest.NewRequest("GET", "/test", nil)
	req2.RemoteAddr = "127.0.0.2:12345"
	rr2 := httptest.NewRecorder()

	handler.ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusOK {
		t.Errorf("Request from different IP should succeed, got %v", rr2.Code)
	}
}

