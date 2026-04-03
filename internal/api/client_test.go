package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/planitaicojp/resas-cli/internal/api"
)

func TestClientGet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-KEY") != "test-key" {
			t.Errorf("X-API-KEY = %q, want %q", r.Header.Get("X-API-KEY"), "test-key")
		}
		if r.Header.Get("User-Agent") == "" {
			t.Error("User-Agent is empty")
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"result": []map[string]any{{"prefCode": 13, "prefName": "東京都"}},
		})
	}))
	defer srv.Close()

	c := api.NewClient("test-key")
	c.BaseURL = srv.URL

	var result struct {
		Result []struct {
			PrefCode int    `json:"prefCode"`
			PrefName string `json:"prefName"`
		} `json:"result"`
	}
	err := c.Get("/api/v1/prefectures", &result)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Result) != 1 {
		t.Fatalf("len(result) = %d, want 1", len(result.Result))
	}
	if result.Result[0].PrefName != "東京都" {
		t.Errorf("prefName = %q, want %q", result.Result[0].PrefName, "東京都")
	}
}

func TestClientAuthError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "Forbidden."})
	}))
	defer srv.Close()

	c := api.NewClient("bad-key")
	c.BaseURL = srv.URL

	var result any
	err := c.Get("/api/v1/prefectures", &result)
	if err == nil {
		t.Fatal("expected error for 403")
	}
}

func TestClientRetryOn429(t *testing.T) {
	attempts := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer srv.Close()

	c := api.NewClient("test-key")
	c.BaseURL = srv.URL

	var result map[string]string
	err := c.Get("/test", &result)
	if err != nil {
		t.Fatal(err)
	}
	if attempts != 3 {
		t.Errorf("attempts = %d, want 3", attempts)
	}
}
