package errors

import "fmt"

type serverError struct {
	Msg  string
	Code int
}

func (e *serverError) Error() string {
	return fmt.Sprintf("error: %d: %s", e.Code, e.Msg)
}

var (
	// Register errors
	ErrInvalidRequestFormat = &serverError{Msg: `Invalid request format`, Code: 400}
	ErrLoginIsAlreadyUsed   = &serverError{Msg: `Login already used`, Code: 409}
	ErrInternalServerError  = &serverError{Msg: `Internal server error`, Code: 500}
	// Login errors
	ErrInvalidLoginPassPair = &serverError{Msg: `Invalid login/password pair`, Code: 401}
	// Orders errors
	ErrUserIsNotauthenticated   = &serverError{Msg: `User is not authenticated`, Code: 401}
	ErrOrderAlreadyUploaded     = &serverError{Msg: `Order number has already been uploaded by another user`, Code: 409}
	ErrInvalidOrderNumberForamt = &serverError{Msg: `Invalid order number format`, Code: 422}
	ErrNoData                   = &serverError{Msg: `No data for responce`, Code: 204}
	// Balance errors
	ErrInsufficientFunds  = &serverError{Msg: `Insufficient funds in the account`, Code: 402}
	ErrInvalidOrderNumber = &serverError{Msg: `Invalid order number`, Code: 422}
	// Withdraws errors
	ErrNoWithdraws = &serverError{Msg: `No withdraws`, Code: 204}
)
