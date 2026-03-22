package transaction

import "github.com/google/uuid"

type CreateTransactionRequest struct {
	ToAccountNumber string `json:"to_account_number"`
	Amount          int64  `json:"amount"`
	Currency        string `json:"currency"`
	IdempotencyKey  string `json:"idempotency_key"`
}

type TransactionResponse struct {
	ID             uuid.UUID `json:"id"`
	FromAccountID  uuid.UUID `json:"from_account_id"`
	ToAccountID    uuid.UUID `json:"to_account_id"`
	Amount         int64     `json:"amount"`
	Currency       string    `json:"currency"`
	Status         string    `json:"status"`
	IdempotencyKey string    `json:"idempotency_key,omitempty"`
	CreatedAt      string    `json:"created_at"`
}
