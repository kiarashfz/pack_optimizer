package customerrrors

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// ErrorHandler is a Central error handler
func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	// Fetch request ID from context
	reqID, tErr := c.Locals("request_id").(string)
	if !tErr {
		// If request ID is not set, generate a new one
		reqID = c.Get("X-Request-ID")
		if reqID == "" {
			reqID = "unknown"
		}
	}

	// Log the error with structured fields
	log.Error().
		Err(err).
		Str("request_id", reqID).
		Str("method", c.Method()).
		Str("path", c.Path()).
		Str("ip", c.IP()).
		Str("user_agent", c.Get("User-Agent")).
		Int("status_code", code).
		Msg("request failed")

	// Respond with generic error
	return c.Status(code).JSON(fiber.Map{
		"error":  err.Error(),
		"req_id": reqID,
	})
}
