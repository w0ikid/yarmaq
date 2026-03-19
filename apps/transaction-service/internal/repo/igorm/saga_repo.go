package igorm

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/pkg/models"
	"github.com/w0ikid/yarmaq/pkg/models/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SagaStepRepo struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewSagaStepRepo(db *gorm.DB, logger *zap.SugaredLogger) *SagaStepRepo {
	return &SagaStepRepo{
		db:     db,
		logger: logger.Named("saga_step_repo"),
	}
}

func (r *SagaStepRepo) tx(ctx context.Context) *gorm.DB {
	if tx := RetrieveTx(ctx); tx != nil {
		return tx
	}
	return r.db.WithContext(ctx)
}

func (r *SagaStepRepo) Create(ctx context.Context, step models.SagaStep) (*models.SagaStep, error) {
	e := entity.FromSagaStepDTO(step)
	if err := r.tx(ctx).Create(&e).Error; err != nil {
		r.logger.Errorw("failed to create saga step", "error", err, "transaction_id", step.TransactionID)
		return nil, err
	}
	dto := e.ToDTO()
	return &dto, nil
}

func (r *SagaStepRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.SagaStep, error) {
	var e entity.SagaStep
	err := r.tx(ctx).Where("id = ?", id).First(&e).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	dto := e.ToDTO()
	return &dto, nil
}

func (r *SagaStepRepo) GetByTransactionID(ctx context.Context, transactionID uuid.UUID) ([]models.SagaStep, error) {
	var entities []entity.SagaStep
	err := r.tx(ctx).Where("transaction_id = ?", transactionID).Order("created_at asc").Find(&entities).Error
	if err != nil {
		return nil, err
	}
	
	dtos := make([]models.SagaStep, len(entities))
	for i, e := range entities {
		dtos[i] = e.ToDTO()
	}
	return dtos, nil
}

func (r *SagaStepRepo) Update(ctx context.Context, step models.SagaStep) (*models.SagaStep, error) {
	e := entity.FromSagaStepDTO(step)
	if err := r.tx(ctx).Save(&e).Error; err != nil {
		r.logger.Errorw("failed to update saga step", "id", step.ID, "error", err)
		return nil, err
	}
	dto := e.ToDTO()
	return &dto, nil
}
