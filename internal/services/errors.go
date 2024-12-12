package services

import "errors"

var (
	ErrUserRegistration   = errors.New("Registration: Registration new user error")
	ErrUserAuthentication = errors.New("Authentication: Authentication user error")

	ErrIncorrectOrderNumber       = errors.New("LuhnAlghoritm: Incorrect OrderNumber")
	ErrUnknownValidationAlgorithm = errors.New("ValidateOrderNumber: An unknown algorithm was specified")
)
