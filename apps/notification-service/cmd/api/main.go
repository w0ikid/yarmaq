package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	pkgLogger "github.com/w0ikid/yarmaq/pkg/logger"
	"github.com/w0ikid/yarmaq/pkg/config"
	"github.com/w0ikid/yarmaq/apps/notification-service/internal"
)

func main() {
	// CONFIG
	cfg := config.Load("NOTIFICATION")

	// LOGGER
	logger, err := pkgLogger.New(cfg.LogLevel)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	// CONTEXT APP
	ctx := context.Background()
	app, err := internal.NewApp(ctx, cfg, logger.Sugar())
	if err != nil {
		logger.Fatal("failed to start app", zap.Error(err))
	}

	// SERVER
	serverErrCh := make(chan error, 1)
	go func() {
		serverErrCh <- app.Start(ctx)
	}()

	// SHUTDOWN
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrCh:
		if err != nil {
			logger.Fatal("app server failed", zap.Error(err))
		}
	case sig := <-shutdown:
		logger.Info("shutdown signal received", zap.String("signal", sig.String()))
	}

	stopCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	if err := app.Stop(stopCtx); err != nil {
		logger.Fatal("failed to stop app", zap.Error(err))
	}

	logger.Info("stopped gracefully")
}
