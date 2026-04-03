package errors

import "fmt"

const (
	ExitOK         = 0
	ExitGeneral    = 1
	ExitAuth       = 2
	ExitNotFound   = 3
	ExitValidation = 4
	ExitAPI        = 5
	ExitNetwork    = 6
	ExitCancelled  = 10
)

type ExitCoder interface {
	ExitCode() int
}

type AuthError struct {
	Message string
}

func (e *AuthError) Error() string     { return "エラー: " + e.Message }
func (e *AuthError) ExitCode() int     { return ExitAuth }
func (e *AuthError) ErrorCode() string { return "auth_error" }

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string     { return "エラー: " + e.Message }
func (e *NotFoundError) ExitCode() int     { return ExitNotFound }
func (e *NotFoundError) ErrorCode() string { return "not_found" }

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string     { return "エラー: " + e.Message }
func (e *ValidationError) ExitCode() int     { return ExitValidation }
func (e *ValidationError) ErrorCode() string { return "validation_error" }

type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string     { return fmt.Sprintf("APIエラー (%d): %s", e.StatusCode, e.Message) }
func (e *APIError) ExitCode() int     { return ExitAPI }
func (e *APIError) ErrorCode() string { return "api_error" }

type NetworkError struct {
	Err error
}

func (e *NetworkError) Error() string {
	if e.Err != nil {
		return "ネットワークエラー: " + e.Err.Error()
	}
	return "ネットワークエラー"
}
func (e *NetworkError) ExitCode() int     { return ExitNetwork }
func (e *NetworkError) ErrorCode() string { return "network_error" }

type CancelledError struct{}

func (e *CancelledError) Error() string     { return "キャンセルされました" }
func (e *CancelledError) ExitCode() int     { return ExitCancelled }
func (e *CancelledError) ErrorCode() string { return "cancelled" }

func New(msg string) error {
	return fmt.Errorf("エラー: %s", msg)
}

func GetExitCode(err error) int {
	if err == nil {
		return ExitOK
	}
	if ec, ok := err.(ExitCoder); ok {
		return ec.ExitCode()
	}
	return ExitGeneral
}

type ErrorJSON struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	ExitCode int    `json:"exit_code"`
}

type errorCoder interface {
	ErrorCode() string
}

func ToJSON(err error) ErrorJSON {
	code := "general_error"
	if ec, ok := err.(errorCoder); ok {
		code = ec.ErrorCode()
	}
	return ErrorJSON{
		Error: ErrorDetail{
			Code:     code,
			Message:  err.Error(),
			ExitCode: GetExitCode(err),
		},
	}
}
