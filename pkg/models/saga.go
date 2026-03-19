package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	SagaStepNameHold     = "HOLD"
	SagaStepNameDeposit  = "DEPOSIT"
	SagaStepNameWithdraw = "WITHDRAW"
	SagaStepNameRefund   = "REFUND"

	SagaStatusPending   = "PENDING"
	SagaStatusCompleted = "COMPLETED"
	SagaStatusFailed    = "FAILED"
)

type SagaStep struct {
	ID            uuid.UUID `json:"id"`
	TransactionID uuid.UUID `json:"transaction_id"`
	StepName      string    `json:"step_name"`
	Status        string    `json:"status"`
	Error         *string   `json:"error,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}
