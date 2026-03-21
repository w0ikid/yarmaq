package ledger

import (
	"context"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/service/ledger"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/usecase"
	"github.com/w0ikid/yarmaq/pkg/models"
)

type GetLedgerUsecase struct {
	usecase.BaseUsecase
	LedgerService interface {
		GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]models.Ledger, error)
		GetAll(ctx context.Context) ([]models.Ledger, error)
	}
}

func NewGetLedgerUsecase(base usecase.BaseUsecase, ledgerService ledger.Service) GetLedgerUsecase {
	return GetLedgerUsecase{
		BaseUsecase:   base,
		LedgerService: ledgerService,
	}
}

func (uc *GetLedgerUsecase) ExecuteByAccountID(ctx context.Context, accountID uuid.UUID) ([]models.Ledger, error) {
	uc.Logger.Infow("starting GetLedgerUsecase.ExecuteByAccountID", "account_id", accountID)
	return uc.LedgerService.GetByAccountID(ctx, accountID)
}

func (uc *GetLedgerUsecase) ExecuteAll(ctx context.Context) ([]models.Ledger, error) {
	uc.Logger.Infow("starting GetLedgerUsecase.ExecuteAll")
	return uc.LedgerService.GetAll(ctx)
}
