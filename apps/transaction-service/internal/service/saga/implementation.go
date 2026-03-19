package saga

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/pkg/models"
	"go.uber.org/zap"
)

type Service interface {
	CreateStep(ctx context.Context, step models.SagaStep) (*models.SagaStep, error)
	GetSteps(ctx context.Context, transactionID uuid.UUID) ([]models.SagaStep, error)
	UpdateStepStatus(ctx context.Context, stepID uuid.UUID, status string, errStr *string) error
}

type implementation struct {
	repo   SagaStepRepo
	logger *zap.SugaredLogger
}

func NewService(repo SagaStepRepo, logger *zap.SugaredLogger) Service {
	return &implementation{
		repo:   repo,
		logger: logger.Named("saga_service"),
	}
}

func (s *implementation) CreateStep(ctx context.Context, step models.SagaStep) (*models.SagaStep, error) {
	s.logger.Infow("creating saga step", "transaction_id", step.TransactionID, "step", step.StepName)
	return s.repo.Create(ctx, step)
}

func (s *implementation) GetSteps(ctx context.Context, transactionID uuid.UUID) ([]models.SagaStep, error) {
	return s.repo.GetByTransactionID(ctx, transactionID)
}

func (s *implementation) UpdateStepStatus(ctx context.Context, stepID uuid.UUID, status string, errStr *string) error {
	s.logger.Infow("updating saga step status", "id", stepID, "status", status)
	
    step, err := s.repo.GetByID(ctx, stepID)
    if err != nil {
        return err
    }
    if step == nil {
        return fmt.Errorf("saga step not found: %s", stepID)
    }

    step.Status = status
    step.Error = errStr
    
    _, err = s.repo.Update(ctx, *step)
    return err
}
