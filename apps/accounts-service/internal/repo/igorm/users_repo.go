package igorm

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/w0ikid/yarmaq/pkg/models"
	"github.com/w0ikid/yarmaq/pkg/models/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UsersRepo struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewUsersRepo(db *gorm.DB, logger *zap.SugaredLogger) *UsersRepo {
	return &UsersRepo{
		db:     db,
		logger: logger.Named("users_repo"),
	}
}

func (r *UsersRepo) tx(ctx context.Context) *gorm.DB {
	if tx := RetrieveTx(ctx); tx != nil {
		return tx
	}
	return r.db.WithContext(ctx)
}

func (r *UsersRepo) Create(ctx context.Context, user models.User) (*models.User, error) {
	e := entity.FromDTO(user)

	r.logger.Infow("executing db create", "zitadel_id", e.ZitadelUserID, "email", e.Email)

	if err := r.tx(ctx).Create(&e).Error; err != nil {
		r.logger.Errorw("failed to create user", "error", err, "zitadel_id", e.ZitadelUserID)
		return nil, err
	}

	r.logger.Infow("db create success", "id", e.ID)
	dto := e.ToDTO()
	return &dto, nil
}

func (r *UsersRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var e entity.User
	err := r.tx(ctx).Where("id = ?", id).First(&e).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.Errorw("failed to get user by id", "id", id, "error", err)
		return nil, err
	}
	dto := e.ToDTO()
	return &dto, nil
}

func (r *UsersRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var e entity.User
	err := r.tx(ctx).Where("email = ?", email).First(&e).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.Errorw("failed to get user by email", "email", email, "error", err)
		return nil, err
	}
	dto := e.ToDTO()
	return &dto, nil
}

func (r *UsersRepo) GetByZitadelID(ctx context.Context, zitadelID string) (*models.User, error) {
	var e entity.User
	err := r.tx(ctx).Where("zitadel_user_id = ?", zitadelID).First(&e).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.Errorw("failed to get user by zitadel id", "zitadel_id", zitadelID, "error", err)
		return nil, err
	}
	dto := e.ToDTO()
	return &dto, nil
}

func (r *UsersRepo) GetAll(ctx context.Context) ([]models.User, error) {
	var entities []entity.User
	if err := r.tx(ctx).Find(&entities).Error; err != nil {
		r.logger.Errorw("failed to get all users", "error", err)
		return nil, err
	}
	users := make([]models.User, len(entities))
	for i, e := range entities {
		users[i] = e.ToDTO()
	}
	return users, nil
}

func (r *UsersRepo) Update(ctx context.Context, user models.User) (*models.User, error) {
	var e entity.User
	err := r.tx(ctx).Where("id = ?", user.ID).First(&e).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	e.Email = user.Email
	e.Username = user.Username
	e.Roles = pq.StringArray(user.Roles)
	e.ImageURL = user.ImageURL
	e.IsActive = user.IsActive
	e.Address = user.Address
	if err := r.tx(ctx).Save(&e).Error; err != nil {
		r.logger.Errorw("failed to update user", "id", user.ID, "error", err)
		return nil, err
	}
	dto := e.ToDTO()
	return &dto, nil
}

func (r *UsersRepo) UpdateRoles(ctx context.Context, zitadelID string, roles []string) (*models.User, error) {
	var e entity.User
	err := r.tx(ctx).Where("zitadel_user_id = ?", zitadelID).First(&e).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.Errorw("failed to find user for role update", "zitadel_id", zitadelID, "error", err)
		return nil, err
	}
	e.Roles = pq.StringArray(roles)
	if err := r.tx(ctx).Save(&e).Error; err != nil {
		r.logger.Errorw("failed to update roles", "zitadel_id", zitadelID, "error", err)
		return nil, err
	}
	dto := e.ToDTO()
	return &dto, nil
}

func (r *UsersRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.tx(ctx).Where("id = ?", id).Delete(&entity.User{}).Error; err != nil {
		r.logger.Errorw("failed to delete user", "id", id, "error", err)
		return err
	}
	return nil
}
