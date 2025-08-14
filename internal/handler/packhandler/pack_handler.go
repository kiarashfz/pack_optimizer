// Package packhandler provides HTTP handlers for pack-related operations.
package packhandler

import (
	"errors"
	"pack_optimizer/internal/domain"
	"pack_optimizer/internal/handler/customerrrors"
	"pack_optimizer/internal/usecase/packusecase"
	"pack_optimizer/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type PackHandler struct {
	packUseCase *packusecase.PackUseCase
}

func NewPackHandler(packUseCase *packusecase.PackUseCase) *PackHandler {
	return &PackHandler{packUseCase: packUseCase}
}

func (h *PackHandler) CalculatePacks(c *fiber.Ctx) error {
	var req CalculatePacksReq

	// Parse and validate JSON body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	// Validate the input
	if err := validator.Validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	output, err := h.packUseCase.CalculatePacks(c.Context(), req.Quantity)
	if err != nil {
		if errors.Is(err, domain.ErrNoPacksAvailable) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}

		// For all other errors, we return a 500.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": customerrrors.ErrUnexpected})
	}
	return c.Status(fiber.StatusOK).JSON(output)
}
