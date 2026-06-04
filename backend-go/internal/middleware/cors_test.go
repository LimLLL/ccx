package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BenedictKing/ccx/internal/config"
	"github.com/gin-gonic/gin"
)

func setupRouterWithCORS(envCfg *config.EnvConfig) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(CORSMiddleware(envCfg))
	r.Use(WebAuthMiddleware(envCfg, nil))
	r.GET("/api/channels", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	return r
}

func TestCORSMiddleware_PreflightBypassesAuthWhenCORSDisabled(t *testing.T) {
	envCfg := &config.EnvConfig{
		EnableCORS:     false,
		EnableWebUI:    true,
		ProxyAccessKey: "test-key",
		CORSOrigin:     "*",
	}
	r := setupRouterWithCORS(envCfg)

	req := httptest.NewRequest(http.MethodOptions, "/api/channels", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Access-Control-Request-Method", http.MethodGet)
	req.Header.Set("Access-Control-Request-Headers", "x-api-key")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Fatalf("Access-Control-Allow-Origin = %q, want request origin", got)
	}
	if got := w.Header().Get("Access-Control-Allow-Headers"); got == "" {
		t.Fatal("Access-Control-Allow-Headers should be set")
	}
}

func TestCORSMiddleware_ActualResponseIncludesHeadersWhenCORSDisabled(t *testing.T) {
	envCfg := &config.EnvConfig{
		EnableCORS:     false,
		EnableWebUI:    true,
		ProxyAccessKey: "test-key",
		CORSOrigin:     "*",
	}
	r := setupRouterWithCORS(envCfg)

	req := httptest.NewRequest(http.MethodGet, "/api/channels", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("x-api-key", "test-key")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Fatalf("Access-Control-Allow-Origin = %q, want request origin", got)
	}
}
