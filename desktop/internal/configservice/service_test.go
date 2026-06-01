package configservice

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// ── Part A: 纯函数表驱动测试 ──────────────────────────────

func TestExtractTopLevelTomlString(t *testing.T) {
	cases := []struct {
		name    string
		content string
		key     string
		wantVal string
		wantOK  bool
	}{
		{"正常提取", `model_provider = "ccx"`, "model_provider", "ccx", true},
		{"不存在 key", `model_provider = "ccx"`, "other", "", false},
		{"空内容", "", "key", "", false},
		{"值含特殊字符", `key = "http://127.0.0.1:3688/v1"`, "key", "http://127.0.0.1:3688/v1", true},
		{"带注释", `model_provider = "ccx"  # comment`, "model_provider", "ccx", true},
		{"多行取第一个", "a = \"1\"\na = \"2\"", "a", "1", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, ok := extractTopLevelTomlString(c.content, c.key)
			if ok != c.wantOK {
				t.Fatalf("ok = %v, want %v", ok, c.wantOK)
			}
			if got != c.wantVal {
				t.Errorf("value = %q, want %q", got, c.wantVal)
			}
		})
	}
}

func TestExtractNamedTomlBlock(t *testing.T) {
	cases := []struct {
		name    string
		content string
		table   string
		wantOK  bool
	}{
		{
			"正常提取",
			"[model_providers.ccx]\nname = \"CCX\"\nbase_url = \"http://localhost\"\n\n[model_providers.openai]\nname = \"OpenAI\"\n",
			"model_providers.ccx",
			true,
		},
		{"不存在", "[other]\nkey = \"val\"\n", "model_providers.ccx", false},
		{"空内容", "", "model_providers.ccx", false},
		{
			"最后一个 block",
			"[other]\nkey = \"val\"\n[model_providers.ccx]\nname = \"CCX\"\n",
			"model_providers.ccx",
			true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, ok := extractNamedTomlBlock(c.content, c.table)
			if ok != c.wantOK {
				t.Fatalf("ok = %v, want %v", ok, c.wantOK)
			}
		})
	}
}

func TestFindNamedTomlBlock(t *testing.T) {
	content := "[other]\nkey = \"val\"\n[model_providers.ccx]\nname = \"CCX\"\nbase_url = \"x\"\n\n[model_providers.openai]\nname = \"OpenAI\"\n"
	start, end, ok := findNamedTomlBlock(content, "model_providers.ccx")
	if !ok {
		t.Fatal("expected to find block")
	}
	block := content[start:end]
	if !strings.Contains(block, "name = \"CCX\"") {
		t.Errorf("block does not contain expected content: %q", block)
	}
	if strings.Contains(block, "model_providers.openai") {
		t.Errorf("block should not contain next table")
	}
}

func TestUpsertTopLevelTomlString(t *testing.T) {
	t.Run("替换已有", func(t *testing.T) {
		got := upsertTopLevelTomlString(`model_provider = "openai"`, "model_provider", "ccx")
		if !strings.Contains(got, `"ccx"`) {
			t.Errorf("expected ccx, got %q", got)
		}
		if strings.Contains(got, `"openai"`) {
			t.Errorf("should not contain old value")
		}
	})
	t.Run("新增 key", func(t *testing.T) {
		got := upsertTopLevelTomlString("other = \"val\"\n", "model_provider", "ccx")
		if !strings.Contains(got, `model_provider = "ccx"`) {
			t.Errorf("expected new key, got %q", got)
		}
	})
	t.Run("空内容", func(t *testing.T) {
		got := upsertTopLevelTomlString("", "key", "val")
		if !strings.Contains(got, `key = "val"`) {
			t.Errorf("expected key in empty content, got %q", got)
		}
	})
}

func TestUpsertNamedTomlBlock(t *testing.T) {
	t.Run("替换已有", func(t *testing.T) {
		content := "[model_providers.ccx]\nold = \"data\"\n\n[other]\nk = \"v\"\n"
		block := "[model_providers.ccx]\nnew = \"data\"\n"
		got := upsertNamedTomlBlock(content, "model_providers.ccx", block)
		if !strings.Contains(got, `new = "data"`) {
			t.Errorf("expected new block, got %q", got)
		}
		if strings.Contains(got, `old = "data"`) {
			t.Errorf("should not contain old block")
		}
	})
	t.Run("新增 block", func(t *testing.T) {
		got := upsertNamedTomlBlock("existing = \"val\"\n", "model_providers.ccx", "[model_providers.ccx]\nname = \"CCX\"\n")
		if !strings.Contains(got, `[model_providers.ccx]`) {
			t.Errorf("expected new block, got %q", got)
		}
	})
}

func TestRestoreTopLevelTomlString(t *testing.T) {
	t.Run("恢复原值", func(t *testing.T) {
		orig := "original"
		got := restoreTopLevelTomlString(`model_provider = "ccx"`, "model_provider", &orig)
		if !strings.Contains(got, `"original"`) {
			t.Errorf("expected original, got %q", got)
		}
	})
	t.Run("nil 删除行", func(t *testing.T) {
		got := restoreTopLevelTomlString("model_provider = \"ccx\"\nother = \"val\"\n", "model_provider", nil)
		if strings.Contains(got, "model_provider") {
			t.Errorf("should have removed line, got %q", got)
		}
		if !strings.Contains(got, "other") {
			t.Errorf("should keep other lines")
		}
	})
}

