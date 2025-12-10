package handler

import (
	"course-planner-api/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CourseHandler struct {
	service service.CourseService
}

func NewCourseHandler(s service.CourseService) *CourseHandler {
	return &CourseHandler{service: s}
}

type CreateCourseRequest struct {
	Kode string `json:"kode"`
	Nama string `json:"nama"`
	SKS  int    `json:"sks"`
}

type UpdateCourseRequest struct {
	Kode *string `json:"kode,omitempty"`
	Nama *string `json:"nama,omitempty"`
	SKS  *int    `json:"sks,omitempty"`
}

// ListCourses godoc
// @Summary List all courses
// @Tags Admin - Courses
// @Security BearerAuth
// @Success 200 {array} models.Course
// @Router /api/admin/courses [get]
func (h *CourseHandler) ListCourses(c *fiber.Ctx) error {
	courses, err := h.service.GetAllCourses()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(courses)
}

// CreateCourse godoc
// @Summary Create a new course
// @Tags Admin - Courses
// @Security BearerAuth
// @Accept json
// @Param body body CreateCourseRequest true "Course data"
// @Success 201 {object} models.Course
// @Router /api/admin/courses [post]
func (h *CourseHandler) CreateCourse(c *fiber.Ctx) error {
	var req CreateCourseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if req.Kode == "" || req.Nama == "" || req.SKS <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "kode, nama, and sks are required (sks > 0)"})
	}

	course, err := h.service.CreateCourse(req.Kode, req.Nama, req.SKS)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Course created successfully",
		"course":  course,
	})
}

// GetCourse godoc
// @Summary Get course by ID
// @Tags Admin - Courses
// @Security BearerAuth
// @Param id path string true "Course ID"
// @Success 200 {object} models.Course
// @Router /api/admin/courses/{id} [get]
func (h *CourseHandler) GetCourse(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid course id"})
	}

	course, err := h.service.GetCourseByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "course not found"})
	}

	return c.JSON(course)
}

// UpdateCourse godoc
// @Summary Update course
// @Tags Admin - Courses
// @Security BearerAuth
// @Accept json
// @Param id path string true "Course ID"
// @Param body body UpdateCourseRequest true "Course data"
// @Success 200 {object} models.Course
// @Router /api/admin/courses/{id} [patch]
func (h *CourseHandler) UpdateCourse(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid course id"})
	}

	var req UpdateCourseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	course, err := h.service.UpdateCourse(id, req.Kode, req.Nama, req.SKS)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Course updated successfully",
		"course":  course,
	})
}

// DeleteCourse godoc
// @Summary Delete course
// @Tags Admin - Courses
// @Security BearerAuth
// @Param id path string true "Course ID"
// @Success 204 "No Content"
// @Router /api/admin/courses/{id} [delete]
func (h *CourseHandler) DeleteCourse(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid course id"})
	}

	if err := h.service.DeleteCourse(id); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
