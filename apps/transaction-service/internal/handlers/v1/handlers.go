package v1

import (
	"github.com/w0ikid/yarmaq/apps/transaction-service/internal/handlers/v1/transaction"
	"github.com/w0ikid/yarmaq/pkg/jwks"
	"go.uber.org/zap"
)

type Dependencies struct {
	Logger *zap.SugaredLogger

	TransactionDeps transaction.HandlerDeps
	JWKS            *jwks.JWKS
}

type Handlers struct {
	Transaction transaction.Handler
	JWKS        *jwks.JWKS
}

func NewHandlers(deps Dependencies) *Handlers {
	return &Handlers{
		Transaction: transaction.NewHandler(deps.TransactionDeps),
		JWKS:        deps.JWKS,
	}
}
