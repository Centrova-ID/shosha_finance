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
	userRepo := repository.NewUserRepository(db)
	incomeRepo := repository.NewIncomeEntryRepository(db)
	expenseRepo := repository.NewExpenseEntryRepository(db)

	txService := service.NewTransactionService(txRepo)
	branchService := service.NewBranchService(branchRepo)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	incomeService := service.NewIncomeEntryService(incomeRepo, branchRepo)
	expenseService := service.NewExpenseEntryService(expenseRepo, branchRepo)

	if err := authService.CreateDefaultUsers(); err != nil {
		log.Warn().Err(err).Msg("Failed to create default users")
	}

	if err := branchService.CreateDefaultBranches(); err != nil {
		log.Warn().Err(err).Msg("Failed to create default branches")
	}

	// Initialize sync worker
	syncWorker := worker.NewSyncWorker(db, cfg)
	if cfg.CloudAPIURL != "" {
		syncWorker.Start()
	} else {
		log.Warn().Msg("Sync worker disabled: CLOUD_API_URL not set")
	}

	txHandler := handler.NewTransactionHandler(txService)
	dashboardHandler := handler.NewDashboardHandler(txService)
	systemHandler := handler.NewSystemHandler(txService, syncWorker)
	authHandler := handler.NewAuthHandler(authService)
	branchHandler := handler.NewBranchHandler(branchService)
	incomeHandler := handler.NewIncomeEntryHandler(incomeService)
	expenseHandler := handler.NewExpenseEntryHandler(expenseService)

	app := fiber.New(fiber.Config{
		AppName: "Shosha Finance Local",
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(middleware.SetupCORS())

	api := app.Group("/api/v1")

	// Public routes
	api.Post("/auth/login", authHandler.Login)
	api.Get("/health", systemHandler.HealthCheck)

	// Protected routes
	protected := api.Group("", middleware.JWTAuth(authService))

	protected.Get("/auth/me", authHandler.Me)
	protected.Post("/auth/logout", authHandler.Logout)

	protected.Post("/transactions", txHandler.Create)
	protected.Get("/transactions", txHandler.GetAll)
	protected.Get("/transactions/:id", txHandler.GetByID)

	protected.Get("/branches", branchHandler.GetAll)
	protected.Get("/branches/active", branchHandler.GetActive)
	protected.Get("/branches/:id", branchHandler.GetByID)
	protected.Post("/branches", branchHandler.Create)
	protected.Put("/branches/:id", branchHandler.Update)
	protected.Delete("/branches/:id", branchHandler.Delete)

	protected.Get("/dashboard/summary", dashboardHandler.GetSummary)

	// Income entries
	protected.Post("/income-entries", incomeHandler.Create)
	protected.Get("/income-entries", incomeHandler.GetAll)
	protected.Get("/income-entries/range", incomeHandler.GetByDateRange)
	protected.Get("/income-entries/:id", incomeHandler.GetByID)
	protected.Put("/income-entries/:id", incomeHandler.Update)
	protected.Delete("/income-entries/:id", incomeHandler.Delete)

	// Expense entries
	protected.Post("/expense-entries", expenseHandler.Create)
	protected.Get("/expense-entries", expenseHandler.GetAll)
	protected.Get("/expense-entries/range", expenseHandler.GetByDateRange)
	protected.Get("/expense-entries/:id", expenseHandler.GetByID)
	protected.Put("/expense-entries/:id", expenseHandler.Update)
	protected.Delete("/expense-entries/:id", expenseHandler.Delete)

	protected.Get("/system/status", systemHandler.GetStatus)

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
}
