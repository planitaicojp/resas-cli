package config_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/planitaicojp/resas-cli/internal/config"
)

func TestCacheWriteAndRead(t *testing.T) {
	dir := t.TempDir()
	cache := config.NewCache(dir)

	data := []map[string]any{{"prefCode": 13, "prefName": "東京都"}}
	raw, _ := json.Marshal(data)

	if err := cache.Write("prefectures", raw); err != nil {
		t.Fatal(err)
	}

	got, err := cache.Read("prefectures")
	if err != nil {
		t.Fatal(err)
	}

	var result []map[string]any
	if err := json.Unmarshal(got, &result); err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 || result[0]["prefName"] != "東京都" {
		t.Errorf("unexpected: %+v", result)
	}
}

func TestCacheExpired(t *testing.T) {
	dir := t.TempDir()
	cache := config.NewCache(dir)
	cache.TTL = 1 * time.Millisecond

	data := []byte(`[{"prefCode":13}]`)
	if err := cache.Write("test", data); err != nil {
		t.Fatal(err)
	}

	time.Sleep(5 * time.Millisecond)

	_, err := cache.Read("test")
	if err != config.ErrCacheExpired {
		t.Errorf("expected ErrCacheExpired, got %v", err)
	}
}

func TestCacheMiss(t *testing.T) {
	dir := t.TempDir()
	cache := config.NewCache(dir)

	_, err := cache.Read("nonexistent")
	if err != config.ErrCacheMiss {
		t.Errorf("expected ErrCacheMiss, got %v", err)
	}
}

func TestCacheDir(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("RESAS_CACHE_DIR", dir)

	got := config.CacheDir()
	if got != dir {
		t.Errorf("CacheDir() = %q, want %q", got, dir)
	}
}

func TestCacheDelete(t *testing.T) {
	dir := t.TempDir()
	cache := config.NewCache(dir)

	data := []byte(`[{"prefCode":13}]`)
	if err := cache.Write("test", data); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(dir, "test.json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("cache file not created")
	}

	if err := cache.Delete("test"); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("cache file still exists after delete")
	}
}
