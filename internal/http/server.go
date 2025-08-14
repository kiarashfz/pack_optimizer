// Package http initializes and configures the HTTP server.
package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"

	"pack_optimizer/templates"
)

// NewServer initializes and returns a Fiber app with routes and middleware.
func NewServer() *fiber.App {
	// Create template engine from embedded templates
	engine := html.NewFileSystem(http.FS(templates.FS), ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(recover.New())
	app.Use(logger.New())

	return app
}
