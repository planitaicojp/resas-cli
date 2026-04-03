package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	EnvConfigDir = "RESAS_CONFIG_DIR"
	EnvAPIKey    = "RESAS_API_KEY"
	EnvFormat    = "RESAS_FORMAT"
	EnvNoInput   = "RESAS_NO_INPUT"
)

type Config struct {
	APIKey   string   `json:"api_key,omitempty"`
	Defaults Defaults `json:"defaults"`
}

type Defaults struct {
	Format   string `json:"format,omitempty"`
	PrefCode int    `json:"pref_code,omitempty"`
}

func Dir() string {
	if d := os.Getenv(EnvConfigDir); d != "" {
		return d
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "resas")
}

func Path() string {
	return filepath.Join(Dir(), "config.json")
}

func Load() (*Config, error) {
	cfg := &Config{
		Defaults: Defaults{
			Format: "table",
		},
	}

	data, err := os.ReadFile(Path())
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	if cfg.Defaults.Format == "" {
		cfg.Defaults.Format = "table"
	}
	return cfg, nil
}

func (c *Config) Save() error {
	dir := Dir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(Path(), data, 0644)
}

func ResolveAPIKey(flagValue string, cfg *Config) string {
	if flagValue != "" {
		return flagValue
	}
	if v := os.Getenv(EnvAPIKey); v != "" {
		return v
	}
	return cfg.APIKey
}

func ResolveFormat(flagValue string, cfg *Config) string {
	if flagValue != "" {
		return flagValue
	}
	if v := os.Getenv(EnvFormat); v != "" {
		return v
	}
	if cfg.Defaults.Format != "" {
		return cfg.Defaults.Format
	}
	return "table"
}
