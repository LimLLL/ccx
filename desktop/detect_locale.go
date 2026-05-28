package main

import (
	"os"
	"strings"
)

// envLocaleFallback 按 POSIX 约定顺序读取 locale 环境变量。
// LANGUAGE 是 GNU 扩展的冒号分隔列表（"zh_CN:en_US"），仅取首项。
// 在 Windows 默认环境中这些变量通常为空，由各平台决定是否再走原生 API。
func envLocaleFallback() string {
	for _, key := range []string{"LANGUAGE", "LC_ALL", "LC_MESSAGES", "LANG"} {
		v := os.Getenv(key)
		if v == "" {
			continue
		}
		if i := strings.Index(v, ":"); i > 0 {
			v = v[:i]
		}
		// "C" / "POSIX" 是 ASCII 兜底，留给 NormalizeLocale 落到英文即可。
		return v
	}
	return ""
}
