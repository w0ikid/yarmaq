package igorm

import (
	"context"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/pkg/models"
	"github.com/w0ikid/yarmaq/pkg/models/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type LedgerRepo struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewLedgerRepo(db *gorm.DB, logger *zap.SugaredLogger) *LedgerRepo {
	return &LedgerRepo{
		db:     db,
		logger: logger.Named("ledger_repo"),
	}
}

func (r *LedgerRepo) tx(ctx context.Context) *gorm.DB {
	if tx := RetrieveTx(ctx); tx != nil {
		return tx
	}
	return r.db.WithContext(ctx)
}

func (r *LedgerRepo) Create(ctx context.Context, entry models.Ledger) (*models.Ledger, error) {
	e := entity.FromLedgerDTO(entry)
	if err := r.tx(ctx).Create(&e).Error; err != nil {
		r.logger.Errorw("failed to create ledger entry", "error", err, "account_id", entry.AccountID)
		return nil, err
	}
	dto := e.ToDTO()
	return &dto, nil
}

func (r *LedgerRepo) GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]models.Ledger, error) {
	var entities []entity.Ledger
	err := r.tx(ctx).Where("account_id = ?", accountID).Order("created_at DESC").Find(&entities).Error
	if err != nil {
		return nil, err
	}
	dtos := make([]models.Ledger, len(entities))
	for i, e := range entities {
		dtos[i] = e.ToDTO()
	}
	return dtos, nil
}

func (r *LedgerRepo) GetAll(ctx context.Context) ([]models.Ledger, error) {
	var entities []entity.Ledger
	if err := r.tx(ctx).Find(&entities).Error; err != nil {
		return nil, err
	}
	dtos := make([]models.Ledger, len(entities))
	for i, e := range entities {
		dtos[i] = e.ToDTO()
	}
	return dtos, nil
}
