package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/pkg/models"
)

type IContextTransaction interface {
	StartTransaction(ctx context.Context) (context.Context, error)
	FinalizeTransaction(ctx context.Context, err *error) error
}

type IOutboxRepo interface {
	Create(ctx context.Context, event models.Outbox) (*models.Outbox, error)
	GetAll(ctx context.Context) ([]models.Outbox, error)
	GetUnsent(ctx context.Context) ([]models.Outbox, error)
	Update(ctx context.Context, event models.Outbox) (*models.Outbox, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type ITransactionRepo interface {
	Create(ctx context.Context, transaction models.Transaction) (*models.Transaction, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Transaction, error)
	Update(ctx context.Context, transaction models.Transaction) (*models.Transaction, error)
	GetByIdempotencyKey(ctx context.Context, key string) (*models.Transaction, error)
}

type ISagaStepRepo interface {
	Create(ctx context.Context, step models.SagaStep) (*models.SagaStep, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.SagaStep, error)
	GetByTransactionID(ctx context.Context, transactionID uuid.UUID) ([]models.SagaStep, error)
	Update(ctx context.Context, step models.SagaStep) (*models.SagaStep, error)
}

type Repository struct {
	ContextTransaction IContextTransaction
	Outbox             IOutboxRepo
	Transaction        ITransactionRepo
	SagaStep           ISagaStepRepo
}
