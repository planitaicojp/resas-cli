package api_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/planitaicojp/resas-cli/internal/api"
)

func TestGetPopulationComposition(t *testing.T) {
	fixture, err := os.ReadFile("../../test/fixtures/population_composition.json")
	if err != nil {
		t.Fatal(err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/population/composition/perYear" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("prefCode") != "13" {
			t.Errorf("prefCode = %q", r.URL.Query().Get("prefCode"))
		}
		if r.URL.Query().Get("cityCode") != "-" {
			t.Errorf("cityCode = %q", r.URL.Query().Get("cityCode"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(fixture)
	}))
	defer srv.Close()

	c := api.NewClient("test-key")
	c.BaseURL = srv.URL

	pop := api.NewPopulationAPI(c)
	result, err := pop.GetComposition(13, "-")
	if err != nil {
		t.Fatal(err)
	}
	if result.BoundaryYear != 2020 {
		t.Errorf("BoundaryYear = %d, want 2020", result.BoundaryYear)
	}
	if len(result.Data) != 2 {
		t.Fatalf("len(Data) = %d, want 2", len(result.Data))
	}
	if result.Data[0].Label != "総人口" {
		t.Errorf("Label = %q, want %q", result.Data[0].Label, "総人口")
	}
}
