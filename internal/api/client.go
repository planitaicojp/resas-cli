package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	cerrors "github.com/planitaicojp/resas-cli/internal/errors"
)

const (
	defaultBaseURL = "https://opendata.resas-portal.go.jp"
	maxRetries     = 3
)

var UserAgent = "planitai/resas-cli/dev"
var Verbose bool

type Client struct {
	HTTP    *http.Client
	APIKey  string
	BaseURL string
}

func NewClient(apiKey string) *Client {
	return &Client{
		HTTP: &http.Client{
			Timeout: 30 * time.Second,
		},
		APIKey:  apiKey,
		BaseURL: defaultBaseURL,
	}
}

func (c *Client) Get(path string, result any) error {
	url := c.BaseURL + path

	var resp *http.Response
	var err error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		req, reqErr := http.NewRequest("GET", url, nil)
		if reqErr != nil {
			return &cerrors.NetworkError{Err: reqErr}
		}
		req.Header.Set("X-API-KEY", c.APIKey)
		req.Header.Set("User-Agent", UserAgent)

		if Verbose {
			fmt.Fprintf(os.Stderr, ">>> GET %s\n", url)
		}

		resp, err = c.HTTP.Do(req)
		if err != nil {
			if attempt == maxRetries {
				return &cerrors.NetworkError{Err: err}
			}
			time.Sleep(time.Duration(attempt+1) * time.Second)
			continue
		}

		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
			_ = resp.Body.Close()
			if attempt < maxRetries {
				time.Sleep(time.Duration(attempt+1) * time.Second)
				continue
			}
		}
		break
	}

	defer func() { _ = resp.Body.Close() }()

	if Verbose {
		fmt.Fprintf(os.Stderr, "<<< %d %s\n", resp.StatusCode, resp.Status)
	}

	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized {
		return &cerrors.AuthError{Message: "APIキーが無効です。resas config set api_key <KEY> で正しいキーを設定してください。"}
	}

	if resp.StatusCode == http.StatusNotFound {
		return &cerrors.NotFoundError{Message: "指定されたデータが見つかりません。"}
	}

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return &cerrors.APIError{
			StatusCode: resp.StatusCode,
			Message:    string(body),
		}
	}

	return json.NewDecoder(resp.Body).Decode(result)
}
