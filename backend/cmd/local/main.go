package main

import (
	"os"
	"os/signal"
	"syscall"

	"shosha-finance/internal/config"
	"shosha-finance/internal/database"
	"shosha-finance/internal/handler"
	"shosha-finance/internal/middleware"
	"shosha-finance/internal/repository"
	"shosha-finance/internal/service"
	"shosha-finance/internal/worker"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	cfg := config.LoadLocalConfig()

	db, err := database.NewConnection(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	if err := database.Migrate(db); err != nil {
		log.Fatal().Err(err).Msg("Failed to run migrations")
	}

	txRepo := repository.NewTransactionRepository(db)
	branchRepo := repository.NewBranchRepository(db)

	txService := service.NewTransactionService(txRepo)
	branchService := service.NewBranchService(branchRepo)

	txHandler := handler.NewTransactionHandler(txService, cfg.BranchID)
	dashboardHandler := handler.NewDashboardHandler(txService, cfg.BranchID)
	systemHandler := handler.NewSystemHandler(txService, cfg.CloudAPIURL)

	app := fiber.New(fiber.Config{
		AppName: "Shosha Finance Local",
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(middleware.SetupCORS())

	api := app.Group("/api/v1")

	api.Post("/transactions", txHandler.Create)
	api.Get("/transactions", txHandler.GetAll)
	api.Get("/transactions/:id", txHandler.GetByID)

	api.Get("/dashboard/summary", dashboardHandler.GetSummary)

	api.Get("/system/status", systemHandler.GetStatus)
	api.Get("/health", systemHandler.HealthCheck)

	syncWorker := worker.NewSyncWorker(txService, cfg)
	if cfg.BranchAPIKey != "" && cfg.CloudAPIURL != "" {
		syncWorker.Start()
	} else {
		log.Warn().Msg("Sync worker disabled: missing BRANCH_API_KEY or CLOUD_API_URL")
	}

	go func() {
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	log.Info().Str("port", cfg.Port).Msg("Local API server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down...")
	syncWorker.Stop()
	app.Shutdown()

	_ = branchService
}
