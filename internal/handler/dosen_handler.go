package handler

import (
	"course-planner-api/internal/service"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DosenHandler struct {
	Service *service.DosenPAService
}

func NewDosenHandler(service *service.DosenPAService) *DosenHandler {
	return &DosenHandler{Service: service}
}

func (h *DosenHandler) ListStudents(c *fiber.Ctx) error {
	dosenID, err := getUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized", "details": err.Error()})
	}

	students, err := h.Service.ListMahasiswa(dosenID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil daftar mahasiswa", "details": err.Error()})
	}

	return c.JSON(fiber.Map{"data": students})
}

func (h *DosenHandler) GetMahasiswaKRS(c *fiber.Ctx) error {
	dosenID, err := getUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized", "details": err.Error()})
	}

	mahasiswaID, err := uuid.Parse(c.Params("mahasiswaId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Mahasiswa ID tidak valid"})
	}

	krs, err := h.Service.GetMahasiswaKRS(dosenID, mahasiswaID)
	if err != nil {
		if strings.Contains(err.Error(), "bukan bimbingan") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "KRS tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil KRS", "details": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data":     krs,
		"semester": krs.Semester,
	})
}

func (h *DosenHandler) RemoveMahasiswaClass(c *fiber.Ctx) error {
	dosenID, err := getUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized", "details": err.Error()})
	}

	mahasiswaID, err := uuid.Parse(c.Params("mahasiswaId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Mahasiswa ID tidak valid"})
	}

	classID, err := uuid.Parse(c.Params("classId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Class ID tidak valid"})
	}

	if err := h.Service.RemoveMahasiswaClass(dosenID, mahasiswaID, classID); err != nil {
		if strings.Contains(err.Error(), "bukan bimbingan") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Matakuliah tidak ditemukan di KRS mahasiswa ini"})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Matakuliah berhasil dihapus dari KRS mahasiswa ini"})
}

func (h *DosenHandler) UpdateMahasiswaClass(c *fiber.Ctx) error {
	dosenID, err := getUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized", "details": err.Error()})
	}

	mahasiswaID, err := uuid.Parse(c.Params("mahasiswaId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Mahasiswa ID tidak valid"})
	}

	classID, err := uuid.Parse(c.Params("classId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Class ID tidak valid"})
	}

	var body struct {
		NewClassID string `json:"new_class_id"`
	}
	if err := c.BodyParser(&body); err != nil || body.NewClassID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Request body tidak valid atau new_class_id kosong"})
	}

	newClassID, err := uuid.Parse(body.NewClassID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "new_class_id tidak valid"})
	}

	if err := h.Service.UpdateMahasiswaClass(dosenID, mahasiswaID, classID, newClassID); err != nil {
		if strings.Contains(err.Error(), "bukan bimbingan") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Matakuliah tidak ditemukan di KRS mahasiswa ini"})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Kelas matakuliah berhasil diperbarui"})
}

func (h *DosenHandler) ApproveMahasiswaClass(c *fiber.Ctx) error {
	return h.updateStatus(c, true)
}

func (h *DosenHandler) RejectMahasiswaClass(c *fiber.Ctx) error {
	return h.updateStatus(c, false)
}

func (h *DosenHandler) updateStatus(c *fiber.Ctx, approve bool) error {
	dosenID, err := getUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized", "details": err.Error()})
	}

	mahasiswaID, err := uuid.Parse(c.Params("mahasiswaId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Mahasiswa ID tidak valid"})
	}

	classID, err := uuid.Parse(c.Params("classId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Class ID tidak valid"})
	}

	var actionErr error
	if approve {
		actionErr = h.Service.ApproveMahasiswaClass(dosenID, mahasiswaID, classID)
	} else {
		actionErr = h.Service.RejectMahasiswaClass(dosenID, mahasiswaID, classID)
	}

	if actionErr != nil {
		if strings.Contains(actionErr.Error(), "bukan bimbingan") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": actionErr.Error()})
		}
		if errors.Is(actionErr, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Matakuliah tidak ditemukan di KRS mahasiswa ini"})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": actionErr.Error()})
	}

	if approve {
		return c.JSON(fiber.Map{"message": "Matakuliah disetujui."})
	}
	return c.JSON(fiber.Map{"message": "Matakuliah ditolak."})
}