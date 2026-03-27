package handlers

import (
	v1 "github.com/w0ikid/yarmaq/apps/notification-service/internal/handlers/v1"
	"github.com/w0ikid/yarmaq/pkg/jwks"
)

type Depedencies struct {
	JWKS            *jwks.JWKS
}

type Handlers struct {
	V1   *v1.Handlers
	JWKS *jwks.JWKS
}

func NewHandlers(deps Depedencies) *Handlers {
	return &Handlers{
		V1: v1.NewHandlers(v1.Dependencies{
			JWKS:            deps.JWKS,
		}),
		JWKS: deps.JWKS,
	}
}