func TestRestoreNamedTomlBlock(t *testing.T) {
	t.Run("nil 删除 block", func(t *testing.T) {
		content := "[model_providers.ccx]\nname = \"CCX\"\n\n[other]\nk = \"v\"\n"
		got := restoreNamedTomlBlock(content, "model_providers.ccx", nil)
		if strings.Contains(got, "model_providers.ccx") {
			t.Errorf("should have removed block, got %q", got)
		}
		if !strings.Contains(got, "[other]") {
			t.Errorf("should keep other block")
		}
	})
	t.Run("恢复原 block", func(t *testing.T) {
		orig := "[model_providers.ccx]\nname = \"Original\"\n"
		content := "[model_providers.ccx]\nname = \"CCX\"\n"
		got := restoreNamedTomlBlock(content, "model_providers.ccx", &orig)
		if !strings.Contains(got, `"Original"`) {
			t.Errorf("expected original, got %q", got)
		}
	})
}

func TestDetectClaudeProvider(t *testing.T) {
	cases := []struct {
		baseURL string
		want    string
	}{
		{"", ""},
		{"http://127.0.0.1:3688", ProviderCCX},
		{"http://localhost:3688", ProviderCCX},
		{"https://api.deepseek.com/anthropic", ProviderDeepSeek},
		{"https://api.mimo.xiaomi.com/v1", ProviderMiMo},
		{"https://xiaomimimo.com/v1", ProviderMiMo},
		{"https://cp.compshare.cn", ProviderCompshare},
		{"https://custom-api.example.com/v1", ProviderCustom},
	}
	for _, c := range cases {
		t.Run(c.want, func(t *testing.T) {
			got := detectClaudeProvider(c.baseURL)
			if got != c.want {
				t.Errorf("detectClaudeProvider(%q) = %q, want %q", c.baseURL, got, c.want)
			}
		})
	}
}

func TestNormalizeClaudeProvider(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"", ProviderCCX},
		{"ccx", ProviderCCX},
		{"CCX", ProviderCCX},
		{"deepseek", ProviderDeepSeek},
		{"DeepSeek", ProviderDeepSeek},
		{"mimo", ProviderMiMo},
		{"MIMO", ProviderMiMo},
		{"compshare", ProviderCompshare},
		{"Compshare", ProviderCompshare},
		{"custom-provider", "custom-provider"},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			if got := normalizeClaudeProvider(c.input); got != c.want {
				t.Errorf("normalizeClaudeProvider(%q) = %q, want %q", c.input, got, c.want)
			}
		})
	}
}

func TestIsLocalBaseURL(t *testing.T) {
	cases := []struct {
		value string
		want  bool
	}{
		{"http://127.0.0.1:3688", true},
		{"http://localhost:3688", true},
		{"https://api.deepseek.com", false},
		{"", false},
	}
	for _, c := range cases {
		if got := isLocalBaseURL(c.value); got != c.want {
			t.Errorf("isLocalBaseURL(%q) = %v, want %v", c.value, got, c.want)
		}
	}
}

func TestClaudeBaseURL(t *testing.T) {
	if got := claudeBaseURL(3688); got != "http://127.0.0.1:3688" {
		t.Errorf("got %q", got)
	}
}

func TestCodexBaseURL(t *testing.T) {
	if got := codexBaseURL(3688); got != "http://127.0.0.1:3688/v1" {
		t.Errorf("got %q", got)
	}
}

func TestCodexProviderBlock(t *testing.T) {
	block := codexProviderBlock("http://127.0.0.1:3688/v1")
	if !strings.Contains(block, `[model_providers.ccx]`) {
		t.Errorf("missing table header")
	}
	if !strings.Contains(block, `base_url = "http://127.0.0.1:3688/v1"`) {
		t.Errorf("missing base_url")
	}
}

func TestAppendUnique(t *testing.T) {
	got := appendUnique([]string{"a", "b"}, "c")
	if len(got) != 3 {
		t.Fatalf("expected 3, got %d", len(got))
	}
	got = appendUnique(got, "b")
	if len(got) != 3 {
		t.Errorf("should not add duplicate, got %d", len(got))
	}
	got = appendUnique(got, "")
	if len(got) != 3 {
		t.Errorf("should not add empty, got %d", len(got))
	}
}

func TestAppendUniqueMany(t *testing.T) {
	got := appendUniqueMany([]string{"a"}, []string{"b", "a", "c"})
	if len(got) != 3 {
		t.Errorf("expected 3, got %d", len(got))
	}
}

func TestProviderFromStoreKey(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"claude:deepseek", "deepseek"},
		{"channel:mimo", "mimo"},
		{"standalone", "standalone"},
	}
	for _, c := range cases {
		if got := providerFromStoreKey(c.input); got != c.want {
			t.Errorf("providerFromStoreKey(%q) = %q, want %q", c.input, got, c.want)
		}
	}
}

func TestUsageFromStoreKey(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"claude:deepseek", "agent-direct"},
		{"codex:openai", "codex-direct"},
		{"channel:mimo", "channel"},
		{"standalone", "manual"},
	}
	for _, c := range cases {
		if got := usageFromStoreKey(c.input); got != c.want {
			t.Errorf("usageFromStoreKey(%q) = %q, want %q", c.input, got, c.want)
		}
	}
}

