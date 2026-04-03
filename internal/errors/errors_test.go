package errors_test

import (
	"testing"

	cerrors "github.com/planitaicojp/resas-cli/internal/errors"
)

func TestExitCodes(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{"AuthError", &cerrors.AuthError{Message: "test"}, cerrors.ExitAuth},
		{"NotFoundError", &cerrors.NotFoundError{Message: "test"}, cerrors.ExitNotFound},
		{"ValidationError", &cerrors.ValidationError{Message: "test"}, cerrors.ExitValidation},
		{"APIError", &cerrors.APIError{StatusCode: 500, Message: "test"}, cerrors.ExitAPI},
		{"NetworkError", &cerrors.NetworkError{Err: nil}, cerrors.ExitNetwork},
		{"CancelledError", &cerrors.CancelledError{}, cerrors.ExitCancelled},
		{"generic error", cerrors.New("generic"), cerrors.ExitGeneral},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cerrors.GetExitCode(tt.err)
			if got != tt.want {
				t.Errorf("GetExitCode() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestErrorMessages(t *testing.T) {
	tests := []struct {
		err  error
		want string
	}{
		{&cerrors.AuthError{Message: "APIキーが無効です"}, "エラー: APIキーが無効です"},
		{&cerrors.ValidationError{Message: "都道府県コードが不正です"}, "エラー: 都道府県コードが不正です"},
		{&cerrors.APIError{StatusCode: 429, Message: "レート制限"}, "APIエラー (429): レート制限"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestErrorJSON(t *testing.T) {
	err := &cerrors.AuthError{Message: "APIキーが未設定"}
	j := cerrors.ToJSON(err)
	if j.Error.Code != "auth_error" {
		t.Errorf("Code = %q, want %q", j.Error.Code, "auth_error")
	}
	if j.Error.ExitCode != cerrors.ExitAuth {
		t.Errorf("ExitCode = %d, want %d", j.Error.ExitCode, cerrors.ExitAuth)
	}
}
