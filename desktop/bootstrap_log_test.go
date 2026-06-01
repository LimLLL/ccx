package main

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestDefaultBootstrapLogDirLinux(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Linux bootstrap log dir only applies on Linux")
	}
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	t.Setenv("XDG_STATE_HOME", "")

	want := filepath.Join(tmpHome, ".local", "state", "ccx")
	if got := defaultBootstrapLogDir(); got != want {
		t.Errorf("defaultBootstrapLogDir() = %q, want %q", got, want)
	}
}

func TestDefaultBootstrapLogDirLinuxRespectsXDGStateHome(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Linux bootstrap log dir only applies on Linux")
	}
	stateHome := t.TempDir()
	t.Setenv("XDG_STATE_HOME", stateHome)

	want := filepath.Join(stateHome, "ccx")
	if got := defaultBootstrapLogDir(); got != want {
		t.Errorf("defaultBootstrapLogDir() = %q, want %q", got, want)
	}
}
