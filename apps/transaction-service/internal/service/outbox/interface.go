package outbox

import (
	"context"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/pkg/models"
)

type OutboxRepo interface {
	Create(ctx context.Context, event models.Outbox) (*models.Outbox, error)
	GetAll(ctx context.Context) ([]models.Outbox, error)
	Delete(ctx context.Context, id uuid.UUID) error
}


