package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/pkg/models"
)

type Ledger struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id" json:"id"`
	AccountID     uuid.UUID  `gorm:"type:uuid;not null;column:account_id" json:"account_id"`
	Amount        int64      `gorm:"type:bigint;not null;column:amount" json:"amount"`
	OperationType string     `gorm:"type:varchar(50);not null;column:operation_type" json:"operation_type"`
	ReferenceID   *uuid.UUID `gorm:"type:uuid;column:reference_id" json:"reference_id,omitempty"`
	CreatedAt     time.Time  `gorm:"autoCreateTime;column:created_at" json:"created_at"`
}

func (l Ledger) ToDTO() models.Ledger {
	return models.Ledger{
		ID:            l.ID,
		AccountID:     l.AccountID,
		Amount:        l.Amount,
		OperationType: l.OperationType,
		ReferenceID:   l.ReferenceID,
		CreatedAt:     l.CreatedAt,
	}
}

func FromLedgerDTO(l models.Ledger) Ledger {
	return Ledger{
		ID:            l.ID,
		AccountID:     l.AccountID,
		Amount:        l.Amount,
		OperationType: l.OperationType,
		ReferenceID:   l.ReferenceID,
		CreatedAt:     l.CreatedAt,
	}
}
