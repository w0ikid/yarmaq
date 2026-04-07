package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/pkg/models"
)

type Transaction struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id" json:"id"`
	Type           string     `gorm:"type:varchar(20);not null;default:'TRANSFER';column:type" json:"type"`
	FromAccountID  uuid.UUID  `gorm:"type:uuid;not null;column:from_account_id" json:"from_account_id"`
	ToAccountID    uuid.UUID  `gorm:"type:uuid;not null;column:to_account_id" json:"to_account_id"`
	Amount         int64      `gorm:"type:bigint;not null;column:amount" json:"amount"`
	Currency       string     `gorm:"type:varchar(3);not null;default:'KZT';column:currency" json:"currency"`
	TargetAmount   *int64     `gorm:"type:bigint;column:target_amount" json:"target_amount"`
	TargetCurrency *string    `gorm:"type:varchar(3);column:target_currency" json:"target_currency"`
	ExchangeRate   *float64   `gorm:"type:numeric;column:exchange_rate" json:"exchange_rate"`
	Status         string     `gorm:"type:varchar(20);not null;default:'PENDING';column:status" json:"status"`
	IdempotencyKey string     `gorm:"type:varchar(255);unique;column:idempotency_key" json:"idempotency_key,omitempty"`
	CreatedAt      time.Time  `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt      *time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}

func (t Transaction) ToDTO() models.Transaction {
	return models.Transaction{
		ID:             t.ID,
		Type:           t.Type,
		FromAccountID:  t.FromAccountID,
		ToAccountID:    t.ToAccountID,
		Amount:         t.Amount,
		Currency:       t.Currency,
		TargetAmount:   t.TargetAmount,
		TargetCurrency: t.TargetCurrency,
		ExchangeRate:   t.ExchangeRate,
		Status:         t.Status,
		IdempotencyKey: t.IdempotencyKey,
		CreatedAt:      t.CreatedAt,
		UpdatedAt:      t.UpdatedAt,
	}
}

func FromTransactionDTO(t models.Transaction) Transaction {
	return Transaction{
		ID:             t.ID,
		Type:           t.Type,
		FromAccountID:  t.FromAccountID,
		ToAccountID:    t.ToAccountID,
		Amount:         t.Amount,
		Currency:       t.Currency,
		TargetAmount:   t.TargetAmount,
		TargetCurrency: t.TargetCurrency,
		ExchangeRate:   t.ExchangeRate,
		Status:         t.Status,
		IdempotencyKey: t.IdempotencyKey,
		CreatedAt:      t.CreatedAt,
		UpdatedAt:      t.UpdatedAt,
	}
}
