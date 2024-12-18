package accrual

import "fmt"

type accrualError struct {
	Msg        string
	Code       int
	RetryAfter int
}

func (e *accrualError) Error() string {
	return fmt.Sprintf("error: %d: %s", e.Code, e.Msg)
}

var (
	ErrOrderDoesNotRegistered     = &accrualError{Msg: "order does not registered in accrual system", Code: 204}
	ErrServerRequestLimitExceeded = &accrualError{Msg: "server request limit exceeded", Code: 429, RetryAfter: 0}
	ErrInternalServerError        = &accrualError{Msg: "internal server error", Code: 500}
	ErrUnknownStatus              = &accrualError{Msg: "unknown status", Code: 500}
)