func TestGetNestedString(t *testing.T) {
	data := map[string]any{
		"env": map[string]any{
			"ANTHROPIC_BASE_URL": "http://localhost",
		},
	}
	if v, ok := getNestedString(data, "env", "ANTHROPIC_BASE_URL"); !ok || v != "http://localhost" {
		t.Errorf("got (%q, %v)", v, ok)
	}
	if _, ok := getNestedString(data, "env", "MISSING"); ok {
		t.Error("expected not found")
	}
	if _, ok := getNestedString(data, "missing"); ok {
		t.Error("expected not found for top-level")
	}
}

func TestOptionalString(t *testing.T) {
	if got := optionalString("val", true); got == nil || *got != "val" {
		t.Errorf("expected non-nil")
	}
	if got := optionalString("val", false); got != nil {
		t.Errorf("expected nil")
	}
}

func TestRestoreStringField(t *testing.T) {
	data := map[string]any{"key": "old"}
	restoreStringField(data, "key", strPtr("new"))
	if data["key"] != "new" {
		t.Errorf("expected new, got %v", data["key"])
	}
	restoreStringField(data, "key", nil)
	if _, ok := data["key"]; ok {
		t.Error("expected deleted")
	}
}

func strPtr(s string) *string { return &s }

// ── Part B: Service 集成测试（t.TempDir） ─────────────────

func newTestService(t *testing.T) *Service {
	t.Helper()
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	dataDir := filepath.Join(t.TempDir(), "ccx-data")
	svc, err := New(dataDir)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	return svc
}

func TestNew(t *testing.T) {
	svc := newTestService(t)
	if svc.homeDir == "" {
		t.Error("homeDir should not be empty")
	}
	if svc.stateDir == "" {
		t.Error("stateDir should not be empty")
	}
	if _, err := os.Stat(svc.stateDir); os.IsNotExist(err) {
		t.Error("stateDir should be created")
	}
}

func TestNewWithDefaultDataDir(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	if runtime.GOOS == "linux" {
		t.Setenv("XDG_STATE_HOME", "")
	}
	svc, err := New("")
	if err != nil {
		t.Fatalf("New('') failed: %v", err)
	}
	expected := filepath.Join(tmpHome, ".config", "ccx-desktop", "agent-config-state")
	if runtime.GOOS == "linux" {
		expected = filepath.Join(tmpHome, ".local", "state", "ccx", "agent-config-state")
	}
	if svc.stateDir != expected {
		t.Errorf("stateDir = %q, want %q", svc.stateDir, expected)
	}
}

func TestGetStatusClaude_NoConfig(t *testing.T) {
	svc := newTestService(t)
	status, err := svc.GetStatus(PlatformClaude, 3688)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if status.Configured {
		t.Error("should not be configured when no settings.json exists")
	}
	// 空 base_url 时 detectClaudeProvider 返回 ""
	if status.Provider != "" {
		t.Errorf("provider = %q, want empty", status.Provider)
	}
}

func TestGetStatusClaude_CCXProvider(t *testing.T) {
	svc := newTestService(t)
	settingsPath := filepath.Join(svc.homeDir, ".claude", "settings.json")
	os.MkdirAll(filepath.Dir(settingsPath), 0o755)
	data := map[string]any{
		"env": map[string]any{
			"ANTHROPIC_BASE_URL": "http://127.0.0.1:3688",
		},
	}
	writeJSON(settingsPath, data)

	status, err := svc.GetStatus(PlatformClaude, 3688)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if status.Provider != ProviderCCX {
		t.Errorf("provider = %q, want %q", status.Provider, ProviderCCX)
	}
	if !status.MatchesCurrentPort {
		t.Error("should match current port")
	}
	if status.NeedsUpdate {
		t.Error("should not need update when port matches")
	}
}

func TestGetStatusClaude_DeepSeekProvider(t *testing.T) {
	svc := newTestService(t)
	settingsPath := filepath.Join(svc.homeDir, ".claude", "settings.json")
	os.MkdirAll(filepath.Dir(settingsPath), 0o755)
	data := map[string]any{
		"env": map[string]any{
			"ANTHROPIC_BASE_URL": "https://api.deepseek.com/anthropic",
		},
	}
	writeJSON(settingsPath, data)

	status, err := svc.GetStatus(PlatformClaude, 3688)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if status.Provider != ProviderDeepSeek {
		t.Errorf("provider = %q, want %q", status.Provider, ProviderDeepSeek)
	}
	if !status.Configured {
		t.Error("should be configured for DeepSeek")
	}
}

func TestGetStatusClaude_CompshareProvider(t *testing.T) {
	svc := newTestService(t)
	settingsPath := filepath.Join(svc.homeDir, ".claude", "settings.json")
	os.MkdirAll(filepath.Dir(settingsPath), 0o755)
	data := map[string]any{
		"env": map[string]any{
			"ANTHROPIC_BASE_URL": "https://cp.compshare.cn",
		},
	}
	writeJSON(settingsPath, data)

	status, err := svc.GetStatus(PlatformClaude, 3688)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if status.Provider != ProviderCompshare {
		t.Errorf("provider = %q, want %q", status.Provider, ProviderCompshare)
	}
	if !status.Configured {
		t.Error("should be configured for Compshare")
	}
}

