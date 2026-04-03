package cmdutil

import (
	"github.com/planitaicojp/resas-cli/internal/api"
	"github.com/planitaicojp/resas-cli/internal/config"
	cerrors "github.com/planitaicojp/resas-cli/internal/errors"
)

func NewClient(flagAPIKey string) (*api.Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	apiKey := config.ResolveAPIKey(flagAPIKey, cfg)
	if apiKey == "" {
		return nil, &cerrors.AuthError{
			Message: "APIキーが設定されていません。以下のいずれかで設定してください:\n  resas config set api_key <KEY>\n  export RESAS_API_KEY=<KEY>\n  --api-key <KEY>",
		}
	}

	return api.NewClient(apiKey), nil
}
