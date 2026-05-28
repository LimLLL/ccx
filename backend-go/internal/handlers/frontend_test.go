package handlers

import "testing"

func TestIsAPIPath(t *testing.T) {
	cases := []struct {
		name string
		path string
		want bool
	}{
		// 裸前缀 —— 本次修复点
		{"bare /v1", "/v1", true},
		{"bare /v1beta", "/v1beta", true},
		{"bare /api", "/api", true},
		{"bare /admin", "/admin", true},

		// 子路径
		{"v1 sub", "/v1/messages", true},
		{"v1 deep sub", "/v1/messages/count_tokens", true},
		{"v1beta sub", "/v1beta/models/gemini-2.0:generateContent", true},
		{"api sub", "/api/messages/channels", true},
		{"admin sub", "/admin/whatever", true},

		// 必须当作非 API（避免误把前端路由扣下）
		{"root", "/", false},
		{"empty", "", false},
		{"v1-prefixed but unrelated", "/v1custom", false},
		{"apifoo", "/apifoo", false},
		{"adminfoo", "/adminfoo", false},
		{"settings page", "/settings", false},
		{"static asset", "/assets/index.js", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := isAPIPath(tc.path)
			if got != tc.want {
				t.Fatalf("isAPIPath(%q) = %v, want %v", tc.path, got, tc.want)
			}
		})
	}
}
