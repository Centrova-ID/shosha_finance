package database

import (
	"fmt"

	"shosha-finance/internal/config"
	"shosha-finance/internal/models"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewConnection(cfg *config.Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	switch cfg.DBDriver {
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(cfg.SQLitePath), gormConfig)
	case "postgres":
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
		)
		db, err = gorm.Open(postgres.Open(dsn), gormConfig)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.DBDriver)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Info().Str("driver", cfg.DBDriver).Msg("Database connected successfully")
	return db, nil
}

func Migrate(db *gorm.DB) error {
	log.Info().Msg("Running database migrations...")

	err := db.AutoMigrate(
		&models.Branch{},
		&models.Transaction{},
	)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Info().Msg("Database migrations completed")
	return nil
}
