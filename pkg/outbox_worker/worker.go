package outbox_worker

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/w0ikid/yarmaq/pkg/models/entity"
)

type Publisher interface {
	Publish(ctx context.Context, eventType string, payload []byte) error
}

type Worker struct {
	db        *gorm.DB
	publisher Publisher
	interval  time.Duration
	batchSize int
	logger    *zap.SugaredLogger
}

func NewWorker(db *gorm.DB, publisher Publisher, logger *zap.SugaredLogger, opts ...Option) *Worker {
	w := &Worker{
		db:        db,
		publisher: publisher,
		interval:  5 * time.Second,
		batchSize: 100,
		logger:    logger.With("component", "outbox-worker"),
	}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

type Option func(*Worker)

func WithInterval(d time.Duration) Option {
	return func(w *Worker) { w.interval = d }
}

func WithBatchSize(n int) Option {
	return func(w *Worker) { w.batchSize = n }
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	w.logger.Info("outbox worker started", "interval", w.interval)

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("outbox worker stopped")
			return
		case <-ticker.C:
			if err := w.processOnce(ctx); err != nil {
				w.logger.Error("outbox processing failed", "err", err)
			}
		}
	}
}

func (w *Worker) processOnce(ctx context.Context) error {
	return w.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var events []entity.Outbox

		err := tx.
			Clauses(clause.Locking{
				Strength: "UPDATE",
				Options:  "SKIP LOCKED",
			}).
			Where("sent_at IS NULL").
			Order("created_at").
			Limit(w.batchSize).
			Find(&events).Error
		if err != nil {
			return err
		}

		if len(events) == 0 {
			return nil
		}

		for _, e := range events {
			if err := w.publisher.Publish(ctx, e.EventType, e.Payload); err != nil {
				return err
			}

			if err := tx.Model(&e).Update("sent_at", time.Now()).Error; err != nil {
				return err
			}
		}

		w.logger.Info("outbox batch processed", "count", len(events))
		return nil
	})
}
