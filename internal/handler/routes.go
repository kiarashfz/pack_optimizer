package handler

import (
	"github.com/gofiber/fiber/v2"
	"pack_optimizer/internal/handler/pack_handler"
)

// SetupRoutes registers all application routes.
func SetupRoutes(app *fiber.App, packHandler *pack_handler.PackHandler) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})

	// api V1
	apiV1 := app.Group("/api/v1")
	// packs
	apiV1.Post("/packs/calculate", packHandler.CalculatePacks)
}
