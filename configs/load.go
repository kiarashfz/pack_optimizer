// Package configs provides functionality to load and validate application configuration.
package configs

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/spf13/viper"
)

// LoadConfig loads config from .env file and environment variables
func LoadConfig() Config {
	v := viper.New()

	// --- Load .env file ---
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	// --- Environment variables override ---
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		log.Printf("No .env file found, relying on environment variables only: %v", err)
	}

	v.SetEnvKeyReplacer(strings.NewReplacer(".", ""))

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("unable to decode config: %v", err)
	}

	// --- Validation for required fields using a validation library ---
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	// Generate DSNs from loaded values
	cfg.DB.GormDSN = fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.SSLMode,
	)

	cfg.DB.MigrateDSN = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name, cfg.DB.SSLMode,
	)

	return cfg
}
