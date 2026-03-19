package outbox

import (
	"github.com/w0ikid/yarmaq/apps/transaction-service/internal/service/outbox"
	"github.com/w0ikid/yarmaq/apps/transaction-service/internal/usecase"
)


type OutboxDomain struct {
}

func NewDomain(baseusecase usecase.BaseUsecase, outboxService outbox.Service) OutboxDomain {
	return OutboxDomain{}
}
