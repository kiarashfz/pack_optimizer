package main

import (
	"pack_optimizer/configs"
	"pack_optimizer/db"
	"pack_optimizer/internal/http"
	"pack_optimizer/pkg/logpkg"

	"github.com/rs/zerolog/log"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	logpkg.InitLogger()

	// Load environment variables and configurations
	cfg := configs.LoadConfig()

	// Connect to DB
	gormDB, err := gorm.Open(postgres.Open(cfg.DB.GormDSN), &gorm.Config{})
	if err != nil {
		log.Info().Msg("Starting server on :3000")
	}

	// Run migrations
	db.RunMigrations(cfg.DB.MigrateDSN)

	// --- Server Setup and Execution ---
	server := http.NewServer(gormDB)

	// Run the server and handle graceful shutdown.
	server.Run(cfg.App)
}
