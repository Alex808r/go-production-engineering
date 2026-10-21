package domain

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// OperationType represents the type of transaction (DEPOSIT or WITHDRAW).
type OperationType string

const (
	OperationDeposit  OperationType = "DEPOSIT"
	OperationWithdraw OperationType = "WITHDRAW"
)

// Wallet represents a customer wallet account.
type Wallet struct {
	ID        uuid.UUID       `json:"walletId"`
	Balance   decimal.Decimal `json:"balance"`
	CreatedAt time.Time       `json:"createdAt,omitempty"`
	UpdatedAt time.Time       `json:"updatedAt,omitempty"`
}

// WalletOperationRequest represents the incoming request payload.
// Supports both "valletId" (from spec typo) and "walletId" seamlessly.
type WalletOperationRequest struct {
	WalletID      uuid.UUID       `json:"walletId"`
	OperationType OperationType   `json:"operationType"`
	Amount        decimal.Decimal `json:"amount"`
}

// UnmarshalJSON implements custom JSON unmarshaling to handle both "valletId" and "walletId".
func (r *WalletOperationRequest) UnmarshalJSON(data []byte) error {
	type rawRequest struct {
		WalletID      *uuid.UUID       `json:"walletId"`
		ValletID      *uuid.UUID       `json:"valletId"`
		OperationType string           `json:"operationType"`
		Amount        *decimal.Decimal `json:"amount"`
	}

	var raw rawRequest
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("invalid json payload: %w", err)
	}

	if raw.WalletID != nil {
		r.WalletID = *raw.WalletID
	} else if raw.ValletID != nil {
		r.WalletID = *raw.ValletID
	} else {
		return ErrInvalidWalletID
	}

	if r.WalletID == uuid.Nil {
		return ErrInvalidWalletID
	}

	op := strings.ToUpper(strings.TrimSpace(raw.OperationType))
	switch OperationType(op) {
	case OperationDeposit:
		r.OperationType = OperationDeposit
	case OperationWithdraw:
		r.OperationType = OperationWithdraw
	default:
		return ErrInvalidOperation
	}

	if raw.Amount == nil || raw.Amount.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidAmount
	}
	r.Amount = *raw.Amount

	return nil
}

// WalletResponse represents the response containing wallet information.
type WalletResponse struct {
	WalletID uuid.UUID       `json:"walletId"`
	Balance  decimal.Decimal `json:"balance"`
}
