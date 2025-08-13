package pack_handler

import (
	"pack_optimizer/internal/usecase/pack_usecase"
	"pack_optimizer/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type PackHandler struct {
	packUseCase *pack_usecase.PackUseCase
}

func NewPackHandler(packUseCase *pack_usecase.PackUseCase) *PackHandler {
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(output)
}
