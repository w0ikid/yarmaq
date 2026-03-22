package v1

import (
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/handlers/v1/account"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/handlers/v1/internals"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/handlers/v1/ledger"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/handlers/v1/webhook"
	"github.com/w0ikid/yarmaq/pkg/jwks"
	"go.uber.org/zap"
)

type Dependencies struct {
	Logger *zap.SugaredLogger

	AccountDeps  account.HandlerDeps
	InternalDeps internals.HandlerDeps
	LedgerDeps   ledger.HandlerDeps
	WebhookDeps webhook.HandlerDeps
	JWKS        *jwks.JWKS
}

type Handlers struct {
	Account  account.Handler
	Internal internals.Handler
	Ledger   ledger.Handler
	Webhook webhook.Handler
	JWKS    *jwks.JWKS
}

func NewHandlers(deps Dependencies) *Handlers {
	return &Handlers{
		Account:  account.NewHandler(deps.AccountDeps),
		Internal: internals.NewHandler(deps.InternalDeps),
		Ledger:   ledger.NewHandler(deps.LedgerDeps),
		Webhook:  webhook.NewHandler(deps.WebhookDeps),
		JWKS:     deps.JWKS,
	}
}
