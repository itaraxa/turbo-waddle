package services

import (
	"errors"

	e "github.com/itaraxa/turbo-waddle/internal/errors"
)

const (
	LUHN = `Luhn`
	NONE = `None`
)

/*
ValidateOrderNumber validates order number by specified algorithm

Args:

	orderNumber string: order number for validation
	algorithm string: algorithm used for validation (Luhn or None). None returns TRUE always

Returns:

	result bool: result of validation
	err error
*/
func ValidateOrderNumber(orderNumber string, algorithm string) (result bool, err error) {
	switch algorithm {
	case LUHN:
		return LuhnAlghoritm(orderNumber)
	case NONE:
		return true, nil
	default:
		return false, errors.Join(ErrUnknownValidationAlgorithm, e.ErrInternalServerError)
	}
}

/*
LuhnAlghoritm validates the orderNumber using the Luhn algorithm

Args:

	orderNumber string: order number as a string. This allows for the use of long numbers

Returns:

	result bool: True - if order number passed the validation, False - if the number failed the validation or an error occurred
	err error: nil or error if the order number is too short or contains non-numeric characters
*/
func LuhnAlghoritm(orderNumber string) (result bool, err error) {
	if len(orderNumber) < 2 {
		return false, errors.Join(ErrIncorrectOrderNumber,
			e.ErrInvalidOrderNumberForamt,
			errors.New("order number is too short"),
		)
	}
	seed := orderNumber[:len(orderNumber)-1]
	sum, parity := 0, len(seed)%2
	for i, n := range seed {
		if isNotNumber(n) {
			return false, errors.Join(ErrIncorrectOrderNumber,
				e.ErrInvalidOrderNumberForamt,
				errors.New("order number contains non-numeric symbol"),
			)
		}
		d := int(n - '0')
		if i%2 != parity {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
	}

	result = (sum+int(orderNumber[len(orderNumber)-1]-'0'))%10 == 0
	return
}

/*
isNotNumber checks if the specified rune is a digit

Args:

	n rune

Returns:

	bool
*/
func isNotNumber(n rune) bool {
	return n < '0' || '9' < n
}
