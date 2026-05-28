//go:build windows

package main

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// LOCALE_NAME_MAX_LENGTH，含 null 终结符的 wchar_t 数。
const localeNameMaxLength = 85

// detectSystemLocale 通过 Win32 GetUserDefaultLocaleName 读取用户默认区域，
// 返回 BCP-47 形式（如 "zh-CN"、"en-US"），失败时回落 POSIX 环境变量。
//
// Windows 默认不设 LC_ALL/LC_MESSAGES/LANG，仅依赖环境变量会导致首启英文。
func detectSystemLocale() string {
	if v := envLocaleFallback(); v != "" {
		return v
	}

	dll := windows.NewLazySystemDLL("kernel32.dll")
	proc := dll.NewProc("GetUserDefaultLocaleName")

	buf := make([]uint16, localeNameMaxLength)
	ret, _, _ := proc.Call(
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(len(buf)),
	)
	if ret == 0 {
		return ""
	}
	// ret 是写入的 wchar_t 数（含末尾 NUL），UTF16ToString 会自行截断。
	return syscall.UTF16ToString(buf[:ret])
}
