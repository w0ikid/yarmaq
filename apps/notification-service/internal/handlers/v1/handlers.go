package v1

import (
	"github.com/w0ikid/yarmaq/pkg/jwks"
	"go.uber.org/zap"
)

type Dependencies struct {
	Logger *zap.SugaredLogger

	JWKS            *jwks.JWKS
}

type Handlers struct {
	JWKS        *jwks.JWKS
}

func NewHandlers(deps Dependencies) *Handlers {
	return &Handlers{
		JWKS:        deps.JWKS,
	}
}
