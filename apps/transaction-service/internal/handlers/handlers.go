package handlers

import (
	v1 "github.com/w0ikid/yarmaq/apps/transaction-service/internal/handlers/v1"
	"github.com/w0ikid/yarmaq/apps/transaction-service/internal/handlers/v1/transaction"
	"github.com/w0ikid/yarmaq/pkg/jwks"
)

type Depedencies struct {
	TransactionDeps transaction.HandlerDeps
	JWKS            *jwks.JWKS
}

type Handlers struct {
	V1   *v1.Handlers
	JWKS *jwks.JWKS
}

func NewHandlers(deps Depedencies) *Handlers {
	return &Handlers{
		V1: v1.NewHandlers(v1.Dependencies{
			TransactionDeps: deps.TransactionDeps,
			JWKS:            deps.JWKS,
		}),
		JWKS: deps.JWKS,
	}
}
