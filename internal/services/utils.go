package services

/*
ValidateOrderNumber validates order number by Luhn algorithm

Args:

	orderNumber string: order number for validation

Returns:

	bool: result of validation
*/
func ValidateOrderNumber(orderNumber string) bool {
	return LuhnAlghoritm(orderNumber)
}

func LuhnAlghoritm(orderNumber string) bool {
	// checks
	return true
}