func TestApplyAndRestoreClaude(t *testing.T) {
	svc := newTestService(t)
	settingsPath := filepath.Join(svc.homeDir, ".claude", "settings.json")
	os.MkdirAll(filepath.Dir(settingsPath), 0o755)
	original := map[string]any{
		"env": map[string]any{
			"ANTHROPIC_BASE_URL": "https://original.example.com",
		},
	}
	writeJSON(settingsPath, original)

	err := svc.Apply(ApplyAgentConfigRequest{Platform: PlatformClaude}, 3688, "test-key")
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	// 验证注入后配置
	var after map[string]any
	readJSON(settingsPath, &after)
	env := after["env"].(map[string]any)
	if env["ANTHROPIC_BASE_URL"] != "http://127.0.0.1:3688" {
		t.Errorf("base_url = %v", env["ANTHROPIC_BASE_URL"])
	}
	if env["ANTHROPIC_AUTH_TOKEN"] != "test-key" {
		t.Errorf("auth_token = %v", env["ANTHROPIC_AUTH_TOKEN"])
	}

	// Restore
	err = svc.Restore(PlatformClaude)
	if err != nil {
		t.Fatalf("Restore failed: %v", err)
	}
	var restored map[string]any
	readJSON(settingsPath, &restored)
	env = restored["env"].(map[string]any)
	if env["ANTHROPIC_BASE_URL"] != "https://original.example.com" {
		t.Errorf("restored base_url = %v", env["ANTHROPIC_BASE_URL"])
	}
}

func TestGetStatusCodex_NoConfig(t *testing.T) {
	svc := newTestService(t)
	status, err := svc.GetStatus(PlatformCodex, 3688)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if status.Configured {
		t.Error("should not be configured when no config.toml exists")
	}
}

func TestApplyAndRestoreCodex(t *testing.T) {
	svc := newTestService(t)
	configPath := filepath.Join(svc.homeDir, ".codex", "config.toml")
	authPath := filepath.Join(svc.homeDir, ".codex", "auth.json")
	os.MkdirAll(filepath.Dir(configPath), 0o755)

	err := svc.Apply(ApplyAgentConfigRequest{Platform: PlatformCodex}, 3688, "test-key")
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	// 验证 config.toml — 新格式 openai_base_url 模式
	content, _ := os.ReadFile(configPath)
	if !strings.Contains(string(content), `model_provider = "openai"`) {
		t.Errorf("config.toml should contain model_provider = openai")
	}
	if !strings.Contains(string(content), `openai_base_url = "http://127.0.0.1:3688/v1"`) {
		t.Errorf("config.toml should contain openai_base_url")
	}
	if strings.Contains(string(content), `[model_providers.ccx]`) {
		t.Errorf("config.toml should NOT contain legacy [model_providers.ccx] block")
	}

	// 验证 auth.json
	var authData map[string]any
	readJSON(authPath, &authData)
	if authData["OPENAI_API_KEY"] != "test-key" {
		t.Errorf("auth OPENAI_API_KEY = %v", authData["OPENAI_API_KEY"])
	}

	// Restore
	err = svc.Restore(PlatformCodex)
	if err != nil {
		t.Fatalf("Restore failed: %v", err)
	}
}

func TestSaveAndLoadProviderKeys(t *testing.T) {
	svc := newTestService(t)

	err := svc.SaveProviderKeyAsset(ProviderKeyAsset{
		Provider: ProviderDeepSeek,
		APIKey:   "sk-test-key",
		BaseURL:  "https://api.deepseek.com",
	})
	if err != nil {
		t.Fatalf("SaveProviderKeyAsset failed: %v", err)
	}
	err = svc.SaveProviderKeyAsset(ProviderKeyAsset{
		Provider: ProviderCompshare,
		APIKey:   "cs-test-key",
		BaseURL:  "https://cp.compshare.cn",
	})
	if err != nil {
		t.Fatalf("SaveProviderKeyAsset failed: %v", err)
	}

	keys := svc.GetSavedProviderKeys()
	if keys["channel:"+ProviderDeepSeek] != "sk-test-key" {
		t.Errorf("channel key = %q", keys["channel:"+ProviderDeepSeek])
	}
	if keys[PlatformClaude+":"+ProviderDeepSeek] != "sk-test-key" {
		t.Errorf("claude key = %q", keys[PlatformClaude+":"+ProviderDeepSeek])
	}
	if keys["channel:"+ProviderCompshare] != "cs-test-key" {
		t.Errorf("compshare channel key = %q", keys["channel:"+ProviderCompshare])
	}
	if keys[PlatformClaude+":"+ProviderCompshare] != "cs-test-key" {
		t.Errorf("compshare claude key = %q", keys[PlatformClaude+":"+ProviderCompshare])
	}

	assets := svc.GetProviderKeyAssets()
	foundDeepSeek := false
	foundCompshare := false
	for _, a := range assets {
		switch a.Provider {
		case ProviderDeepSeek:
			foundDeepSeek = true
			if a.APIKey != "sk-test-key" {
				t.Errorf("asset APIKey = %q", a.APIKey)
			}
		case ProviderCompshare:
			foundCompshare = true
			if a.APIKey != "cs-test-key" {
				t.Errorf("compshare asset APIKey = %q", a.APIKey)
			}
		}
	}
	if !foundDeepSeek {
		t.Error("DeepSeek asset not found")
	}
	if !foundCompshare {
		t.Error("Compshare asset not found")
	}
}

