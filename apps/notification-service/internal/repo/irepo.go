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

type INotificationRepo interface {
	Create(ctx context.Context, notification models.Notification) (*models.Notification, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Notification, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.Notification, error)
	Update(ctx context.Context, notification models.Notification) (*models.Notification, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type Repository struct {
	ContextTransaction IContextTransaction
	Outbox             IOutboxRepo
	Notification       INotificationRepo
}
