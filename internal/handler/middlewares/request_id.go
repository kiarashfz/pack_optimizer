// Package middlewares provides HTTP middleware utilities for the application.
package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RequestIDMiddleware to generate request ID
func RequestIDMiddleware(c *fiber.Ctx) error {
	reqID := c.Get("X-Request-ID")
	if reqID == "" {
		reqID = uuid.New().String()
	}
	c.Locals("request_id", reqID)
	return c.Next()
}