// ── P1 回归: Codex 第三方 provider 状态识别 ──

func TestGetStatusCodex_ThirdPartyProvider(t *testing.T) {
	svc := newTestService(t)
	configPath := filepath.Join(svc.homeDir, ".codex", "config.toml")
	authPath := filepath.Join(svc.homeDir, ".codex", "auth.json")
	os.MkdirAll(filepath.Dir(configPath), 0o755)

	// 模拟已配置 dashscope 的 config.toml
	tomlContent := `model_provider = "dashscope"

[model_providers.dashscope]
base_url = "https://dashscope.aliyuncs.com/compatible-mode/v1"
`
	os.WriteFile(configPath, []byte(tomlContent), 0o644)
	writeJSON(authPath, map[string]any{"OPENAI_API_KEY": "sk-ds-key"})

	status, err := svc.GetStatus(PlatformCodex, 3688)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if status.Provider != ProviderDashScope {
		t.Errorf("Provider = %q, want %q", status.Provider, ProviderDashScope)
	}
	if !status.Configured {
		t.Error("Configured should be true for third-party provider")
	}
}

func TestGetStatusCodex_OpenCodeZenProvider(t *testing.T) {
	svc := newTestService(t)
	configPath := filepath.Join(svc.homeDir, ".codex", "config.toml")
	authPath := filepath.Join(svc.homeDir, ".codex", "auth.json")
	os.MkdirAll(filepath.Dir(configPath), 0o755)

	tomlContent := `model_provider = "opencode-zen"

[model_providers.opencode-zen]
base_url = "https://api.opencode.ai/v1"
`
	os.WriteFile(configPath, []byte(tomlContent), 0o644)
	writeJSON(authPath, map[string]any{"OPENAI_API_KEY": "sk-zen-key"})

	status, err := svc.GetStatus(PlatformCodex, 3688)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if status.Provider != ProviderOpenCodeZen {
		t.Errorf("Provider = %q, want %q", status.Provider, ProviderOpenCodeZen)
	}
	if !status.Configured {
		t.Error("Configured should be true for opencode-zen provider")
	}
}

// ── P1 回归: Codex 恢复清理第三方 provider block ──

func TestApplyAndRestoreCodex_ThirdPartyCleanup(t *testing.T) {
	svc := newTestService(t)
	configPath := filepath.Join(svc.homeDir, ".codex", "config.toml")
	authPath := filepath.Join(svc.homeDir, ".codex", "auth.json")
	os.MkdirAll(filepath.Dir(configPath), 0o755)

	// 先 Apply 第三方 provider
	err := svc.Apply(ApplyAgentConfigRequest{
		Platform: PlatformCodex,
		Provider: ProviderDashScope,
		APIKey:   "sk-ds-key",
		BaseURL:  "https://dashscope.aliyuncs.com/compatible-mode/v1",
	}, 0, "")
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	content, _ := os.ReadFile(configPath)
	s := string(content)
	if !strings.Contains(s, `model_provider = "dashscope"`) {
		t.Error("config.toml should contain model_provider = dashscope after apply")
	}
	if !strings.Contains(s, `[model_providers.dashscope]`) {
		t.Error("config.toml should contain dashscope provider block after apply")
	}
	authData, _, _ := readJSONMap(authPath)
	if authData["OPENAI_API_KEY"] != "sk-ds-key" {
		t.Errorf("OPENAI_API_KEY = %v, want sk-ds-key", authData["OPENAI_API_KEY"])
	}
	if authData["auth_mode"] != "chatgpt" {
		t.Errorf("auth_mode = %v, want chatgpt", authData["auth_mode"])
	}

	// Restore
	err = svc.Restore(PlatformCodex)
	if err != nil {
		t.Fatalf("Restore failed: %v", err)
	}

	content, _ = os.ReadFile(configPath)
	s = string(content)
	if strings.Contains(s, `[model_providers.dashscope]`) {
		t.Error("config.toml should NOT contain dashscope provider block after restore")
	}
}

// ── Codex openai_base_url 新格式测试 ──

func TestGetStatusCodex_NewStyleCCXProxy(t *testing.T) {
	svc := newTestService(t)
	configPath := filepath.Join(svc.homeDir, ".codex", "config.toml")
	authPath := filepath.Join(svc.homeDir, ".codex", "auth.json")
	os.MkdirAll(filepath.Dir(configPath), 0o755)

	tomlContent := `model_provider = "openai"
openai_base_url = "http://127.0.0.1:3688/v1"
`
	os.WriteFile(configPath, []byte(tomlContent), 0o644)
	writeJSON(authPath, map[string]any{"OPENAI_API_KEY": "test-key"})

	status, err := svc.GetStatus(PlatformCodex, 3688)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if status.Provider != ProviderCCX {
		t.Errorf("Provider = %q, want %q", status.Provider, ProviderCCX)
	}
	if !status.Configured {
		t.Error("Configured should be true for new-style CCX proxy")
	}
	if !status.MatchesCurrentPort {
		t.Error("MatchesCurrentPort should be true when port matches")
	}
}

