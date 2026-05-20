//go:build linux

package updater

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// platformInstall 替换当前 AppImage 文件，再 exec 新的。
//
// 流程：
//  1. 检测当前可执行路径（os.Executable）
//  2. 把下载的新 AppImage 复制到 <exe>.new
//  3. chmod +x
//  4. 启动 shell 脚本：sleep 1 && mv <new> <exe> && exec <exe>
//  5. 当前进程立即退出，腾出文件锁
//
// 注意：仅支持 AppImage 形态；如果用户用 deb/rpm 安装则跳过自动替换。
func platformInstall(localPath string) error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("获取当前可执行路径失败: %w", err)
	}

	// 简单启发式：可执行不在 /usr/* 才认为是 AppImage 直装模式
	if strings.HasPrefix(exePath, "/usr/") {
		return fmt.Errorf("检测到系统安装版本，请通过包管理器更新")
	}

	newPath := exePath + ".new"
	if err := copyFile(localPath, newPath); err != nil {
		return fmt.Errorf("复制新版本失败: %w", err)
	}
	if err := os.Chmod(newPath, 0o755); err != nil {
		return fmt.Errorf("设置可执行权限失败: %w", err)
	}

	script := fmt.Sprintf(`sleep 1 && mv -f %q %q && exec %q &`, newPath, exePath, exePath)
	if err := exec.Command("sh", "-c", script).Start(); err != nil {
		return fmt.Errorf("启动替换脚本失败: %w", err)
	}
	go func() {
		os.Exit(0)
	}()
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o755)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
