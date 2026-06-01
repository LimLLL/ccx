package appdirs

import (
	"os"
	"path/filepath"
	"runtime"
)

const (
	linuxAppDir = "ccx"
	appDir      = "ccx-desktop"
)

func DataDir() string {
	if runtime.GOOS == "linux" {
		if base := linuxStateBaseDir(); base != "" {
			return filepath.Join(base, linuxAppDir)
		}
	}

	base, err := os.UserConfigDir()
	if err != nil || base == "" {
		base, _ = os.UserHomeDir()
	}
	if base == "" {
		base = "."
	}
	return filepath.Join(base, appDir)
}

func DataDirForHome(homeDir string) string {
	if runtime.GOOS == "linux" {
		base := os.Getenv("XDG_STATE_HOME")
		if base == "" && homeDir != "" {
			base = filepath.Join(homeDir, ".local", "state")
		}
		if base != "" {
			return filepath.Join(base, linuxAppDir)
		}
	}

	return filepath.Join(homeDir, ".config", appDir)
}

func linuxStateBaseDir() string {
	base := os.Getenv("XDG_STATE_HOME")
	if base != "" {
		return base
	}
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return ""
	}
	return filepath.Join(home, ".local", "state")
}