func TestGetStatusCodex_NewStyleCCXProxy_WrongPort(t *testing.T) {
	svc := newTestService(t)
	configPath := filepath.Join(svc.homeDir, ".codex", "config.toml")
	authPath := filepath.Join(svc.homeDir, ".codex", "auth.json")
	os.MkdirAll(filepath.Dir(configPath), 0o755)

	tomlContent := `model_provider = "openai"
openai_base_url = "http://127.0.0.1:9999/v1"
`
	os.WriteFile(configPath, []byte(tomlContent), 0o644)
	writeJSON(authPath, map[string]any{"OPENAI_API_KEY": "test-key"})

	status, err := svc.GetStatus(PlatformCodex, 3688)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if status.Provider != ProviderCCX {
		t.Errorf("Provider = %q, want %q", status.Provider, ProviderCCX)
	}
	if !status.NeedsUpdate {
		t.Error("NeedsUpdate should be true when port mismatches")
	}
	if status.MatchesCurrentPort {
		t.Error("MatchesCurrentPort should be false when port mismatches")
	}
}

func TestGetStatusCodex_OpenAIDirect(t *testing.T) {
	svc := newTestService(t)
	configPath := filepath.Join(svc.homeDir, ".codex", "config.toml")
	authPath := filepath.Join(svc.homeDir, ".codex", "auth.json")
	os.MkdirAll(filepath.Dir(configPath), 0o755)

	tomlContent := `model_provider = "openai"
`
	os.WriteFile(configPath, []byte(tomlContent), 0o644)
	writeJSON(authPath, map[string]any{"OPENAI_API_KEY": "sk-openai-key"})

	status, err := svc.GetStatus(PlatformCodex, 3688)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if status.Provider != ProviderOpenAI {
		t.Errorf("Provider = %q, want %q", status.Provider, ProviderOpenAI)
	}
	if !status.Configured {
		t.Error("Configured should be true for OpenAI direct")
	}
}

func TestGetStatusCodex_LegacyCCX(t *testing.T) {
	svc := newTestService(t)
	configPath := filepath.Join(svc.homeDir, ".codex", "config.toml")
	authPath := filepath.Join(svc.homeDir, ".codex", "auth.json")
	os.MkdirAll(filepath.Dir(configPath), 0o755)

	tomlContent := `model_provider = "ccx"

[model_providers.ccx]
name = "CCX Proxy"
base_url = "http://127.0.0.1:3688/v1"
wire_api = "responses"
`
	os.WriteFile(configPath, []byte(tomlContent), 0o644)
	writeJSON(authPath, map[string]any{"OPENAI_API_KEY": "test-key"})

	status, err := svc.GetStatus(PlatformCodex, 3688)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if status.Provider != ProviderCCX {
		t.Errorf("Provider = %q, want %q", status.Provider, ProviderCCX)
	}
	if !status.MatchesCurrentPort {
		t.Error("MatchesCurrentPort should be true for legacy CCX with matching port")
	}
}

func TestApplyCodexOpenAI_RemovesBaseURL(t *testing.T) {
	svc := newTestService(t)
	configPath := filepath.Join(svc.homeDir, ".codex", "config.toml")
	authPath := filepath.Join(svc.homeDir, ".codex", "auth.json")
	os.MkdirAll(filepath.Dir(configPath), 0o755)

	// 先写入一个含 openai_base_url 的 config
	tomlContent := `model_provider = "openai"
openai_base_url = "http://127.0.0.1:3688/v1"
`
	os.WriteFile(configPath, []byte(tomlContent), 0o644)
	writeJSON(authPath, map[string]any{"OPENAI_API_KEY": "sk-openai-key"})

	// 应用 OpenAI direct（apiKey 须通过 req.APIKey 传入）
	err := svc.Apply(ApplyAgentConfigRequest{Platform: PlatformCodex, Provider: ProviderOpenAI, APIKey: "sk-openai-key"}, 0, "")
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	content, _ := os.ReadFile(configPath)
	s := string(content)
	if strings.Contains(s, "openai_base_url") {
		t.Error("config.toml should NOT contain openai_base_url after switching to OpenAI direct")
	}
	if !strings.Contains(s, `model_provider = "openai"`) {
		t.Error("config.toml should contain model_provider = openai")
	}
}

