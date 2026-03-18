package ledger

import (
	"context"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/pkg/models"
	"go.uber.org/zap"
)

type Service interface {
	GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]models.Ledger, error)
	GetAll(ctx context.Context) ([]models.Ledger, error)
}

type implementation struct {
	repo   LedgerRepo
	logger *zap.SugaredLogger
}

func NewService(repo LedgerRepo, logger *zap.SugaredLogger) Service {
	return &implementation{
		repo:   repo,
		logger: logger.Named("ledger_service"),
	}
}

func (s *implementation) GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]models.Ledger, error) {
	return s.repo.GetByAccountID(ctx, accountID)
}

func (s *implementation) GetAll(ctx context.Context) ([]models.Ledger, error) {
	return s.repo.GetAll(ctx)
}
