package domain

import "errors"

var (
	// ErrWalletNotFound is returned when the specified wallet does not exist.
	ErrWalletNotFound = errors.New("wallet not found")

	// ErrInsufficientFunds is returned when attempting to withdraw more than the available balance.
	ErrInsufficientFunds = errors.New("insufficient funds")

	// ErrInvalidOperation is returned when the operation type is not DEPOSIT or WITHDRAW.
	ErrInvalidOperation = errors.New("invalid operation type: must be DEPOSIT or WITHDRAW")

	// ErrInvalidAmount is returned when the amount is less than or equal to zero.
	ErrInvalidAmount = errors.New("amount must be greater than zero")

	// ErrInvalidWalletID is returned when the wallet UUID is invalid or nil.
	ErrInvalidWalletID = errors.New("invalid wallet id format")
)
