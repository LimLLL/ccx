//go:build darwin

package updater

import (
	"fmt"
	"os/exec"
)

// platformInstall 在 macOS 上调用 `open` 让 Finder 接管 DMG 挂载。
//
// 用户从挂载窗口拖拽 CCX Desktop.app 到 Applications 完成安装。
// 我们不直接替换 .app（避免 ad-hoc 签名残留 + 文件占用问题）。
func platformInstall(localPath string) error {
	cmd := exec.Command("/usr/bin/open", localPath)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("打开 DMG 失败: %w", err)
	}
	return nil
}
