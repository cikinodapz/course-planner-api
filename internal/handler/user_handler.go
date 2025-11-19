package handler

import (
	"course-planner-api/internal/service"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service *service.UserService
}

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{s}
}

func (h *UserHandler) ListStudents(c *fiber.Ctx) error {
	users, err := h.service.ListStudents()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message":  "Failed to retrieve students",
			"error": err.Error()})
	}
	return c.JSON(Response{
	Message: "Successfully retrieved student data",
	Data:    users,
	})
}