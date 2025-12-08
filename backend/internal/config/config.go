package config

import (
	"os"
	"strconv"
)

type Config struct {
	AppMode       string
	Port          string
	DBDriver      string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	SQLitePath    string
	CloudAPIURL   string
	BranchAPIKey  string
	BranchID      string
	SyncInterval  int
}

func LoadLocalConfig() *Config {
	return &Config{
		AppMode:      getEnv("APP_MODE", "local"),
		Port:         getEnv("PORT", "8080"),
		DBDriver:     "sqlite",
		SQLitePath:   getEnv("SQLITE_PATH", "./shosha_finance.db"),
		CloudAPIURL:  getEnv("CLOUD_API_URL", "http://localhost:3000"),
		BranchAPIKey: getEnv("BRANCH_API_KEY", ""),
		BranchID:     getEnv("BRANCH_ID", ""),
		SyncInterval: getEnvInt("SYNC_INTERVAL", 30),
	}
}

func LoadCloudConfig() *Config {
	return &Config{
		AppMode:    getEnv("APP_MODE", "cloud"),
		Port:       getEnv("PORT", "3000"),
		DBDriver:   "postgres",
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASS", ""),
		DBName:     getEnv("DB_NAME", "shosha_finance"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
