package domain

import "errors"

var (
	ErrOverflow           = errors.New("integer overflow in bonding curve calculation")
	ErrInvalidQuantity    = errors.New("quantity must be greater than zero")
	ErrInsufficientSupply = errors.New("insufficient supply to sell")
)
