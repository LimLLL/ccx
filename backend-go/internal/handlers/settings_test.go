package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BenedictKing/ccx/internal/config"
	"github.com/BenedictKing/ccx/internal/metrics"
	"github.com/gin-gonic/gin"
)

func newSettingsTestConfigManager(t *testing.T) *config.ConfigManager {
	t.Helper()

	configPath := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(configPath, []byte(`{"upstream":[]}`), 0644); err != nil {
		t.Fatalf("写入测试配置失败: %v", err)
	}

	cfgManager, err := config.NewConfigManager(configPath, "")
	if err != nil {
		t.Fatalf("初始化配置管理器失败: %v", err)
	}
	t.Cleanup(func() { _ = cfgManager.Close() })
	return cfgManager
}

func performSettingsJSON(handler gin.HandlerFunc, method string, body string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(method, "/api/settings/circuit-breaker", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	handler(c)
	return w
}

func TestGetCircuitBreaker_ReturnsToolCallIdleTimeout(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := performSettingsJSON(GetCircuitBreaker(func() metrics.CircuitBreakerParams {
		return metrics.CircuitBreakerParams{
			WindowSize:                   10,
			FailureThreshold:             0.5,
			ConsecutiveFailuresThreshold: 3,
			StreamFirstContentTimeoutMs:  30000,
			StreamInactivityTimeoutMs:    20000,
			StreamToolCallIdleTimeoutMs:  30000,
		}
	}), http.MethodGet, "")

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}
	if got := int(body["streamToolCallIdleTimeoutMs"].(float64)); got != 30000 {
		t.Fatalf("streamToolCallIdleTimeoutMs = %d, want 30000", got)
	}
}

func TestSetCircuitBreaker_AcceptsToolCallIdleTimeout(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfgManager := newSettingsTestConfigManager(t)

	w := performSettingsJSON(SetCircuitBreaker(cfgManager), http.MethodPut, `{"streamToolCallIdleTimeoutMs":300000}`)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", w.Code, http.StatusOK, w.Body.String())
	}

	value := cfgManager.GetCircuitBreakerConfig().StreamToolCallIdleTimeoutMs
	if value == nil || *value != 300000 {
		t.Fatalf("saved streamToolCallIdleTimeoutMs = %v, want 300000", value)
	}
}

func TestSetCircuitBreaker_AcceptsInactivityTimeout(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfgManager := newSettingsTestConfigManager(t)

	w := performSettingsJSON(SetCircuitBreaker(cfgManager), http.MethodPut, `{"streamInactivityTimeoutMs":180000}`)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", w.Code, http.StatusOK, w.Body.String())
	}

	value := cfgManager.GetCircuitBreakerConfig().StreamInactivityTimeoutMs
	if value == nil || *value != 180000 {
		t.Fatalf("saved streamInactivityTimeoutMs = %v, want 180000", value)
	}
}

func TestSetCircuitBreaker_RejectsInvalidToolCallIdleTimeout(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfgManager := newSettingsTestConfigManager(t)

	w := performSettingsJSON(SetCircuitBreaker(cfgManager), http.MethodPut, `{"streamToolCallIdleTimeoutMs":29999}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
	if !strings.Contains(w.Body.String(), "streamToolCallIdleTimeoutMs") {
		t.Fatalf("response body %q should mention streamToolCallIdleTimeoutMs", w.Body.String())
	}
}
