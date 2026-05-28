//go:build !windows && !darwin

package main

// detectSystemLocale 在 Linux/BSD 等 Unix 平台读 POSIX 环境变量链。
func detectSystemLocale() string {
	return envLocaleFallback()
}
