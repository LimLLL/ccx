package updater

import "testing"

func TestCompareSemver(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"1.0.0", "1.0.0", 0},
		{"1.0.1", "1.0.0", 1},
		{"1.0.0", "1.0.1", -1},
		{"2.0.0", "1.99.99", 1},
		{"1.0.0", "1.0.0-rc1", 1},
		{"1.0.0-rc2", "1.0.0-rc1", 1},
		{"1.0.0-rc1", "1.0.0", -1},
		{"1.2", "1.2.0", 0},
		{"1.2.3.4", "1.2.3", 1},
	}
	for _, c := range cases {
		got, err := compareSemver(c.a, c.b)
		if err != nil {
			t.Fatalf("compareSemver(%q,%q): unexpected error: %v", c.a, c.b, err)
		}
		if got != c.want {
			t.Errorf("compareSemver(%q,%q) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestCompareSemverError(t *testing.T) {
	if _, err := compareSemver("abc", "1.0.0"); err == nil {
		t.Error("expected error on invalid semver, got nil")
	}
}

func TestPlatformAssetName(t *testing.T) {
	name := PlatformAssetName("2.7.9")
	if name == "" {
		t.Error("PlatformAssetName 返回空")
	}
}
