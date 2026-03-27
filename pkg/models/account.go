package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	AccountTypeUser   = "USER"
	AccountTypeSystem = "SYSTEM"
)

// MODEL
type Account struct {
	ID        uuid.UUID  `json:"id"`
	UserID    *string    `json:"user_id"`
	Type      string     `json:"type"`
	Number    string     `json:"number"`
	Balance   int64      `json:"balance"`
	Currency  string     `json:"currency"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// DTO
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
	UserID    string    `json:"user_id,omitempty"`
	Number    string    `json:"number"`
	Balance   int64     `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt string    `json:"created_at"`
}

// EVENTS
type AccountCreatedEvent struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	Number   string `json:"number"`
	Currency string `json:"currency"`
	Email    string `json:"email,omitempty"`
}
