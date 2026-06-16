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
	config.SetRuntimeTimeouts(123000, 45000)
	t.Cleanup(func() { config.SetRuntimeTimeouts(300000, 60000) })

	w := performSettingsJSON(GetCircuitBreaker(func() metrics.CircuitBreakerParams {
		return metrics.CircuitBreakerParams{
			WindowSize:                   10,
			FailureThreshold:             0.5,
			ConsecutiveFailuresThreshold: 3,
			StreamFirstContentTimeoutMs:  30000,
			StreamInactivityTimeoutMs:    20000,
			StreamToolCallIdleTimeoutMs:  30000,
		}
	}, &config.EnvConfig{RequestTimeout: 300000, ResponseHeaderTimeout: 60}), http.MethodGet, "")

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
	if got := int(body["requestTimeoutMs"].(float64)); got != 123000 {
		t.Fatalf("requestTimeoutMs = %d, want 123000", got)
	}
	if got := int(body["responseHeaderTimeoutMs"].(float64)); got != 45000 {
		t.Fatalf("responseHeaderTimeoutMs = %d, want 45000", got)
	}
}

func TestSetCircuitBreaker_AcceptsRequestLifecycleTimeouts(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfgManager := newSettingsTestConfigManager(t)

	w := performSettingsJSON(SetCircuitBreaker(cfgManager), http.MethodPut, `{"requestTimeoutMs":300000,"responseHeaderTimeoutMs":180000}`)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", w.Code, http.StatusOK, w.Body.String())
	}

	cfg := cfgManager.GetCircuitBreakerConfig()
	if cfg.RequestTimeoutMs == nil || *cfg.RequestTimeoutMs != 300000 {
		t.Fatalf("saved requestTimeoutMs = %v, want 300000", cfg.RequestTimeoutMs)
	}
	if cfg.ResponseHeaderTimeoutMs == nil || *cfg.ResponseHeaderTimeoutMs != 180000 {
		t.Fatalf("saved responseHeaderTimeoutMs = %v, want 180000", cfg.ResponseHeaderTimeoutMs)
	}
}

func TestSetCircuitBreaker_RejectsInvalidRequestLifecycleTimeouts(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfgManager := newSettingsTestConfigManager(t)

	w := performSettingsJSON(SetCircuitBreaker(cfgManager), http.MethodPut, `{"requestTimeoutMs":300001}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("requestTimeoutMs status = %d, want %d", w.Code, http.StatusBadRequest)
	}
	if !strings.Contains(w.Body.String(), "requestTimeoutMs") {
		t.Fatalf("response body %q should mention requestTimeoutMs", w.Body.String())
	}

	w = performSettingsJSON(SetCircuitBreaker(cfgManager), http.MethodPut, `{"responseHeaderTimeoutMs":999}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("responseHeaderTimeoutMs status = %d, want %d", w.Code, http.StatusBadRequest)
	}
	if !strings.Contains(w.Body.String(), "responseHeaderTimeoutMs") {
		t.Fatalf("response body %q should mention responseHeaderTimeoutMs", w.Body.String())
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

func TestGetHistoricalImageTurnLimit_ReturnsValue(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfgManager := newSettingsTestConfigManager(t)
	if err := cfgManager.SetHistoricalImageTurnLimit(3); err != nil {
		t.Fatalf("set limit failed: %v", err)
	}

	w := performSettingsJSON(GetHistoricalImageTurnLimit(cfgManager), http.MethodGet, "")
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}
	if got := int(body["historicalImageTurnLimit"].(float64)); got != 3 {
		t.Fatalf("historicalImageTurnLimit = %d, want 3", got)
	}
}

func TestSetHistoricalImageTurnLimit_Accepts(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfgManager := newSettingsTestConfigManager(t)

	w := performSettingsJSON(SetHistoricalImageTurnLimit(cfgManager), http.MethodPut, `{"limit":5}`)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", w.Code, http.StatusOK, w.Body.String())
	}
	if got := cfgManager.GetHistoricalImageTurnLimit(); got != 5 {
		t.Fatalf("saved limit = %d, want 5", got)
	}
}

func TestSetHistoricalImageTurnLimit_ZeroNormalizesToDefault(t *testing.T) {
	// 全局功能始终开启：limit=0 归一为默认值 5（不再表示"不限制"）
	gin.SetMode(gin.TestMode)
	cfgManager := newSettingsTestConfigManager(t)
	_ = cfgManager.SetHistoricalImageTurnLimit(10)

	w := performSettingsJSON(SetHistoricalImageTurnLimit(cfgManager), http.MethodPut, `{"limit":0}`)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", w.Code, http.StatusOK, w.Body.String())
	}
	if got := cfgManager.GetHistoricalImageTurnLimit(); got != 5 {
		t.Fatalf("saved limit = %d, want 5 (default)", got)
	}
}

func TestSetHistoricalImageTurnLimit_BelowMinNormalizesToMin(t *testing.T) {
	// 0 < limit < 3 归一为最低值 3
	gin.SetMode(gin.TestMode)
	cfgManager := newSettingsTestConfigManager(t)

	w := performSettingsJSON(SetHistoricalImageTurnLimit(cfgManager), http.MethodPut, `{"limit":1}`)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", w.Code, http.StatusOK, w.Body.String())
	}
	if got := cfgManager.GetHistoricalImageTurnLimit(); got != 3 {
		t.Fatalf("saved limit = %d, want 3 (min)", got)
	}
}

func TestSetHistoricalImageTurnLimit_NegativeNormalizesToDefault(t *testing.T) {
	// 负数视为无效，归一为默认值 5
	gin.SetMode(gin.TestMode)
	cfgManager := newSettingsTestConfigManager(t)

	w := performSettingsJSON(SetHistoricalImageTurnLimit(cfgManager), http.MethodPut, `{"limit":-1}`)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", w.Code, http.StatusOK, w.Body.String())
	}
	if got := cfgManager.GetHistoricalImageTurnLimit(); got != 5 {
		t.Fatalf("saved limit = %d, want 5 (default)", got)
	}
}
