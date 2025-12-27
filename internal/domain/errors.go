package domain

import "errors"

var (
	ErrOverflow           = errors.New("integer overflow in bonding curve calculation")
	ErrInvalidQuantity    = errors.New("quantity must be greater than zero")
	ErrInsufficientSupply = errors.New("insufficient supply to sell")
	ErrAlreadyOptedIn     = errors.New("user has already opted in")
	ErrNotOptedIn         = errors.New("user has not opted in")
	ErrMarketNotFound     = errors.New("market not found")
	ErrMarketNotActive    = errors.New("market is not active")
	ErrMarketSunsetting   = errors.New("market is sunsetting, buys are disabled")
	ErrInsufficientFunds  = errors.New("insufficient funds")
	ErrInsufficientShares = errors.New("insufficient shares")
	ErrTradingPaused      = errors.New("trading is paused")
)