func TestApplyCodexOpenAI_OAuthMode(t *testing.T) {
	svc := newTestService(t)
	configPath := filepath.Join(svc.homeDir, ".codex", "config.toml")
	authPath := filepath.Join(svc.homeDir, ".codex", "auth.json")
	os.MkdirAll(filepath.Dir(configPath), 0o755)

	// 先写入一个含 openai_base_url 的 config（CCX 代理残留）
	tomlContent := `model_provider = "openai"
openai_base_url = "http://127.0.0.1:3688/v1"
`
	os.WriteFile(configPath, []byte(tomlContent), 0o644)
	writeJSON(authPath, map[string]any{"OPENAI_API_KEY": "ccx-proxy-key", "auth_mode": "chatgpt"})

	// 应用 OpenAI direct，不提供 apiKey → 应走 OAuth 模式
	err := svc.Apply(ApplyAgentConfigRequest{Platform: PlatformCodex, Provider: ProviderOpenAI}, 0, "")
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	// 验证 config.toml 清理了 proxy 残留
	configContent, _ := os.ReadFile(configPath)
	s := string(configContent)
	if strings.Contains(s, "openai_base_url") {
		t.Error("config.toml should NOT contain openai_base_url after switching to OpenAI direct")
	}
	if !strings.Contains(s, `model_provider = "openai"`) {
		t.Error("config.toml should contain model_provider = openai")
	}

	// 验证 auth.json：auth_mode = "chatgpt"，OPENAI_API_KEY = nil
	authContent, _ := os.ReadFile(authPath)
	var authData map[string]any
	json.Unmarshal(authContent, &authData)
	if authData["auth_mode"] != "chatgpt" {
		t.Errorf("auth.json should have auth_mode = chatgpt, got %v", authData["auth_mode"])
	}
	if authData["OPENAI_API_KEY"] != nil {
		t.Errorf("auth.json should have OPENAI_API_KEY = null, got %v", authData["OPENAI_API_KEY"])
	}
}

func TestApplyCodexOpenAI_WithApiKey(t *testing.T) {
	svc := newTestService(t)
	configPath := filepath.Join(svc.homeDir, ".codex", "config.toml")
	authPath := filepath.Join(svc.homeDir, ".codex", "auth.json")
	os.MkdirAll(filepath.Dir(configPath), 0o755)

	tomlContent := `model_provider = "ccx"
`
	os.WriteFile(configPath, []byte(tomlContent), 0o644)
	writeJSON(authPath, map[string]any{"OPENAI_API_KEY": nil, "auth_mode": "chatgpt"})

	// 应用 OpenAI direct，提供了 apiKey → 应写入 key 并设置 auth_mode = "api"
	err := svc.Apply(ApplyAgentConfigRequest{Platform: PlatformCodex, Provider: ProviderOpenAI, APIKey: "sk-my-openai-key"}, 0, "")
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	// 验证 auth.json：OPENAI_API_KEY 写入，auth_mode = "api"
	authContent, _ := os.ReadFile(authPath)
	var authData map[string]any
	json.Unmarshal(authContent, &authData)
	if authData["OPENAI_API_KEY"] != "sk-my-openai-key" {
		t.Errorf("auth.json should have OPENAI_API_KEY = sk-my-openai-key, got %v", authData["OPENAI_API_KEY"])
	}
	if authData["auth_mode"] != "api" {
		t.Errorf("auth.json should have auth_mode = api, got %v", authData["auth_mode"])
	}
}

func TestApplyCodex_MigratesFromLegacyCCX(t *testing.T) {
	svc := newTestService(t)
	configPath := filepath.Join(svc.homeDir, ".codex", "config.toml")
	authPath := filepath.Join(svc.homeDir, ".codex", "auth.json")
	os.MkdirAll(filepath.Dir(configPath), 0o755)

	// 写入旧格式 ccx 配置
	tomlContent := `model_provider = "ccx"

[model_providers.ccx]
name = "CCX Proxy"
base_url = "http://127.0.0.1:9999/v1"
wire_api = "responses"
`
	os.WriteFile(configPath, []byte(tomlContent), 0o644)
	writeJSON(authPath, map[string]any{"OPENAI_API_KEY": "old-key"})

	// 应用新格式
	err := svc.Apply(ApplyAgentConfigRequest{Platform: PlatformCodex}, 3688, "new-key")
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	content, _ := os.ReadFile(configPath)
	s := string(content)
	if !strings.Contains(s, `model_provider = "openai"`) {
		t.Error("config.toml should contain model_provider = openai after migration")
	}
	if !strings.Contains(s, `openai_base_url = "http://127.0.0.1:3688/v1"`) {
		t.Error("config.toml should contain openai_base_url with correct port")
	}
	if strings.Contains(s, `[model_providers.ccx]`) {
		t.Error("config.toml should NOT contain legacy [model_providers.ccx] block after migration")
	}
}

// ── Codex 模式切换测试 ──

func TestGetStatusCodex_QuickMode(t *testing.T) {
	svc := newTestService(t)
	configPath := filepath.Join(svc.homeDir, ".codex", "config.toml")
	authPath := filepath.Join(svc.homeDir, ".codex", "auth.json")
	os.MkdirAll(filepath.Dir(configPath), 0o755)

	tomlContent := `model_provider = "openai"
openai_base_url = "http://127.0.0.1:3688/v1"
`
	os.WriteFile(configPath, []byte(tomlContent), 0o644)
	writeJSON(authPath, map[string]any{"OPENAI_API_KEY": "test-key"})

	status, err := svc.GetStatus(PlatformCodex, 3688)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if status.Provider != ProviderCCX {
		t.Errorf("Provider = %q, want %q", status.Provider, ProviderCCX)
	}
	if status.Mode != "" {
		t.Errorf("Mode = %q, want empty for quick openai_base_url mode", status.Mode)
	}
	if !status.MatchesCurrentPort {
		t.Error("MatchesCurrentPort should be true")
	}
}

