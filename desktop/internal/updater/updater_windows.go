//go:build windows

package updater

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

// platformInstall 在 Windows 上启动 NSIS 安装器并立即退出当前进程，
// 让安装器解除 .exe 文件锁定后完成自更新。
func platformInstall(localPath string) error {
	cmd := exec.Command(localPath)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动 NSIS 安装器失败: %w", err)
	}
	// 让安装器有机会 spawn 自身进程后再退出本进程
	go func() {
		time.Sleep(500 * time.Millisecond)
		os.Exit(0)
	}()
	return nil
}
