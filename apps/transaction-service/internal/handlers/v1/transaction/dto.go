package transaction

import "github.com/google/uuid"

type CreateTransactionRequest struct {
	Type            string `json:"type,omitempty"`
	ToAccountNumber string `json:"to_account_number"`
	Amount          int64  `json:"amount"`
	Currency        string `json:"currency"`
	IdempotencyKey  string `json:"idempotency_key"`
}

type TransferRequest struct {
	ToAccountNumber string `json:"to_account_number"`
	Amount          int64  `json:"amount"`
	Currency        string `json:"currency"`
	IdempotencyKey  string `json:"idempotency_key"`
}

type DepositRequest struct {
	Amount         int64  `json:"amount"`
	Currency       string `json:"currency"`
	IdempotencyKey string `json:"idempotency_key"`
}

type WithdrawRequest struct {
	Amount         int64  `json:"amount"`
	Currency       string `json:"currency"`
	IdempotencyKey string `json:"idempotency_key"`
}

type ExchangeRequest struct {
	Amount         int64  `json:"amount"`
	FromCurrency   string `json:"from_currency"`
	ToCurrency     string `json:"to_currency"`
	IdempotencyKey string `json:"idempotency_key"`
}

type TransactionResponse struct {
	ID             uuid.UUID `json:"id"`
	Type           string    `json:"type"`
	FromAccountID  uuid.UUID `json:"from_account_id"`
	ToAccountID    uuid.UUID `json:"to_account_id"`
	Amount         int64     `json:"amount"`
	Currency       string    `json:"currency"`
	Status         string    `json:"status"`
	IdempotencyKey string    `json:"idempotency_key,omitempty"`
	CreatedAt      string    `json:"created_at"`
}
