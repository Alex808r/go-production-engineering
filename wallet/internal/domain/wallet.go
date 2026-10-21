package domain

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// OperationType задает допустимые типы финансовых транзакций.
type OperationType string

const (
	OperationDeposit  OperationType = "DEPOSIT"  // Операция пополнения баланса
	OperationWithdraw OperationType = "WITHDRAW" // Операция списания средств
)

// Wallet представляет доменную модель учетной записи кошелька.
// -----------------------------------------------------------------------------
// Внимание: для финансового баланса используется decimal.Decimal вместо float64.
// Числа с плавающей точкой в стандарте IEEE 754 не способны точно представлять
// десятичные дроби (например, 0.1 + 0.2 != 0.3), что в банках приводит к потерям денег.
// -----------------------------------------------------------------------------
type Wallet struct {
	ID        uuid.UUID       `json:"walletId"`
	Balance   decimal.Decimal `json:"balance"`
	CreatedAt time.Time       `json:"createdAt,omitempty"`
	UpdatedAt time.Time       `json:"updatedAt,omitempty"`
}

// WalletOperationRequest структура входящего запроса на изменение баланса.
// -----------------------------------------------------------------------------
// Поддерживает как опечатку из технического задания ("valletId"),
// так и общепринятый стандарт ("walletId") через кастомный десериализатор UnmarshalJSON.
// -----------------------------------------------------------------------------
type WalletOperationRequest struct {
	WalletID      uuid.UUID       `json:"walletId"`
	OperationType OperationType   `json:"operationType"`
	Amount        decimal.Decimal `json:"amount"`
}

// UnmarshalJSON реализует гибкую и безопасную валидацию JSON на входе:
// 1. Извлекает walletId либо valletId.
// 2. Проверяет, что UUID валиден и не является нулевым (uuid.Nil).
// 3. Проверяет тип операции (только DEPOSIT или WITHDRAW без учета регистра и пробелов).
// 4. Проверяет, что сумма строго больше нуля (> 0).
func (r *WalletOperationRequest) UnmarshalJSON(data []byte) error {
	type rawRequest struct {
		WalletID      *uuid.UUID       `json:"walletId"`
		ValletID      *uuid.UUID       `json:"valletId"` // поддержка опечатки из текста ТЗ
		OperationType string           `json:"operationType"`
		Amount        *decimal.Decimal `json:"amount"`
	}

	var raw rawRequest
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("invalid json payload: %w", err)
	}

	// Проверяем оба варианта имени поля
	if raw.WalletID != nil {
		r.WalletID = *raw.WalletID
	} else if raw.ValletID != nil {
		r.WalletID = *raw.ValletID
	} else {
		return ErrInvalidWalletID
	}

	// UUID не должен быть нулевым (00000000-0000-0000-0000-000000000000)
	if r.WalletID == uuid.Nil {
		return ErrInvalidWalletID
	}

	// Нормализация и валидация типа операции
	op := strings.ToUpper(strings.TrimSpace(raw.OperationType))
	switch OperationType(op) {
	case OperationDeposit:
		r.OperationType = OperationDeposit
	case OperationWithdraw:
		r.OperationType = OperationWithdraw
	default:
		return ErrInvalidOperation
	}

	// Сумма должна быть передана и быть строго положительной
	if raw.Amount == nil || raw.Amount.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidAmount
	}
	r.Amount = *raw.Amount

	return nil
}

// WalletResponse представляет формат ответа API с балансом кошелька.
type WalletResponse struct {
	WalletID uuid.UUID       `json:"walletId"`
	Balance  decimal.Decimal `json:"balance"`
}
