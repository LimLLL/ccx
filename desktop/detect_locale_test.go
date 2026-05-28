package main

import (
	"os"
	"testing"
)

func TestEnvLocaleFallback(t *testing.T) {
	// 备份并清空所有相关环境变量，测试结束后恢复。
	keys := []string{"LANGUAGE", "LC_ALL", "LC_MESSAGES", "LANG"}
	saved := make(map[string]string, len(keys))
	for _, k := range keys {
		saved[k] = os.Getenv(k)
		_ = os.Unsetenv(k)
	}
	t.Cleanup(func() {
		for _, k := range keys {
			if v, ok := saved[k]; ok && v != "" {
				_ = os.Setenv(k, v)
			} else {
				_ = os.Unsetenv(k)
			}
		}
	})

	cases := []struct {
		name string
		env  map[string]string
		want string
	}{
		{
			name: "all empty returns empty",
			env:  map[string]string{},
			want: "",
		},
		{
			name: "LANG only",
			env:  map[string]string{"LANG": "en_US.UTF-8"},
			want: "en_US.UTF-8",
		},
		{
			name: "LC_ALL beats LANG",
			env:  map[string]string{"LANG": "en_US.UTF-8", "LC_ALL": "zh_CN.UTF-8"},
			want: "zh_CN.UTF-8",
		},
		{
			name: "LANGUAGE highest priority",
			env: map[string]string{
				"LANG":     "en_US.UTF-8",
				"LC_ALL":   "en_US.UTF-8",
				"LANGUAGE": "zh_CN",
			},
			want: "zh_CN",
		},
		{
			name: "LANGUAGE colon list takes first",
			env: map[string]string{
				"LANGUAGE": "zh_CN:en_US:ja_JP",
			},
			want: "zh_CN",
		},
		{
			name: "LC_MESSAGES between LC_ALL and LANG",
			env: map[string]string{
				"LANG":         "en_US",
				"LC_MESSAGES":  "fr_FR",
			},
			want: "fr_FR",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, k := range keys {
				_ = os.Unsetenv(k)
			}
			for k, v := range tc.env {
				_ = os.Setenv(k, v)
			}
			if got := envLocaleFallback(); got != tc.want {
				t.Fatalf("envLocaleFallback() = %q, want %q", got, tc.want)
			}
		})
	}
}
