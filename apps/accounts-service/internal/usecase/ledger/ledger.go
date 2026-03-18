package ledger

import (
	"context"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/usecase"
	"github.com/w0ikid/yarmaq/pkg/models"
)

type Service interface {
	GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]models.Ledger, error)
	GetAll(ctx context.Context) ([]models.Ledger, error)
}

type LedgerDomain struct {
	GetLedgerUsecase GetLedgerUsecase
}

func NewDomain(baseusecase usecase.BaseUsecase, ledgerService Service) LedgerDomain {
	baseusecase.Logger = baseusecase.Logger.Named("ledger_domain")
	return LedgerDomain{
		GetLedgerUsecase: GetLedgerUsecase{
			BaseUsecase:   baseusecase,
			LedgerService: ledgerService,
		},
	}
}
