package handler

import (
	"course-planner-api/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type DosenManagementHandler struct {
	service service.DosenManagementService
}

func NewDosenManagementHandler(s service.DosenManagementService) *DosenManagementHandler {
	return &DosenManagementHandler{service: s}
}

type UpdateDosenRequest struct {
	Name *string `json:"name,omitempty"`
	NIDN *string `json:"nidn,omitempty"`
}

func (h *DosenManagementHandler) ListDosen(c *fiber.Ctx) error {
	dosens, err := h.service.GetAllDosen()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(dosens)
}

func (h *DosenManagementHandler) GetDosen(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid dosen id"})
	}

	dosen, err := h.service.GetDosenByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "dosen not found"})
	}

	return c.JSON(dosen)
}

func (h *DosenManagementHandler) UpdateDosen(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid dosen id"})
	}

	var req UpdateDosenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	dosen, err := h.service.UpdateDosen(id, req.Name, req.NIDN)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Dosen updated successfully",
		"dosen":   dosen,
	})
}
