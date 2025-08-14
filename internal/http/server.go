// Package http provides the HTTP server implementation for the application.
package http

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"pack_optimizer/configs"
	"pack_optimizer/internal/handler/packhandler"
	"pack_optimizer/internal/repository/sqlrepo"
	"pack_optimizer/internal/usecase/packusecase"
	"pack_optimizer/templates"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"gorm.io/gorm"
)

// Server encapsulates the Fiber app and its configuration.
// It is a single, cohesive component responsible for the HTTP server.
type Server struct {
	App *fiber.App
	DB  *gorm.DB
}

// NewServer initializes and returns a new Server instance.
// This function acts as a factory, building the Fiber app with default
// settings and applying any optional configurations.
func NewServer(db *gorm.DB) *Server {
	// A new, corrected FileSystem wrapper for go:embed
	engine := html.NewFileSystem(http.FS(templates.FS), ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(recover.New()) // Recover from panics
	app.Use(logger.New())  // Log requests to the console

	return &Server{
		App: app,
		DB:  db,
	}
}

// Run starts the server and handles the graceful shutdown process.
func (s *Server) Run(appConfig configs.App) {
	// setup routes
	s.setupRoutes()

	// Create a channel to listen for OS signals.
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a goroutine so it doesn't block.
	go func() {
		// Listen() is blocking, so this will run in the background.
		if err := s.App.Listen(":" + appConfig.Port); err != nil {
			log.Fatalf("Fiber server failed to start: %v", err)
		}
	}()

	// Block until a shutdown signal is received.
	<-shutdownChan

	log.Println("Received shutdown signal. Starting graceful shutdown...")

	// Use Fiber's built-in Shutdown() method.
	if err := s.App.Shutdown(); err != nil {
		log.Fatalf("Error during graceful shutdown: %v", err)
	}
	log.Println("Fiber server has been gracefully shut down.")

	// Gracefully close the database connection.
	sqlDB, err := s.DB.DB()
	if err != nil {
		log.Printf("Failed to get DB instance for closing: %v", err)
	} else {
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error during database shutdown: %v", err)
		}
	}
	log.Println("Database connection has been gracefully closed.")

	log.Println("Server has been successfully shut down.")
}

func (s *Server) setupRoutes() {
	packHandler := packhandler.NewPackHandler(packusecase.NewPackUseCase(sqlrepo.NewPackRepo(s.DB)))

	s.App.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})

	// api
	api := s.App.Group("/api")
	// api v1
	apiV1 := api.Group("/v1")
	// packs
	apiV1.Post("/packs/calculate", packHandler.CalculatePacks)
}