func TestGetStatusCodex_PluginMode(t *testing.T) {
	svc := newTestService(t)
	configPath := filepath.Join(svc.homeDir, ".codex", "config.toml")
	authPath := filepath.Join(svc.homeDir, ".codex", "auth.json")
	os.MkdirAll(filepath.Dir(configPath), 0o755)

	tomlContent := `model_provider = "ccx"

[model_providers.ccx]
name = "CCX Proxy"
base_url = "http://127.0.0.1:3688/v1"
wire_api = "responses"
requires_openai_auth = true
experimental_bearer_token = "test-key"
`
	os.WriteFile(configPath, []byte(tomlContent), 0o644)
	writeJSON(authPath, map[string]any{"OPENAI_API_KEY": "test-key", "auth_mode": "chatgpt"})

	status, err := svc.GetStatus(PlatformCodex, 3688)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if status.Provider != ProviderCCX {
		t.Errorf("Provider = %q, want %q", status.Provider, ProviderCCX)
	}
	if status.Mode != "plugin" {
		t.Errorf("Mode = %q, want plugin", status.Mode)
	}
	if !status.MatchesCurrentPort {
		t.Error("MatchesCurrentPort should be true")
	}
}

func TestApplyCodex_PluginMode(t *testing.T) {
	svc := newTestService(t)
	configPath := filepath.Join(svc.homeDir, ".codex", "config.toml")
	authPath := filepath.Join(svc.homeDir, ".codex", "auth.json")
	os.MkdirAll(filepath.Dir(configPath), 0o755)

	err := svc.Apply(ApplyAgentConfigRequest{Platform: PlatformCodex, Provider: ProviderCCX, Mode: "plugin"}, 3688, "test-key")
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	content, _ := os.ReadFile(configPath)
	s := string(content)
	if !strings.Contains(s, `model_provider = "ccx"`) {
		t.Error("config.toml should contain model_provider = ccx")
	}
	if !strings.Contains(s, `[model_providers.ccx]`) {
		t.Error("config.toml should contain [model_providers.ccx]")
	}
	if !strings.Contains(s, `requires_openai_auth = true`) {
		t.Error("config.toml should contain requires_openai_auth = true")
	}
	if !strings.Contains(s, `experimental_bearer_token = "test-key"`) {
		t.Error("config.toml should contain experimental_bearer_token")
	}
	blockIndex := strings.Index(s, `[model_providers.ccx]`)
	tokenIndex := strings.Index(s, `experimental_bearer_token = "test-key"`)
	if tokenIndex < blockIndex {
		t.Error("experimental_bearer_token should be inside [model_providers.ccx] block")
	}
	if strings.Contains(s, "openai_base_url") {
		t.Error("config.toml should not contain openai_base_url in plugin mode")
	}

	authData, _, _ := readJSONMap(authPath)
	if authData["auth_mode"] != "chatgpt" {
		t.Errorf("auth_mode = %v, want chatgpt", authData["auth_mode"])
	}
	if authData["OPENAI_API_KEY"] != "test-key" {
		t.Errorf("OPENAI_API_KEY = %v, want test-key", authData["OPENAI_API_KEY"])
	}
}

func TestApplyCodex_SwitchFromPluginToQuick(t *testing.T) {
	svc := newTestService(t)
	configPath := filepath.Join(svc.homeDir, ".codex", "config.toml")
	authPath := filepath.Join(svc.homeDir, ".codex", "auth.json")
	os.MkdirAll(filepath.Dir(configPath), 0o755)

	tomlContent := `model_provider = "ccx"

[model_providers.ccx]
name = "CCX Proxy"
base_url = "http://127.0.0.1:3688/v1"
wire_api = "responses"
requires_openai_auth = true
experimental_bearer_token = "old-key"
`
	os.WriteFile(configPath, []byte(tomlContent), 0o644)
	writeJSON(authPath, map[string]any{"OPENAI_API_KEY": "old-key", "auth_mode": "chatgpt"})

	err := svc.Apply(ApplyAgentConfigRequest{Platform: PlatformCodex, Provider: ProviderCCX, Mode: "quick"}, 3688, "new-key")
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	content, _ := os.ReadFile(configPath)
	s := string(content)
	if !strings.Contains(s, `model_provider = "openai"`) {
		t.Error("config.toml should contain model_provider = openai")
	}
	if !strings.Contains(s, `openai_base_url = "http://127.0.0.1:3688/v1"`) {
		t.Error("config.toml should contain openai_base_url")
	}
	if strings.Contains(s, "experimental_bearer_token") {
		t.Error("config.toml should not contain experimental_bearer_token after switching to quick")
	}
	if strings.Contains(s, `[model_providers.ccx]`) {
		t.Error("config.toml should not contain ccx provider block after switching to quick")
	}

	authData, _, _ := readJSONMap(authPath)
	if authData["OPENAI_API_KEY"] != "new-key" {
		t.Errorf("OPENAI_API_KEY = %v, want new-key", authData["OPENAI_API_KEY"])
	}
	if authData["auth_mode"] != "chatgpt" {
		t.Errorf("auth_mode = %v, want chatgpt", authData["auth_mode"])
	}
}

// ── helpers ──

func writeJSON(path string, data any) {
	b, _ := json.MarshalIndent(data, "", "  ")
	os.MkdirAll(filepath.Dir(path), 0o755)
	os.WriteFile(path, append(b, '\n'), 0o644)
}

func readJSON(path string, dest any) {
	b, _ := os.ReadFile(path)
	json.Unmarshal(b, dest)
}
