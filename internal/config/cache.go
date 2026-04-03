package config

import (
	"errors"
	"os"
	"path/filepath"
	"time"
)

const (
	EnvCacheDir     = "RESAS_CACHE_DIR"
	defaultCacheTTL = 30 * 24 * time.Hour
)

var (
	ErrCacheMiss    = errors.New("キャッシュが見つかりません")
	ErrCacheExpired = errors.New("キャッシュの有効期限が切れています")
)

func CacheDir() string {
	if d := os.Getenv(EnvCacheDir); d != "" {
		return d
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".cache", "resas")
}

type Cache struct {
	Dir string
	TTL time.Duration
}

func NewCache(dir string) *Cache {
	return &Cache{
		Dir: dir,
		TTL: defaultCacheTTL,
	}
}

func (c *Cache) Read(key string) ([]byte, error) {
	path := c.path(key)

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrCacheMiss
		}
		return nil, err
	}

	if time.Since(info.ModTime()) > c.TTL {
		return nil, ErrCacheExpired
	}

	return os.ReadFile(path)
}

func (c *Cache) Write(key string, data []byte) error {
	if err := os.MkdirAll(c.Dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(c.path(key), data, 0644)
}

func (c *Cache) Delete(key string) error {
	err := os.Remove(c.path(key))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func (c *Cache) path(key string) string {
	return filepath.Join(c.Dir, key+".json")
}
