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

type TransactionRepo struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewTransactionRepo(db *gorm.DB, logger *zap.SugaredLogger) *TransactionRepo {
	return &TransactionRepo{
		db:     db,
		logger: logger.Named("transaction_repo"),
	}
}

func (r *TransactionRepo) tx(ctx context.Context) *gorm.DB {
	if tx := RetrieveTx(ctx); tx != nil {
		return tx
	}
	return r.db.WithContext(ctx)
}

func (r *TransactionRepo) Create(ctx context.Context, transaction models.Transaction) (*models.Transaction, error) {
	e := entity.FromTransactionDTO(transaction)
	if err := r.tx(ctx).Create(&e).Error; err != nil {
		r.logger.Errorw("failed to create transaction", "error", err, "id", transaction.ID)
		return nil, err
	}
	dto := e.ToDTO()
	return &dto, nil
}

func (r *TransactionRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	var e entity.Transaction
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

func (r *TransactionRepo) Update(ctx context.Context, transaction models.Transaction) (*models.Transaction, error) {
	e := entity.FromTransactionDTO(transaction)
	if err := r.tx(ctx).Save(&e).Error; err != nil {
		r.logger.Errorw("failed to update transaction", "id", transaction.ID, "error", err)
		return nil, err
	}
	dto := e.ToDTO()
	return &dto, nil
}

func (r *TransactionRepo) GetByIdempotencyKey(ctx context.Context, key string) (*models.Transaction, error) {
	var e entity.Transaction
	err := r.tx(ctx).Where("idempotency_key = ?", key).First(&e).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	dto := e.ToDTO()
	return &dto, nil
}
