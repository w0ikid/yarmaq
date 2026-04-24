package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"

	"github.com/w0ikid/yarmaq/pkg/config"
	"github.com/w0ikid/yarmaq/pkg/jwks"
	"github.com/w0ikid/yarmaq/pkg/zitadel"

	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/repo"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/repo/igorm"

	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/container"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/handlers"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/handlers/v1/account"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/handlers/v1/internals"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/handlers/v1/ledger"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/handlers/v1/webhook"

	kafkamodule "github.com/w0ikid/yarmaq/pkg/kafka_module"
	"github.com/w0ikid/yarmaq/pkg/outbox_worker"

	accountsv1 "github.com/w0ikid/yarmaq/pkg/gen/accounts/v1"
	grpc_v1 "github.com/w0ikid/yarmaq/apps/accounts-service/internal/handlers/grpc/v1"
	"google.golang.org/grpc"
	"net"
)

type App struct {
	fapp      *fiber.App
	grpcServer *grpc.Server
	addr      string
	grpcAddr  string
	container *container.Container
	logger    *zap.SugaredLogger
	pg        *repo.Postgres
	cancel    context.CancelFunc

	kafkaPublisher *kafkamodule.Publisher
	outboxWorker   *outbox_worker.Worker
}

func NewApp(ctx context.Context, cfg config.Config, logger *zap.SugaredLogger) (*App, error) {
	appLogger := logger.Named("app")

	// postgres timeout
	pgCtx, cancelPg := context.WithTimeout(ctx, 5*time.Second)
	defer cancelPg()

	pg, err := repo.NewPostgres(pgCtx, cfg.Postgres, appLogger)
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}

	appLogger.Info("jwks: ", cfg.Zitadel.JWKSURL)

	// JWKS client
	jwksClient, err := jwks.New(cfg.Zitadel.JWKSURL)
	if err != nil {
		return nil, fmt.Errorf("init jwks client: %w", err)
	}

	// zitadel client
	zitadelClient, err := zitadel.New(ctx, cfg.Zitadel.Domain, cfg.Zitadel.API, cfg.Zitadel.KeyPath)
	if err != nil {
		return nil, fmt.Errorf("init zitadel client: %w", err)
	}

	appLogger.Info("zitadel client initialized", zitadelClient)
	appLogger.Info("zitadel domain: ", cfg.Zitadel.Domain)
	appLogger.Info("zitadel keyPath: ", cfg.Zitadel.KeyPath)

	// kafka publisher
	kafkaPublisher, err := kafkamodule.NewPublisher(kafkamodule.Config{
		Brokers: cfg.Kafka.Brokers,
	}, appLogger)
	if err != nil {
		return nil, fmt.Errorf("init kafka publisher: %w", err)
	}

	// outbox worker
	outboxWorker := outbox_worker.NewWorker(
		pg.DB(),
		kafkaPublisher,
		appLogger,
		outbox_worker.WithInterval(3*time.Second),
		outbox_worker.WithBatchSize(50),
	)

	// kafka consumers
	// consumers := []*kafkamodule.Consumer{
		
	// }
	// Репозитории
	repositories := igorm.NewGormRepository(pg.DB(), appLogger)

	// DI контейнер
	cont := container.NewContainer(ctx, repositories, zitadelClient, appLogger)

	// Handlers
	h := handlers.NewHandlers(handlers.Depedencies{
		AccountDeps: account.HandlerDeps{
			AccountDomain: cont.AccountDomain,
			Logger:        appLogger,
		},
		InternalDeps: internals.HandlerDeps{
			AccountDomain: cont.AccountDomain,
			Logger:        appLogger,
		},
		LedgerDeps: ledger.HandlerDeps{
			LedgerDomain: cont.LedgerDomain,
			Logger:       appLogger,
		},
		WebhookDeps: webhook.HandlerDeps{
			AccountDomain: cont.AccountDomain,
			Logger:        appLogger,
		},
		JWKS: jwksClient,
	})

	fapp := fiber.New(fiber.Config{
		AppName:      "accounts-service",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	})
	fapp.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			appLogger.Errorw("panic recovered", "error", e)
		},
	}))

	

	// Fiber router
	router := handlers.NewRouter(fapp, h)
	router.SetupRoutes(appLogger)

	// gRPC server
	grpcServer := grpc.NewServer()
	accountsGrpcHandler := grpc_v1.NewAccountsHandler(cont.AccountDomain, appLogger)
	accountsv1.RegisterAccountsServiceServer(grpcServer, accountsGrpcHandler)

	// Контекст приложения
	_, cancel := context.WithCancel(ctx)

	return &App{
		fapp:       fapp,
		grpcServer: grpcServer,
		addr:       ":" + cfg.HTTP.Port,
		grpcAddr:   ":" + cfg.GRPC.Port,
		container:  cont,
		logger:    appLogger,
		pg:        pg,
		cancel:    cancel,

		kafkaPublisher: kafkaPublisher,
		outboxWorker:   outboxWorker,
	}, nil
}

// Start запускает HTTP сервер
func (a *App) Start(ctx context.Context) error {
	go a.outboxWorker.Run(ctx)

	// Start gRPC server
	lis, err := net.Listen("tcp", a.grpcAddr)
	if err != nil {
		return fmt.Errorf("failed to listen gRPC: %w", err)
	}
	go func() {
		a.logger.Info("starting gRPC server", zap.String("addr", a.grpcAddr))
		if err := a.grpcServer.Serve(lis); err != nil {
			a.logger.Error("gRPC server failed", zap.Error(err))
		}
	}()

	a.logger.Info("starting HTTP server", zap.String("addr", a.addr))
	if err := a.fapp.Listen(a.addr); err != nil {
		return fmt.Errorf("fiber server: %w", err)
	}
	return nil
}

func (a *App) Stop(ctx context.Context) error {
	a.cancel()

	shutdownCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var errOccurred bool

	// shut down service
	if err := a.fapp.ShutdownWithContext(shutdownCtx); err != nil {
		a.logger.Error("fiber shutdown failed", zap.Error(err))
		errOccurred = true
	} else {
		a.logger.Info("fiber server stopped gracefully")
	}

	// stop gRPC server
	a.grpcServer.GracefulStop()
	a.logger.Info("gRPC server stopped gracefully")

	// close postgres
	if err := a.pg.Close(); err != nil {
		a.logger.Error("postgres close failed", zap.Error(err))
		errOccurred = true
	} else {
		a.logger.Info("postgres connection closed")
	}
	// close kafka publisher
	if err := a.kafkaPublisher.Close(); err != nil {
		a.logger.Error("kafka publisher close failed", zap.Error(err))
		errOccurred = true
	}

	if errOccurred {
		return fmt.Errorf("some resources failed to close, check logs")
	}

	a.logger.Info("app stopped gracefully")
	return nil
}
