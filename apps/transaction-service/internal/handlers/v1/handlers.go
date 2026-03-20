package v1

import (
	"github.com/w0ikid/yarmaq/apps/transaction-service/internal/handlers/v1/account"
	"github.com/w0ikid/yarmaq/apps/transaction-service/internal/handlers/v1/ledger"
	"github.com/w0ikid/yarmaq/apps/transaction-service/internal/handlers/v1/transaction"
	"github.com/w0ikid/yarmaq/pkg/jwks"
	"go.uber.org/zap"
)

type Dependencies struct {
	Logger *zap.SugaredLogger

	AccountDeps     account.HandlerDeps
	LedgerDeps      ledger.HandlerDeps
	TransactionDeps transaction.HandlerDeps
	JWKS            *jwks.JWKS
}

type Handlers struct {
	Account     account.Handler
	Ledger      ledger.Handler
	Transaction transaction.Handler
	JWKS        *jwks.JWKS
}

func NewHandlers(deps Dependencies) *Handlers {
	return &Handlers{
		Account:     account.NewHandler(deps.AccountDeps),
		Ledger:      ledger.NewHandler(deps.LedgerDeps),
		Transaction: transaction.NewHandler(deps.TransactionDeps),
		JWKS:        deps.JWKS,
	}
}
