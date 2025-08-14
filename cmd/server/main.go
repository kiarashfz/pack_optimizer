package main

import (
	"log"
	"pack_optimizer/configs"
	"pack_optimizer/db"
	"pack_optimizer/internal/http"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables and configurations
	cfg := configs.LoadConfig()

	// Connect to DB
	gormDB, err := gorm.Open(postgres.Open(cfg.DB.GormDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	// Run migrations
	db.RunMigrations(cfg.DB.MigrateDSN)

	// --- Server Setup and Execution ---
	server := http.NewServer(gormDB)

	// Run the server and handle graceful shutdown.
	server.Run(cfg.App)
}
