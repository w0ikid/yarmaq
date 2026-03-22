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

type AccountRepo struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewAccountRepo(db *gorm.DB, logger *zap.SugaredLogger) *AccountRepo {
	return &AccountRepo{
		db:     db,
		logger: logger.Named("account_repo"),
	}
}

func (r *AccountRepo) tx(ctx context.Context) *gorm.DB {
	if tx := RetrieveTx(ctx); tx != nil {
		return tx
	}
	return r.db.WithContext(ctx)
}

func (r *AccountRepo) Create(ctx context.Context, account models.Account) (*models.Account, error) {
	e := entity.FromAccountDTO(account)
	if err := r.tx(ctx).Create(&e).Error; err != nil {
		r.logger.Errorw("failed to create account", "error", err, "user_id", account.UserID)
		return nil, err
	}
	dto := e.ToDTO()
	return &dto, nil
}

func (r *AccountRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Account, error) {
	var e entity.Account
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

func (r *AccountRepo) GetByNumber(ctx context.Context, number string) (*models.Account, error) {
	var e entity.Account
	err := r.tx(ctx).Where("number = ?", number).First(&e).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	dto := e.ToDTO()
	return &dto, nil
}

func (r *AccountRepo) GetByNumberAndCurrency(ctx context.Context, number string, currency string) (*models.Account, error) {
	var e entity.Account
	err := r.tx(ctx).Where("number = ? AND currency = ?", number, currency).First(&e).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	dto := e.ToDTO()
	return &dto, nil
}

func (r *AccountRepo) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Account, error) {
	var e entity.Account
	err := r.tx(ctx).Where("user_id = ?", userID).First(&e).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	dto := e.ToDTO()
	return &dto, nil
}

func (r *AccountRepo) GetByUserIDAndCurrency(ctx context.Context, userID string, currency string) (*models.Account, error) {
	var e entity.Account
	err := r.tx(ctx).Where("user_id = ? AND currency = ?", userID, currency).First(&e).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	dto := e.ToDTO()
	return &dto, nil
}

func (r *AccountRepo) GetByTypeAndCurrency(ctx context.Context, accountType string, currency string) (*models.Account, error) {
	var e entity.Account
	err := r.tx(ctx).Where("type = ? AND currency = ?", accountType, currency).First(&e).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	dto := e.ToDTO()
	return &dto, nil
}

func (r *AccountRepo) GetAllByUserID(ctx context.Context, userID string) ([]models.Account, error) {
	var entities []entity.Account
	if err := r.tx(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&entities).Error; err != nil {
		return nil, err
	}

	accounts := make([]models.Account, 0, len(entities))
	for _, e := range entities {
		accounts = append(accounts, e.ToDTO())
	}

	return accounts, nil
}

func (r *AccountRepo) Update(ctx context.Context, account models.Account) (*models.Account, error) {
	e := entity.FromAccountDTO(account)
	if err := r.tx(ctx).Save(&e).Error; err != nil {
		r.logger.Errorw("failed to update account", "id", account.ID, "error", err)
		return nil, err
	}
	dto := e.ToDTO()
	return &dto, nil
}

func (r *AccountRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.tx(ctx).Where("id = ?", id).Delete(&entity.Account{}).Error; err != nil {
		r.logger.Errorw("failed to delete account", "id", id, "error", err)
		return err
	}
	return nil
}

func (r *AccountRepo) NextSeq(ctx context.Context) (int64, error) {
	var seq int64
	err := r.tx(ctx).Raw("SELECT nextval('account_number_seq')").Scan(&seq).Error
	if err != nil {
		r.logger.Errorw("failed to get next seq", "error", err)
		return 0, err
	}
	return seq, nil
}
