package saga

import (
	"context"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/pkg/models"
)

type SagaStepRepo interface {
	Create(ctx context.Context, step models.SagaStep) (*models.SagaStep, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.SagaStep, error)
	GetByTransactionID(ctx context.Context, transactionID uuid.UUID) ([]models.SagaStep, error)
	Update(ctx context.Context, step models.SagaStep) (*models.SagaStep, error)
}
