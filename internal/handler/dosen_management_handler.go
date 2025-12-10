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

// ListDosen godoc
// @Summary List all dosen
// @Tags Admin - Dosen
// @Security BearerAuth
// @Success 200 {array} models.User
// @Router /api/admin/dosen [get]
func (h *DosenManagementHandler) ListDosen(c *fiber.Ctx) error {
	dosens, err := h.service.GetAllDosen()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(dosens)
}

// GetDosen godoc
// @Summary Get dosen by ID
// @Tags Admin - Dosen
// @Security BearerAuth
// @Param id path string true "Dosen ID"
// @Success 200 {object} models.User
// @Router /api/admin/dosen/{id} [get]
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

// UpdateDosen godoc
// @Summary Update dosen
// @Tags Admin - Dosen
// @Security BearerAuth
// @Accept json
// @Param id path string true "Dosen ID"
// @Param body body UpdateDosenRequest true "Dosen data"
// @Success 200 {object} models.User
// @Router /api/admin/dosen/{id} [patch]
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
