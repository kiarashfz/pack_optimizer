// Package handler initializes and configures the HTTP server.
package handler

import (
	"pack_optimizer/internal/handler/packhandler"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes registers all application routes.
func SetupRoutes(app *fiber.App, packHandler *packhandler.PackHandler) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})

	// api
	api := app.Group("/api")
	// api v1
	apiV1 := api.Group("/v1")
	// packs
	apiV1.Post("/packs/calculate", packHandler.CalculatePacks)
}
