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

	txService := service.NewTransactionService(txRepo)
	branchService := service.NewBranchService(branchRepo)

	syncHandler := handler.NewSyncHandler(txService, branchService)

	app := fiber.New(fiber.Config{
		AppName: "Shosha Finance Cloud",
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(middleware.SetupCORS())

	api := app.Group("/api/v1")

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "OK",
			"data":    fiber.Map{"status": "healthy"},
		})
	})

	syncGroup := api.Group("/sync")
	syncGroup.Use(middleware.BranchAuth(branchService))
	syncGroup.Post("/push", syncHandler.Push)

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
