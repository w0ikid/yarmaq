package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/pkg/models"
)

type Transaction struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id" json:"id"`
	FromAccountID  uuid.UUID  `gorm:"type:uuid;not null;column:from_account_id" json:"from_account_id"`
	ToAccountID    uuid.UUID  `gorm:"type:uuid;not null;column:to_account_id" json:"to_account_id"`
	Amount         int64      `gorm:"type:bigint;not null;column:amount" json:"amount"`
	Currency       string     `gorm:"type:varchar(3);not null;default:'KZT';column:currency" json:"currency"`
	Status         string     `gorm:"type:varchar(20);not null;default:'PENDING';column:status" json:"status"`
	IdempotencyKey string     `gorm:"type:varchar(255);unique;column:idempotency_key" json:"idempotency_key,omitempty"`
	CreatedAt      time.Time  `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt      *time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}

func (t Transaction) ToDTO() models.Transaction {
	return models.Transaction{
		ID:             t.ID,
		FromAccountID:  t.FromAccountID,
		ToAccountID:    t.ToAccountID,
		Amount:         t.Amount,
		Currency:       t.Currency,
		Status:         t.Status,
		IdempotencyKey: t.IdempotencyKey,
		CreatedAt:      t.CreatedAt,
		UpdatedAt:      t.UpdatedAt,
	}
}

func FromTransactionDTO(t models.Transaction) Transaction {
	return Transaction{
		ID:             t.ID,
		FromAccountID:  t.FromAccountID,
		ToAccountID:    t.ToAccountID,
		Amount:         t.Amount,
		Currency:       t.Currency,
		Status:         t.Status,
		IdempotencyKey: t.IdempotencyKey,
		CreatedAt:      t.CreatedAt,
		UpdatedAt:      t.UpdatedAt,
	}
}
