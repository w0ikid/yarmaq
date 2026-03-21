package igorm

import (
	"github.com/w0ikid/yarmaq/apps/transaction-service/internal/repo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func NewGormRepository(db *gorm.DB, logger *zap.SugaredLogger) *repo.Repository {
	log := logger.Named("gorm_repository")
	log.Info("Initializing GORM repository")
	return &repo.Repository{
		ContextTransaction: NewContextTransaction(db),
		Outbox:             NewOutboxRepo(db, log),
		Transaction:        NewTransactionRepo(db, log),
		SagaStep:           NewSagaStepRepo(db, log),
	}
}
