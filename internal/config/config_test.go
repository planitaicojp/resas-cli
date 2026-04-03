package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/planitaicojp/resas-cli/internal/config"
)

func TestLoadConfigDefault(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("RESAS_CONFIG_DIR", dir)

	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Defaults.Format != "table" {
		t.Errorf("default format = %q, want %q", cfg.Defaults.Format, "table")
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("RESAS_CONFIG_DIR", dir)

	cfg := &config.Config{
		APIKey: "test-key-123",
		Defaults: config.Defaults{
			Format:   "json",
			PrefCode: 13,
		},
	}
	if err := cfg.Save(); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(dir, "config.json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("config.json not created")
	}

	loaded, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	if loaded.APIKey != "test-key-123" {
		t.Errorf("APIKey = %q, want %q", loaded.APIKey, "test-key-123")
	}
	if loaded.Defaults.Format != "json" {
		t.Errorf("Format = %q, want %q", loaded.Defaults.Format, "json")
	}
	if loaded.Defaults.PrefCode != 13 {
		t.Errorf("PrefCode = %d, want %d", loaded.Defaults.PrefCode, 13)
	}
}

func TestConfigDir(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("RESAS_CONFIG_DIR", dir)
	if got := config.Dir(); got != dir {
		t.Errorf("Dir() = %q, want %q", got, dir)
	}
}

func TestResolveAPIKey(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("RESAS_CONFIG_DIR", dir)

	cfg := &config.Config{APIKey: "from-config"}

	t.Setenv("RESAS_API_KEY", "from-env")
	got := config.ResolveAPIKey("", cfg)
	if got != "from-env" {
		t.Errorf("got %q, want %q", got, "from-env")
	}

	got = config.ResolveAPIKey("from-flag", cfg)
	if got != "from-flag" {
		t.Errorf("got %q, want %q", got, "from-flag")
	}

	t.Setenv("RESAS_API_KEY", "")
	got = config.ResolveAPIKey("", cfg)
	if got != "from-config" {
		t.Errorf("got %q, want %q", got, "from-config")
	}
}
