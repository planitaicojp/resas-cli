package prompt_test

import (
	"os"
	"testing"

	"github.com/planitaicojp/resas-cli/internal/prompt"
)

func TestSelectItemLabel(t *testing.T) {
	item := prompt.SelectItem{Label: "13: 東京都", Value: "13"}
	if item.Label != "13: 東京都" {
		t.Errorf("Label = %q, want %q", item.Label, "13: 東京都")
	}
	if item.Value != "13" {
		t.Errorf("Value = %q, want %q", item.Value, "13")
	}
}

func TestIsNoInput(t *testing.T) {
	t.Setenv("RESAS_NO_INPUT", "1")
	if !prompt.IsNoInput() {
		t.Error("IsNoInput() should be true when RESAS_NO_INPUT=1")
	}

	t.Setenv("RESAS_NO_INPUT", "")
	// In test env, stdin is not TTY so IsNoInput may still effectively be true
}

func TestIsTTY(t *testing.T) {
	if prompt.IsTTY(os.Stdin) {
		t.Skip("stdin is a TTY in this environment")
	}
}
