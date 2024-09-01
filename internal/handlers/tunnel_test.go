package handlers_test

import (
	"bytes"
	"go-sentry-tunnel/internal/config"
	"go-sentry-tunnel/internal/handlers"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockHTTPClient struct {
	DoFunc func(URL string, payload []byte, timeout time.Duration) (*http.Response, error)
}

// PostWithTimeout calls the mock DoFunc.
func (m *MockHTTPClient) PostWithTimeout(URL string, payload []byte, timeout time.Duration) (*http.Response, error) {
	return m.DoFunc(URL, payload, timeout)
}

func TestHandleTunnel(t *testing.T) {
	t.Parallel()
	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	cfg := &config.Config{
		AllowOrigins:   []string{"http://allowed-origin.com", "*"},
		DSN:            []string{"valid-dsn"},
		RequestTimeout: 5 * time.Second,
	}

	tests := []struct {
		name           string
		origin         string
		body           string
		upstreamStatus int
		expectedStatus int
	}{
		{
			name:           "Valid request with allowed origin and valid DSN",
			origin:         "http://allowed-origin.com",
			body:           `{"dsn":"valid-dsn"}`,
			upstreamStatus: http.StatusOK,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Request with disallowed origin",
			origin:         "http://disallowed-origin.com",
			body:           `{"dsn":"valid-dsn"}`,
			upstreamStatus: http.StatusOK,
			expectedStatus: http.StatusOK, // CORS headers won't be set, but request will still be processed
		},
		{
			name:           "Request with invalid DSN",
			origin:         "http://allowed-origin.com",
			body:           `{"dsn":"invalid-dsn"}`,
			upstreamStatus: http.StatusOK,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Request with empty DSN",
			origin:         "http://allowed-origin.com",
			body:           `{"dsn":""}`,
			upstreamStatus: http.StatusOK,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Request with invalid JSON payload",
			origin:         "http://allowed-origin.com",
			body:           `{"dsn":`,
			upstreamStatus: http.StatusOK,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Request with upstream server returning non-200 status code",
			origin:         "http://allowed-origin.com",
			body:           `{"dsn":"valid-dsn"}`,
			upstreamStatus: http.StatusInternalServerError,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt
		mockClient := &MockHTTPClient{
			DoFunc: func(URL string, payload []byte, timeout time.Duration) (*http.Response, error) {
				return &http.Response{
					StatusCode: tt.upstreamStatus,
					Body:       io.NopCloser(bytes.NewBufferString("mock response")),
				}, nil
			},
		}
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/tunnel", bytes.NewBufferString(tt.body))
			req.Header.Set("Origin", tt.origin)
			w := httptest.NewRecorder()

			handler := handlers.HandleTunnel(mockClient, logger, cfg)
			handler.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}
