package keypool

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"

	"github.com/BenedictKing/ccx/internal/config"
	"github.com/BenedictKing/ccx/internal/ratelimit"
)

type Candidate struct {
	APIKey     string
	Config     config.APIKeyConfig
	Index      int
	Scope      string
	QuotaGroup string
}

type Selection struct {
	APIKey         string
	CredentialID   string
	CredentialName string
	QuotaGroup     string
	LimiterScope   string
	Config         config.APIKeyConfig
}

func HasEffectiveConfig(upstream *config.UpstreamConfig) bool {
	if upstream == nil {
		return false
	}
	for _, cfg := range upstream.APIKeyConfigs {
		if strings.TrimSpace(cfg.Name) != "" || cfg.Enabled != nil || strings.TrimSpace(cfg.QuotaGroup) != "" ||
			cfg.RateLimitRPM > 0 || cfg.RateLimitWindowMinutes > 0 || cfg.RateLimitMaxConcurrent > 0 ||
			cfg.RateLimitAutoFromHeaders != nil || cfg.Weight > 0 || len(cfg.Models) > 0 {
			return true
		}
	}
	return false
}

func Candidates(upstream *config.UpstreamConfig, failedKeys map[string]bool) []Candidate {
	return CandidatesForModel(upstream, failedKeys, "")
}

// CandidatesForModel 返回可用 key 列表，过滤 enabled=false、failedKeys 和模型白名单。
// model 为空时不按模型过滤。
func CandidatesForModel(upstream *config.UpstreamConfig, failedKeys map[string]bool, model string) []Candidate {
	if upstream == nil || len(upstream.APIKeys) == 0 {
		return nil
	}

	configs := config.NormalizeAPIKeyConfigsForView(*upstream)
	byKey := make(map[string]config.APIKeyConfig, len(configs))
	for _, cfg := range configs {
		byKey[cfg.Key] = cfg
	}

	model = strings.TrimSpace(model)
	out := make([]Candidate, 0, len(upstream.APIKeys))
	for i, key := range upstream.APIKeys {
		key = strings.TrimSpace(key)
		if key == "" || failedKeys[key] {
			continue
		}
		cfg := byKey[key]
		if cfg.Key == "" {
			cfg.Key = key
		}
		if cfg.Enabled != nil && !*cfg.Enabled {
			continue
		}
		if model != "" && len(cfg.Models) > 0 && !matchesModel(model, cfg.Models) {
			continue
		}
		quotaGroup := strings.TrimSpace(cfg.QuotaGroup)
		scope := "key:" + stableKeyID(key)
		if quotaGroup != "" {
			scope = "quota:" + stableKeyID("quota:"+quotaGroup)
		}
		out = append(out, Candidate{
			APIKey:     key,
			Config:     cfg,
			Index:      i,
			Scope:      scope,
			QuotaGroup: quotaGroup,
		})
	}

	// 按 weight 降序排序，weight 相同时保持原有顺序（稳定排序）
	if len(out) > 1 {
		sort.SliceStable(out, func(i, j int) bool {
			wi, wj := out[i].Config.Weight, out[j].Config.Weight
			if wi == 0 {
				wi = 1
			}
			if wj == 0 {
				wj = 1
			}
			return wi > wj
		})
	}

	return out
}

// matchesModel 检查 model 是否在允许列表中（支持通配符 *）。
func matchesModel(model string, models []string) bool {
	model = strings.ToLower(strings.TrimSpace(model))
	for _, pattern := range models {
		pattern = strings.ToLower(strings.TrimSpace(pattern))
		if pattern == "" {
			continue
		}
		if pattern == model {
			return true
		}
		if strings.HasPrefix(pattern, "*") && strings.HasSuffix(pattern, "*") {
			if strings.Contains(model, pattern[1:len(pattern)-1]) {
				return true
			}
		} else if strings.HasPrefix(pattern, "*") {
			if strings.HasSuffix(model, pattern[1:]) {
				return true
			}
		} else if strings.HasSuffix(pattern, "*") {
			if strings.HasPrefix(model, pattern[:len(pattern)-1]) {
				return true
			}
		}
	}
	return false
}

func ConfigForCandidate(channel config.UpstreamConfig, cfg config.APIKeyConfig) ratelimit.Config {
	rpm := cfg.RateLimitRPM
	if rpm <= 0 {
		rpm = channel.RateLimitRPM
	}
	windowMinutes := cfg.RateLimitWindowMinutes
	if windowMinutes <= 0 {
		windowMinutes = channel.RateLimitWindowMinutes
	}
	maxConcurrent := cfg.RateLimitMaxConcurrent
	if maxConcurrent <= 0 {
		maxConcurrent = channel.RateLimitMaxConcurrent
	}
	autoFromHeaders := channel.IsRateLimitAutoFromHeadersEnabled()
	if cfg.RateLimitAutoFromHeaders != nil {
		autoFromHeaders = *cfg.RateLimitAutoFromHeaders
	}
	return ratelimit.Config{
		RPM:             rpm,
		WindowSeconds:   windowMinutes * 60,
		MaxConcurrent:   maxConcurrent,
		AutoFromHeaders: autoFromHeaders,
	}
}

func stableKeyID(key string) string {
	sum := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sum[:])[:16]
}
