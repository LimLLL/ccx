package windowstate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadMissingFile(t *testing.T) {
	dir := t.TempDir()
	st, ok, err := Load(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Errorf("expected ok=false for missing file, got %+v", st)
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	in := State{X: 100, Y: 200, Width: 1180, Height: 820, Maximised: false}
	if err := Save(dir, in); err != nil {
		t.Fatalf("save: %v", err)
	}
	out, ok, err := Load(dir)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if !ok {
		t.Fatal("expected ok=true after save")
	}
	if out != in {
		t.Errorf("roundtrip mismatch: got %+v want %+v", out, in)
	}
}

func TestLoadCorruptFileTreatedAsMissing(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "window-state.json"), []byte("not-json"), 0o644); err != nil {
		t.Fatalf("write corrupt: %v", err)
	}
	_, ok, err := Load(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected ok=false on corrupt file")
	}
}

func TestIsValidRejects(t *testing.T) {
	cases := []State{
		{Width: 100, Height: 100},   // too small
		{Width: 99999, Height: 800}, // too large
		{Width: 800, Height: 600, X: -99999, Y: 0},
	}
	for i, c := range cases {
		if IsValid(c) {
			t.Errorf("case %d: expected invalid, got valid: %+v", i, c)
		}
	}
}

func TestIsValidAccepts(t *testing.T) {
	good := State{X: 0, Y: 0, Width: 1180, Height: 820}
	if !IsValid(good) {
		t.Errorf("expected valid, got invalid: %+v", good)
	}
}
