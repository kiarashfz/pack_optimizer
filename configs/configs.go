package configs

import (
	"fmt"
	"log"
	"os"
)

// Config holds all application-wide configurations.
type Config struct {
	Port       string
	GormDSN    string
	MigrateDSN string
}

// LoadConfig loads configuration from environment variables.
func LoadConfig() *Config {
	gormDSN := os.Getenv("POSTGRES_DSN")
	if gormDSN == "" {
		log.Fatal("POSTGRES_DSN environment variable not set")
	}

	migrateDSN := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		"db",
		"5432",
		os.Getenv("POSTGRES_DB"),
		"disable",
	)

	return &Config{
		Port:       getEnv("PORT", "8080"),
		GormDSN:    gormDSN,
		MigrateDSN: migrateDSN,
	}
}

// getEnv gets an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
