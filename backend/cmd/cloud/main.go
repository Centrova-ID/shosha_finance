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

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	cfg := config.LoadCloudConfig()

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

	txService := service.NewTransactionService(txRepo)
	branchService := service.NewBranchService(branchRepo)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)

	// Create default admin user for cloud
	if err := authService.CreateDefaultUsers(); err != nil {
		log.Warn().Err(err).Msg("Failed to create default users")
	}

	syncHandler := handler.NewSyncHandler(txService, branchService)
	authHandler := handler.NewAuthHandler(authService)
	branchHandler := handler.NewBranchHandler(branchService)
	txHandler := handler.NewTransactionHandler(txService)
	dashboardHandler := handler.NewDashboardHandler(txService)

	app := fiber.New(fiber.Config{
		AppName: "Shosha Finance Cloud",
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(middleware.SetupCORS())

	api := app.Group("/api/v1")

	// Public routes
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "OK",
			"data":    fiber.Map{"status": "healthy"},
		})
	})
	api.Post("/auth/login", authHandler.Login)

	// Sync routes (uses API key auth)
	syncGroup := api.Group("/sync")
	syncGroup.Post("/push", syncHandler.Push)
	syncGroup.Get("/pull", syncHandler.Pull)

	// Protected routes (uses JWT auth)
	protected := api.Group("", middleware.JWTAuth(authService))

	protected.Get("/auth/me", authHandler.Me)
	protected.Post("/auth/logout", authHandler.Logout)

	protected.Get("/branches", branchHandler.GetAll)
	protected.Get("/branches/active", branchHandler.GetActive)
	protected.Get("/branches/:id", branchHandler.GetByID)
	protected.Post("/branches", branchHandler.Create)
	protected.Put("/branches/:id", branchHandler.Update)
	protected.Delete("/branches/:id", branchHandler.Delete)

	protected.Get("/transactions", txHandler.GetAll)
	protected.Get("/transactions/:id", txHandler.GetByID)
	protected.Post("/transactions", txHandler.Create)

	protected.Get("/dashboard/summary", dashboardHandler.GetSummary)

	go func() {
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	log.Info().Str("port", cfg.Port).Msg("Cloud API server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down...")
	app.Shutdown()
}
