// Package router provides HTTP route definitions and configuration.
package router

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/rei0721/go-scaffold/internal/handler"
	"github.com/rei0721/go-scaffold/internal/middleware"
	"github.com/rei0721/go-scaffold/pkg/logger"
	"github.com/rei0721/go-scaffold/types/result"
)

func TestHealthEndpoint(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a minimal logger
	log := logger.Default()

	// Create router with nil handler (health endpoint doesn't need it)
	r := New(&handler.UserHandler{}, &handler.RBACHandler{}, log, nil, nil)

	// Setup router with default middleware config
	cfg := middleware.DefaultMiddlewareConfig()
	engine := r.Setup(cfg)

	// Create test request
	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	// Create response recorder
	w := httptest.NewRecorder()

	// Serve the request
	engine.ServeHTTP(w, req)

	// Verify status code
	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify response body
	var resp result.Result[map[string]interface{}]
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Code != 0 {
		t.Errorf("expected code 0, got %d", resp.Code)
	}

	if resp.Data["status"] != "ok" {
		t.Errorf("expected status 'ok', got %v", resp.Data["status"])
	}

	// Verify X-Request-ID header is present
	if w.Header().Get("X-Request-ID") == "" {
		t.Error("expected X-Request-ID header to be present")
	}
}
