package account

import (
	"github.com/google/uuid"
)

type CreateAccountRequest struct {
	Currency string `json:"currency"`
}

type UpdateBalanceRequest struct {
	Amount        int64      `json:"amount"`
	OperationType string     `json:"operation_type"`
	ReferenceID   *uuid.UUID `json:"reference_id"`
}

type AccountResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    string    `json:"user_id"`
	Number    string    `json:"number"`
	Balance   int64     `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt string    `json:"created_at"`
}
