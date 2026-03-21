package ledger

import (
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/service/ledger"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/usecase"
)

type LedgerDomain struct {
	GetLedgerUsecase GetLedgerUsecase
}

func NewDomain(baseusecase usecase.BaseUsecase, ledgerService ledger.Service) LedgerDomain {
	baseusecase.Logger = baseusecase.Logger.Named("ledger_domain")
	return LedgerDomain{
		GetLedgerUsecase: NewGetLedgerUsecase(baseusecase, ledgerService),
	}
}
