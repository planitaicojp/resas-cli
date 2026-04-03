package api_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/planitaicojp/resas-cli/internal/api"
)

func TestGetPrefectures(t *testing.T) {
	fixture, err := os.ReadFile("../../test/fixtures/prefectures.json")
	if err != nil {
		t.Fatal(err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/prefectures" {
			t.Errorf("path = %q, want /api/v1/prefectures", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(fixture)
	}))
	defer srv.Close()

	c := api.NewClient("test-key")
	c.BaseURL = srv.URL

	area := api.NewAreaAPI(c)
	prefs, err := area.GetPrefectures()
	if err != nil {
		t.Fatal(err)
	}
	if len(prefs) != 3 {
		t.Fatalf("len = %d, want 3", len(prefs))
	}
	if prefs[1].PrefName != "東京都" {
		t.Errorf("prefs[1].PrefName = %q, want %q", prefs[1].PrefName, "東京都")
	}
}

func TestGetCities(t *testing.T) {
	fixture, err := os.ReadFile("../../test/fixtures/cities.json")
	if err != nil {
		t.Fatal(err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/cities" {
			t.Errorf("path = %q, want /api/v1/cities", r.URL.Path)
		}
		if r.URL.Query().Get("prefCode") != "13" {
			t.Errorf("prefCode = %q", r.URL.Query().Get("prefCode"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(fixture)
	}))
	defer srv.Close()

	c := api.NewClient("test-key")
	c.BaseURL = srv.URL

	area := api.NewAreaAPI(c)
	cities, err := area.GetCities(13)
	if err != nil {
		t.Fatal(err)
	}
	if len(cities) != 3 {
		t.Fatalf("len = %d, want 3", len(cities))
	}
	if cities[1].CityName != "千代田区" {
		t.Errorf("cities[1].CityName = %q, want %q", cities[1].CityName, "千代田区")
	}
}
