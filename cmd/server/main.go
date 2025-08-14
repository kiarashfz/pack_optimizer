package main

import (
	"log"
	"pack_optimizer/configs"
	"pack_optimizer/db"
	"pack_optimizer/internal/handler"
	"pack_optimizer/internal/handler/pack_handler"
	"pack_optimizer/internal/http"
	"pack_optimizer/internal/repository/sql_repo"
	"pack_optimizer/internal/usecase/pack_usecase"

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

	// Initialize the repository, use case, and handler layers
	packRepo := sql_repo.NewPackRepo(gormDB)
	packUseCase := pack_usecase.NewPackUseCase(packRepo)
	packHandler := pack_handler.NewPackHandler(packUseCase)

	// Initialize and start the HTTP server
	app := http.NewServer()
	handler.SetupRoutes(app, packHandler)

	log.Fatal(app.Listen(":" + cfg.App.Port))
}
