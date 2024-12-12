package services

import "errors"

var (
	ErrIncorrectOrderNumber       = errors.New("LuhnAlghoritm: Incorrect OrderNumber")
	ErrUnknownValidationAlgorithm = errors.New("ValidateOrderNumber: An unknown algorithm was specified")
)
