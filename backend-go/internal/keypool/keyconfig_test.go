package keypool

import (
	"testing"

	"github.com/BenedictKing/ccx/internal/config"
)

func ptrBool(v bool) *bool { return &v }

func TestCandidatesForModel_FiltersByModels(t *testing.T) {
	up := &config.UpstreamConfig{
		APIKeys: []string{"k1", "k2", "k3"},
		APIKeyConfigs: []config.APIKeyConfig{
			{Key: "k1", Models: []string{"claude-sonnet-4-5"}},
			{Key: "k2", Models: []string{"gpt-4*"}},
			{Key: "k3"}, // 无 Models，应匹配所有
		},
	}

	cands := CandidatesForModel(up, nil, "claude-sonnet-4-5")
	if len(cands) != 2 {
		t.Fatalf("want 2 candidates for claude-sonnet-4-5, got %d", len(cands))
	}
	keys := map[string]bool{}
	for _, c := range cands {
		keys[c.APIKey] = true
	}
	if !keys["k1"] || !keys["k3"] {
		t.Fatalf("expected k1 and k3, got %v", keys)
	}
}

func TestCandidatesForModel_WildcardPattern(t *testing.T) {
	up := &config.UpstreamConfig{
		APIKeys: []string{"k1", "k2"},
		APIKeyConfigs: []config.APIKeyConfig{
			{Key: "k1", Models: []string{"gpt-4*"}},
			{Key: "k2", Models: []string{"!gpt-*"}},
		},
	}

	cands := CandidatesForModel(up, nil, "gpt-4o")
	if len(cands) != 1 || cands[0].APIKey != "k1" {
		t.Fatalf("want k1 for gpt-4o, got %v", cands)
	}
}

func TestCandidatesForModel_WeightOrdering(t *testing.T) {
	up := &config.UpstreamConfig{
		APIKeys: []string{"k1", "k2", "k3"},
		APIKeyConfigs: []config.APIKeyConfig{
			{Key: "k1", Weight: 1},
			{Key: "k2", Weight: 5},
			{Key: "k3"}, // 默认 weight=0 => 1
		},
	}

	cands := CandidatesForModel(up, nil, "")
	if len(cands) != 3 {
		t.Fatalf("want 3, got %d", len(cands))
	}
	if cands[0].APIKey != "k2" {
		t.Fatalf("first candidate should be k2 (weight=5), got %s", cands[0].APIKey)
	}
}

func TestCandidatesForModel_EnabledFalseFiltered(t *testing.T) {
	up := &config.UpstreamConfig{
		APIKeys: []string{"k1", "k2"},
		APIKeyConfigs: []config.APIKeyConfig{
			{Key: "k1", Enabled: ptrBool(false)},
			{Key: "k2", Enabled: ptrBool(true)},
		},
	}

	cands := CandidatesForModel(up, nil, "")
	if len(cands) != 1 || cands[0].APIKey != "k2" {
		t.Fatalf("want only k2, got %v", cands)
	}
}

func TestCandidatesForModel_FailedKeysFiltered(t *testing.T) {
	up := &config.UpstreamConfig{
		APIKeys: []string{"k1", "k2"},
		APIKeyConfigs: []config.APIKeyConfig{
			{Key: "k1", Name: "a"},
			{Key: "k2", Name: "b"},
		},
	}

	cands := CandidatesForModel(up, map[string]bool{"k1": true}, "")
	if len(cands) != 1 || cands[0].APIKey != "k2" {
		t.Fatalf("want only k2, got %v", cands)
	}
}

func TestMatchesModel(t *testing.T) {
	tests := []struct {
		model   string
		models  []string
		want    bool
	}{
		{"claude-sonnet-4-5", []string{"claude-sonnet-4-5"}, true},
		{"gpt-4o", []string{"gpt-4*"}, true},
		{"gpt-4o-mini", []string{"gpt-4*"}, true},
		{"claude-opus-4-8", []string{"gpt-4*"}, false},
		{"hello-world", []string{"*world"}, true},
		{"hello-world", []string{"hello-*"}, true},
		{"hello-world", []string{"*lo-wo*"}, true},
	}
	for _, tt := range tests {
		got := matchesModel(tt.model, tt.models)
		if got != tt.want {
			t.Errorf("matchesModel(%q, %v) = %v, want %v", tt.model, tt.models, got, tt.want)
		}
	}
}
