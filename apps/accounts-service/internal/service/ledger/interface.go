package ledger

import (
	"context"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/pkg/models"
)

type LedgerRepo interface {
	GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]models.Ledger, error)
	GetAll(ctx context.Context) ([]models.Ledger, error)
}


