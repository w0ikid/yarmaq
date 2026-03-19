package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/w0ikid/yarmaq/pkg/config"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewPostgres(ctx context.Context, cfg config.PostgresConfig, logger *zap.SugaredLogger) (*Postgres, error) {
	sugar := logger.Named("postgres")
	
	sugar.Infof("connecting to postgres %s:%s/%s", cfg.Host, cfg.Port, cfg.DBName)
	dsn := fmt.Sprintf("%s connect_timeout=5", cfg.DSN())

	var db *gorm.DB
	var err error
	for i := 0; i < 3; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		sugar.Warnf("attempt %d: failed to open postgres connection: %v", i+1, err)
		time.Sleep(500 * time.Millisecond)
	}
	if err != nil {
		sugar.Error("all attempts to connect to postgres failed")
		return nil, fmt.Errorf("open postgres connection: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		sugar.Errorw("get sql db from gorm failed", "error", err)
		return nil, fmt.Errorf("get sql db from gorm: %w", err)
	}

	// Ping с retry
	for i := 0; i < 3; i++ {
		if err = sqlDB.PingContext(ctx); err == nil {
			break
		}
		sugar.Warnf("attempt %d: ping postgres failed: %v", i+1, err)
		time.Sleep(500 * time.Millisecond)
	}
	if err != nil {
		sugar.Error("all ping attempts to postgres failed")
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	sugar.Infof("connected to postgres %s:%s/%s", cfg.Host, cfg.Port, cfg.DBName)

	return &Postgres{
		db:     db,
		logger: sugar,
	}, nil
}

func (p *Postgres) DB() *gorm.DB {
	return p.db
}

func (p *Postgres) Close() error {
	if p == nil || p.db == nil {
		return nil
	}

	sqlDB, err := p.db.DB()
	if err != nil {
		p.logger.Errorw("get sql db from gorm failed", "error", err)
		return fmt.Errorf("get sql db from gorm: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- sqlDB.Close()
	}()

	select {
	case err := <-done:
		if err != nil {
			p.logger.Errorw("close postgres connection failed", "error", err)
			return fmt.Errorf("close postgres connection: %w", err)
		}
	case <-ctx.Done():
		p.logger.Warn("timeout while closing postgres connection")
		return fmt.Errorf("timeout closing postgres")
	}

	p.logger.Info("postgres disconnected")
	return nil
}