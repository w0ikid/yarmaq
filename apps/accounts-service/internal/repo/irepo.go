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

type IAccountRepo interface {
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

type ILedgerRepo interface {
	Create(ctx context.Context, entry models.Ledger) (*models.Ledger, error)
	GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]models.Ledger, error)
	GetAll(ctx context.Context) ([]models.Ledger, error)
}

type IOutboxRepo interface {
	Create(ctx context.Context, event models.Outbox) (*models.Outbox, error)
	GetAll(ctx context.Context) ([]models.Outbox, error)
	GetUnsent(ctx context.Context) ([]models.Outbox, error)
	Update(ctx context.Context, event models.Outbox) (*models.Outbox, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type Repository struct {
	ContextTransaction IContextTransaction
	Account            IAccountRepo
	Ledger             ILedgerRepo
	Outbox             IOutboxRepo
}
