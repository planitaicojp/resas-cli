package output_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/planitaicojp/resas-cli/internal/output"
)

type testRow struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func TestJSONFormatter(t *testing.T) {
	var buf bytes.Buffer
	data := []testRow{{Name: "東京都", Value: 13}}
	err := output.New("json").Format(&buf, data)
	if err != nil {
		t.Fatal(err)
	}
	var result []testRow
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(result) != 1 || result[0].Name != "東京都" {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestTableFormatter(t *testing.T) {
	var buf bytes.Buffer
	data := []testRow{{Name: "東京都", Value: 13}}
	err := output.New("table").Format(&buf, data)
	if err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "NAME") {
		t.Errorf("missing header NAME in: %s", out)
	}
	if !strings.Contains(out, "東京都") {
		t.Errorf("missing data in: %s", out)
	}
}

func TestCSVFormatter(t *testing.T) {
	var buf bytes.Buffer
	data := []testRow{{Name: "東京都", Value: 13}}
	err := output.New("csv").Format(&buf, data)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d: %s", len(lines), buf.String())
	}
	if !strings.Contains(lines[0], "name") {
		t.Errorf("missing header in: %s", lines[0])
	}
	if !strings.Contains(lines[1], "東京都") {
		t.Errorf("missing data in: %s", lines[1])
	}
}

func TestFormatterDefault(t *testing.T) {
	f := output.New("unknown")
	var buf bytes.Buffer
	err := f.Format(&buf, []testRow{{Name: "test", Value: 1}})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "NAME") {
		t.Error("default should be table format")
	}
}

func TestEmptySlice(t *testing.T) {
	for _, format := range []string{"json", "table", "csv"} {
		var buf bytes.Buffer
		err := output.New(format).Format(&buf, []testRow{})
		if err != nil {
			t.Errorf("%s: unexpected error: %v", format, err)
		}
	}
}
