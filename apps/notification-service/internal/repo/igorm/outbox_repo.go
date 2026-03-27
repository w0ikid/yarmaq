package igorm

import (
	"context"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/pkg/models"
	"github.com/w0ikid/yarmaq/pkg/models/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OutboxRepo struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewOutboxRepo(db *gorm.DB, logger *zap.SugaredLogger) *OutboxRepo {
	return &OutboxRepo{
		db:     db,
		logger: logger.Named("outbox_repo"),
	}
}

func (r *OutboxRepo) tx(ctx context.Context) *gorm.DB {
	if tx := RetrieveTx(ctx); tx != nil {
		return tx
	}
	return r.db.WithContext(ctx)
}

func (r *OutboxRepo) Create(ctx context.Context, event models.Outbox) (*models.Outbox, error) {
	e := entity.FromOutboxDTO(event)
	if err := r.tx(ctx).Create(&e).Error; err != nil {
		r.logger.Errorw("failed to create outbox event", "error", err, "event_type", event.EventType)
		return nil, err
	}
	dto := e.ToDTO()
	return &dto, nil
}

func (r *OutboxRepo) GetAll(ctx context.Context) ([]models.Outbox, error) {
	var entities []entity.Outbox
	if err := r.tx(ctx).Find(&entities).Error; err != nil {
		return nil, err
	}
	dtos := make([]models.Outbox, len(entities))
	for i, e := range entities {
		dtos[i] = e.ToDTO()
	}
	return dtos, nil
}

func (r *OutboxRepo) GetUnsent(ctx context.Context) ([]models.Outbox, error) {
	var entities []entity.Outbox
	if err := r.tx(ctx).Where("sent_at IS NULL").Find(&entities).Error; err != nil {
		return nil, err
	}
	dtos := make([]models.Outbox, len(entities))
	for i, e := range entities {
		dtos[i] = e.ToDTO()
	}
	return dtos, nil
}

func (r *OutboxRepo) Update(ctx context.Context, event models.Outbox) (*models.Outbox, error) {
	e := entity.FromOutboxDTO(event)
	if err := r.tx(ctx).Save(&e).Error; err != nil {
		r.logger.Errorw("failed to update outbox event", "id", event.ID, "error", err)
		return nil, err
	}
	dto := e.ToDTO()
	return &dto, nil
}

func (r *OutboxRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.tx(ctx).Where("id = ?", id).Delete(&entity.Outbox{}).Error; err != nil {
		r.logger.Errorw("failed to delete outbox event", "id", id, "error", err)
		return err
	}
	return nil
}
