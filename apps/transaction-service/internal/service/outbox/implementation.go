package outbox

import (
	"context"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/pkg/models"
	"go.uber.org/zap"
)

type Service interface {
	Create(ctx context.Context, event models.Outbox) (*models.Outbox, error)
	GetAll(ctx context.Context) ([]models.Outbox, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type implementation struct {
	repo   OutboxRepo
	logger *zap.SugaredLogger
}

func NewService(repo OutboxRepo, logger *zap.SugaredLogger) Service {
	return &implementation{
		repo:   repo,
		logger: logger.Named("outbox_service"),
	}
}

func (s *implementation) Create(ctx context.Context, event models.Outbox) (*models.Outbox, error) {
	return s.repo.Create(ctx, event)
}

func (s *implementation) GetAll(ctx context.Context) ([]models.Outbox, error) {
	return s.repo.GetAll(ctx)
}

func (s *implementation) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
