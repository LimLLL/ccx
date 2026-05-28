//go:build darwin

package main

import (
	"os/exec"
	"strings"
)

// detectSystemLocale 在 macOS 优先读 POSIX 环境变量；
// GUI 应用从 Finder/Dock 启动通常不继承 shell 环境，
// 此时回退 `defaults read -g AppleLocale`，输出形如 "zh_CN"。
func detectSystemLocale() string {
	if v := envLocaleFallback(); v != "" {
		return v
	}
	out, err := exec.Command("defaults", "read", "-g", "AppleLocale").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
