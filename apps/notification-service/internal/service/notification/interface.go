package notification

import (
	"context"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/pkg/models"
)

type NotificationRepo interface {
	Create(ctx context.Context, notification models.Notification) (*models.Notification, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Notification, error)
	GetByUserID(ctx context.Context, userID string) ([]models.Notification, error)
	Update(ctx context.Context, notification models.Notification) (*models.Notification, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
