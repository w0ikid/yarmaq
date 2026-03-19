package handlers

import (
	v1 "github.com/w0ikid/yarmaq/apps/accounts-service/internal/handlers/v1"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/handlers/v1/account"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/handlers/v1/ledger"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/handlers/v1/users"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/handlers/v1/webhook"
	"github.com/w0ikid/yarmaq/pkg/jwks"
)

type Depedencies struct {
	UsersDeps   users.HandlerDeps
	AccountDeps account.HandlerDeps
	LedgerDeps  ledger.HandlerDeps
	WebhookDeps webhook.HandlerDeps
	JWKS        *jwks.JWKS
}

type Handlers struct {
	V1   *v1.Handlers
	JWKS *jwks.JWKS
}

func NewHandlers(deps Depedencies) *Handlers {
	return &Handlers{
		V1: v1.NewHandlers(v1.Dependencies{
			Logger:      deps.UsersDeps.Logger,
			UsersDeps:   deps.UsersDeps,
			AccountDeps: deps.AccountDeps,
			LedgerDeps:  deps.LedgerDeps,
			WebhookDeps: deps.WebhookDeps,
			JWKS:        deps.JWKS,
		}),
		JWKS: deps.JWKS,
	}
}
