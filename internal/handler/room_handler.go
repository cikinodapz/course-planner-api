package handler

import (
	"course-planner-api/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type RoomHandler struct {
	service service.RoomService
}

func NewRoomHandler(s service.RoomService) *RoomHandler {
	return &RoomHandler{service: s}
}

type CreateRoomRequest struct {
	Nama string `json:"nama"`
}

type UpdateRoomRequest struct {
	Nama *string `json:"nama,omitempty"`
}

// ListRooms godoc
// @Summary List all rooms
// @Tags Admin - Rooms
// @Security BearerAuth
// @Success 200 {array} models.Room
// @Router /api/admin/rooms [get]
func (h *RoomHandler) ListRooms(c *fiber.Ctx) error {
	rooms, err := h.service.GetAllRooms()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(rooms)
}

// CreateRoom godoc
// @Summary Create a new room
// @Tags Admin - Rooms
// @Security BearerAuth
// @Accept json
// @Param body body CreateRoomRequest true "Room data"
// @Success 201 {object} models.Room
// @Router /api/admin/rooms [post]
func (h *RoomHandler) CreateRoom(c *fiber.Ctx) error {
	var req CreateRoomRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if req.Nama == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "nama is required"})
	}

	room, err := h.service.CreateRoom(req.Nama)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Room created successfully",
		"room":    room,
	})
}

// GetRoom godoc
// @Summary Get room by ID
// @Tags Admin - Rooms
// @Security BearerAuth
// @Param id path string true "Room ID"
// @Success 200 {object} models.Room
// @Router /api/admin/rooms/{id} [get]
func (h *RoomHandler) GetRoom(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid room id"})
	}

	room, err := h.service.GetRoomByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "room not found"})
	}

	return c.JSON(room)
}

// UpdateRoom godoc
// @Summary Update room
// @Tags Admin - Rooms
// @Security BearerAuth
// @Accept json
// @Param id path string true "Room ID"
// @Param body body UpdateRoomRequest true "Room data"
// @Success 200 {object} models.Room
// @Router /api/admin/rooms/{id} [patch]
func (h *RoomHandler) UpdateRoom(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid room id"})
	}

	var req UpdateRoomRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	room, err := h.service.UpdateRoom(id, req.Nama)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Room updated successfully",
		"room":    room,
	})
}

// DeleteRoom godoc
// @Summary Delete room
// @Tags Admin - Rooms
// @Security BearerAuth
// @Param id path string true "Room ID"
// @Success 204 "No Content"
// @Router /api/admin/rooms/{id} [delete]
func (h *RoomHandler) DeleteRoom(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid room id"})
	}

	if err := h.service.DeleteRoom(id); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
