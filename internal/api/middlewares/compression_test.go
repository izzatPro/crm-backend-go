package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCompression(t *testing.T) {
	tests := []struct {
		name              string
		acceptEncoding    string
		shouldCompress    bool
		checkContentEncoding bool
	}{
		{
			name:              "gzip accepted",
			acceptEncoding:    "gzip",
			shouldCompress:    true,
			checkContentEncoding: true,
		},
		{
			name:              "no gzip",
			acceptEncoding:    "deflate",
			shouldCompress:    false,
			checkContentEncoding: false,
		},
		{
			name:              "empty accept encoding",
			acceptEncoding:    "",
			shouldCompress:    false,
			checkContentEncoding: false,
		},
		{
			name:              "gzip in multiple encodings",
			acceptEncoding:    "gzip, deflate, br",
			shouldCompress:    true,
			checkContentEncoding: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.acceptEncoding != "" {
				req.Header.Set("Accept-Encoding", tt.acceptEncoding)
			}

			rr := httptest.NewRecorder()
			handler := Compression(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("test response"))
			}))

			handler.ServeHTTP(rr, req)

			if tt.checkContentEncoding {
				if rr.Header().Get("Content-Encoding") != "gzip" {
					t.Errorf("Content-Encoding header not set to gzip")
				}

				// Проверяем, что данные сжаты
				reader, err := gzip.NewReader(rr.Body)
				if err != nil {
					t.Fatalf("Failed to create gzip reader: %v", err)
				}
				defer reader.Close()

				decompressed, err := io.ReadAll(reader)
				if err != nil {
					t.Fatalf("Failed to decompress: %v", err)
				}

				if string(decompressed) != "test response" {
					t.Errorf("Decompressed content mismatch: got %s, want test response", string(decompressed))
				}
			} else {
				if rr.Header().Get("Content-Encoding") == "gzip" {
					t.Errorf("Content should not be compressed")
				}
			}
		})
	}
}

func TestGzipResponseWriter(t *testing.T) {
	rr := httptest.NewRecorder()
	gz := gzip.NewWriter(rr)
	
	grw := &gzipResponseWriter{
		ResponseWriter: rr,
		Writer:         gz,
	}

	testData := []byte("test data")
	n, err := grw.Write(testData)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	if n != len(testData) {
		t.Errorf("Write returned %d bytes, want %d", n, len(testData))
	}

	gz.Close()
}

