package v1

import (
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/handlers/v1/users"
	usersHandlers "github.com/w0ikid/yarmaq/apps/accounts-service/internal/handlers/v1/users"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/handlers/v1/webhook"
	"github.com/w0ikid/yarmaq/pkg/jwks"
	"go.uber.org/zap"
)

type Dependencies struct {
	Logger *zap.SugaredLogger

	UsersDeps   users.HandlerDeps
	WebhookDeps webhook.HandlerDeps
	JWKS        *jwks.JWKS
}

type Handlers struct {
	Users   usersHandlers.Handler
	Webhook webhook.Handler
	JWKS    *jwks.JWKS
}

func NewHandlers(deps Dependencies) *Handlers {
	return &Handlers{
		Users:   usersHandlers.NewHandler(deps.UsersDeps),
		Webhook: webhook.NewHandler(deps.WebhookDeps),
		JWKS:    deps.JWKS,
	}
}
