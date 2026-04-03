package cmdutil

import "github.com/planitaicojp/resas-cli/internal/config"

func GetFormat(flagFormat string) string {
	cfg, _ := config.Load()
	return config.ResolveFormat(flagFormat, cfg)
}
