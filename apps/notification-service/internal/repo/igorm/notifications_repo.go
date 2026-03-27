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

type NotificationRepo struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewNotificationRepo(db *gorm.DB, logger *zap.SugaredLogger) *NotificationRepo {
	return &NotificationRepo{
		db:     db,
		logger: logger.Named("notification_repo"),
	}
}

func (r *NotificationRepo) tx(ctx context.Context) *gorm.DB {
	if tx := RetrieveTx(ctx); tx != nil {
		return tx
	}
	return r.db.WithContext(ctx)
}

func (r *NotificationRepo) Create(ctx context.Context, notification models.Notification) (*models.Notification, error) {
	e := entity.FromNotificationDTO(notification)
	if err := r.tx(ctx).Create(&e).Error; err != nil {
		r.logger.Errorw("failed to create notification", "error", err, "user_id", notification.UserID, "type", notification.Type)
		return nil, err
	}
	dto := e.ToDTO()
	return &dto, nil
}

func (r *NotificationRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Notification, error) {
	var e entity.Notification
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

func (r *NotificationRepo) GetByUserID(ctx context.Context, userID string) ([]models.Notification, error) {
	var entities []entity.Notification
	if err := r.tx(ctx).Where("user_id = ?", userID).Order("created_at desc").Find(&entities).Error; err != nil {
		return nil, err
	}

	notifications := make([]models.Notification, len(entities))
	for i, e := range entities {
		notifications[i] = e.ToDTO()
	}

	return notifications, nil
}

func (r *NotificationRepo) Update(ctx context.Context, notification models.Notification) (*models.Notification, error) {
	e := entity.FromNotificationDTO(notification)
	if err := r.tx(ctx).Save(&e).Error; err != nil {
		r.logger.Errorw("failed to update notification", "id", notification.ID, "error", err)
		return nil, err
	}
	dto := e.ToDTO()
	return &dto, nil
}

func (r *NotificationRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.tx(ctx).Where("id = ?", id).Delete(&entity.Notification{}).Error; err != nil {
		r.logger.Errorw("failed to delete notification", "id", id, "error", err)
		return err
	}
	return nil
}
