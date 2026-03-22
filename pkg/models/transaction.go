package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	TransactionStatusPending    = "PENDING"
	TransactionStatusHolding    = "HOLDING"
	TransactionStatusDepositing = "DEPOSITING"
	TransactionStatusCompleted  = "COMPLETED"
	TransactionStatusFailed     = "FAILED"
)

const (
	TransactionTypeTransfer   = "TRANSFER"
	TransactionTypeDeposit    = "DEPOSIT"
	TransactionTypeWithdrawal = "WITHDRAWAL"
)

func IsValidTransactionType(transactionType string) bool {
	switch transactionType {
	case TransactionTypeTransfer, TransactionTypeDeposit, TransactionTypeWithdrawal:
		return true
	default:
		return false
	}
}

type Transaction struct {
	ID              uuid.UUID  `json:"id"`
	Type            string     `json:"type"` // TRANSFER, DEPOSIT, WITHDRAWAL
	FromAccountID   uuid.UUID  `json:"from_account_id"`
	ToAccountID     uuid.UUID  `json:"to_account_id"`
	ToAccountNumber string     `json:"to_account_number,omitempty"`
	Amount          int64      `json:"amount"`
	Currency        string     `json:"currency"`
	Status          string     `json:"status"`
	IdempotencyKey  string     `json:"idempotency_key,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty"`
}

type TransactionCreatedEvent struct {
	ID            string `json:"id"`
	Type          string `json:"type"` // TRANSFER, DEPOSIT, WITHDRAWAL
	FromAccountID string `json:"from_account_id"`
	ToAccountID   string `json:"to_account_id"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
}

type TransactionStatusChangedEvent struct {
	ID     string `json:"id"`
	Status string `json:"status"` // HOLDING, COMPLETED, FAILED
}
