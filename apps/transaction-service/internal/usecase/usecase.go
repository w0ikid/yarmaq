package usecase

import (
	"github.com/w0ikid/yarmaq/apps/transaction-service/internal/repo"
	"go.uber.org/zap"
)

type BaseUsecase struct {
	Logger *zap.SugaredLogger
	Tx     repo.IContextTransaction
}

func NewBaseUsecase(tx repo.IContextTransaction, logger *zap.SugaredLogger) BaseUsecase {
	return BaseUsecase{
		Logger: logger,
		Tx:     tx,
	}
}
