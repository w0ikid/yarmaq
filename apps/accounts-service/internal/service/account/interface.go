package account

import (
	"context"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/pkg/models"
)

type AccountRepo interface {
	Create(ctx context.Context, account models.Account) (*models.Account, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Account, error)
	GetByNumber(ctx context.Context, number string) (*models.Account, error)
	GetByNumberAndCurrency(ctx context.Context, number string, currency string) (*models.Account, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Account, error)
	GetByUserIDAndCurrency(ctx context.Context, userID string, currency string) (*models.Account, error)
	GetByTypeAndCurrency(ctx context.Context, accountType string, currency string) (*models.Account, error)
	GetAllByUserID(ctx context.Context, userID string) ([]models.Account, error)
	Update(ctx context.Context, account models.Account) (*models.Account, error)
	Delete(ctx context.Context, id uuid.UUID) error

	NextSeq(ctx context.Context) (int64, error)
}

type LedgerRepo interface {
	Create(ctx context.Context, entry models.Ledger) (*models.Ledger, error)
}

type OutboxRepo interface {
	Create(ctx context.Context, event models.Outbox) (*models.Outbox, error)
}
