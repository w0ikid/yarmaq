package account

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *handler) GetByAccountID(c *fiber.Ctx) error {
	idStr := c.Params("accountId")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid account ID"})
	}

	entries, err := h.ledgerDomain.GetLedgerUsecase.ExecuteByAccountID(c.Context(), id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(200).JSON(entries)
}
