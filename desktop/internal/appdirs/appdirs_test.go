package appdirs

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestDataDirLinux(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Linux default data dir only applies on Linux")
	}
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	t.Setenv("XDG_STATE_HOME", "")

	want := filepath.Join(tmpHome, ".local", "state", "ccx")
	if got := DataDir(); got != want {
		t.Errorf("DataDir() = %q, want %q", got, want)
	}
}

func TestDataDirLinuxRespectsXDGStateHome(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Linux default data dir only applies on Linux")
	}
	stateHome := t.TempDir()
	t.Setenv("XDG_STATE_HOME", stateHome)

	want := filepath.Join(stateHome, "ccx")
	if got := DataDir(); got != want {
		t.Errorf("DataDir() = %q, want %q", got, want)
	}
}

func TestDataDirForHome(t *testing.T) {
	home := t.TempDir()
	if runtime.GOOS == "linux" {
		t.Setenv("XDG_STATE_HOME", "")
	}

	want := filepath.Join(home, ".config", "ccx-desktop")
	if runtime.GOOS == "linux" {
		want = filepath.Join(home, ".local", "state", "ccx")
	}
	if got := DataDirForHome(home); got != want {
		t.Errorf("DataDirForHome() = %q, want %q", got, want)
	}
}

func TestDataDirForHomeLinuxRespectsXDGStateHome(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Linux default data dir only applies on Linux")
	}
	stateHome := t.TempDir()
	t.Setenv("XDG_STATE_HOME", stateHome)

	want := filepath.Join(stateHome, "ccx")
	if got := DataDirForHome(t.TempDir()); got != want {
		t.Errorf("DataDirForHome() = %q, want %q", got, want)
	}
}
