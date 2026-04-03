package prompt

import (
	"os"

	"github.com/manifoldco/promptui"
	cerrors "github.com/planitaicojp/resas-cli/internal/errors"
)

type SelectItem struct {
	Label string
	Value string
}

func IsNoInput() bool {
	return os.Getenv("RESAS_NO_INPUT") == "1"
}

var NoInputFlag bool

func IsTTY(f *os.File) bool {
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) != 0
}

func Select(label string, items []SelectItem) (string, error) {
	if NoInputFlag || IsNoInput() || !IsTTY(os.Stdin) {
		return "", &cerrors.ValidationError{
			Message: label + "が指定されていません。--pref-code フラグで指定してください。",
		}
	}

	labels := make([]string, len(items))
	for i, item := range items {
		labels[i] = item.Label
	}

	p := promptui.Select{
		Label: label,
		Items: labels,
		Size:  15,
		Searcher: func(input string, index int) bool {
			return containsIgnoreCase(labels[index], input)
		},
		Stdout: os.Stderr,
	}

	idx, _, err := p.Run()
	if err != nil {
		if err == promptui.ErrInterrupt {
			return "", &cerrors.CancelledError{}
		}
		return "", err
	}

	return items[idx].Value, nil
}

func containsIgnoreCase(s, substr string) bool {
	sLower := []rune(s)
	subLower := []rune(substr)
	if len(subLower) > len(sLower) {
		return false
	}
	for i := 0; i <= len(sLower)-len(subLower); i++ {
		match := true
		for j := range subLower {
			a, b := sLower[i+j], subLower[j]
			if a != b && toLower(a) != toLower(b) {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

func toLower(r rune) rune {
	if r >= 'A' && r <= 'Z' {
		return r + 32
	}
	return r
}
