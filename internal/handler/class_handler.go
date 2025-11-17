package handler

import (
	"course-planner-api/internal/service"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ClassHandler struct {
	classService service.ClassService
}

func NewClassHandler(classService service.ClassService) *ClassHandler {
	return &ClassHandler{classService: classService}
}

type createClassRequest struct {
	CourseID          string `json:"course_id"`
	DosenID           string `json:"dosen_id"`
	NamaKelas         string `json:"nama_kelas"`
	Hari              string `json:"hari"`
	JamMulai          string `json:"jam_mulai"`          // format "15:04"
	JamSelesai        string `json:"jam_selesai"`        // format "15:04"
	RoomID            string `json:"room_id"`
	Kuota             int    `json:"kuota"`
	SemesterPenawaran string `json:"semester_penawaran"` // "ganjil" / "genap"
}

type updateClassRequest struct {
	CourseID          *string `json:"course_id"`
	DosenID           *string `json:"dosen_id"`
	NamaKelas         *string `json:"nama_kelas"`
	Hari              *string `json:"hari"`
	JamMulai          *string `json:"jam_mulai"`          // format "15:04"
	JamSelesai        *string `json:"jam_selesai"`        // format "15:04"
	RoomID            *string `json:"room_id"`
	Kuota             *int    `json:"kuota"`
	SemesterPenawaran *string `json:"semester_penawaran"` // "ganjil" / "genap"
}

func parseTimeHM(value string) (time.Time, error) {
	return time.Parse("15:04", value)
}

func (h *ClassHandler) CreateClass(c *fiber.Ctx) error {
	var body createClassRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	courseID, err := uuid.Parse(body.CourseID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid course_id"})
	}
	dosenID, err := uuid.Parse(body.DosenID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid dosen_id"})
	}
	roomID, err := uuid.Parse(body.RoomID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid room_id"})
	}

	jamMulai, err := parseTimeHM(body.JamMulai)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid jam_mulai, expected HH:MM"})
	}
	jamSelesai, err := parseTimeHM(body.JamSelesai)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid jam_selesai, expected HH:MM"})
	}

	input := service.CreateClassInput{
		CourseID:          courseID,
		DosenID:           dosenID,
		NamaKelas:         body.NamaKelas,
		Hari:              body.Hari,
		JamMulai:          jamMulai,
		JamSelesai:        jamSelesai,
		RoomID:            roomID,
		Kuota:             body.Kuota,
		SemesterPenawaran: body.SemesterPenawaran,
	}

	class, err := h.classService.CreateClass(input)
	if err != nil {
		if err == service.ErrClassTimeConflict {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "room already used at this time"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(class)
}

func (h *ClassHandler) ListClasses(c *fiber.Ctx) error {
	classes, err := h.classService.ListClasses()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(classes)
}

func (h *ClassHandler) GetClass(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	class, err := h.classService.GetClass(id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "class not found"})
	}

	return c.JSON(class)
}

func (h *ClassHandler) UpdateClass(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var body updateClassRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	var input service.UpdateClassInput

	if body.CourseID != nil {
		courseID, err := uuid.Parse(*body.CourseID)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid course_id"})
		}
		input.CourseID = &courseID
	}

	if body.DosenID != nil {
		dosenID, err := uuid.Parse(*body.DosenID)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid dosen_id"})
		}
		input.DosenID = &dosenID
	}

	if body.NamaKelas != nil {
		input.NamaKelas = body.NamaKelas
	}

	if body.Hari != nil {
		input.Hari = body.Hari
	}

	if body.JamMulai != nil {
		jamMulai, err := parseTimeHM(*body.JamMulai)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid jam_mulai, expected HH:MM"})
		}
		input.JamMulai = &jamMulai
	}

	if body.JamSelesai != nil {
		jamSelesai, err := parseTimeHM(*body.JamSelesai)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid jam_selesai, expected HH:MM"})
		}
		input.JamSelesai = &jamSelesai
	}

	if body.RoomID != nil {
		roomID, err := uuid.Parse(*body.RoomID)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid room_id"})
		}
		input.RoomID = &roomID
	}

	if body.Kuota != nil {
		input.Kuota = body.Kuota
	}

	if body.SemesterPenawaran != nil {
		input.SemesterPenawaran = body.SemesterPenawaran
	}

	class, err := h.classService.UpdateClass(id, input)
	if err != nil {
		if err == service.ErrClassTimeConflict {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "room already used at this time"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(class)
}

func (h *ClassHandler) DeleteClass(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	if err := h.classService.DeleteClass(id); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(http.StatusNoContent)
}
