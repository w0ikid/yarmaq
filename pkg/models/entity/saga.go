package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/pkg/models"
)

type SagaStep struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id" json:"id"`
	TransactionID  uuid.UUID `gorm:"type:uuid;not null;column:transaction_id" json:"transaction_id"`
	StepName       string    `gorm:"type:varchar(50);not null;column:step_name" json:"step_name"`
	Status         string    `gorm:"type:varchar(20);not null;default:'PENDING';column:status" json:"status"`
	Error          *string   `gorm:"type:text;column:error" json:"error,omitempty"`
	CreatedAt      time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
}

func (s SagaStep) ToDTO() models.SagaStep {
	return models.SagaStep{
		ID:            s.ID,
		TransactionID: s.TransactionID,
		StepName:      s.StepName,
		Status:        s.Status,
		Error:         s.Error,
		CreatedAt:     s.CreatedAt,
	}
}

func FromSagaStepDTO(s models.SagaStep) SagaStep {
	return SagaStep{
		ID:            s.ID,
		TransactionID: s.TransactionID,
		StepName:      s.StepName,
		Status:        s.Status,
		Error:         s.Error,
		CreatedAt:     s.CreatedAt,
	}
}
