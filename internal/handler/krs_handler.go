package handler

import (
	"course-planner-api/internal/service"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

type KRSHandler struct {
	Service *service.KRSService
}

func NewKRSHandler(service *service.KRSService) *KRSHandler {
	return &KRSHandler{Service: service}
}

// getMahasiswaID extracts the MahasiswaID (User ID) from JWT claims
func getMahasiswaID(c *fiber.Ctx) (uuid.UUID, error) {
	return getUserIDFromContext(c)
}

// ListAvailableClasses (Req. 1)
func (h *KRSHandler) ListAvailableClasses(c *fiber.Ctx) error {
	mahasiswaID, err := getMahasiswaID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized", "details": err.Error()})
	}

	classes, err := h.Service.ListAvailableClasses(mahasiswaID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil daftar matakuliah", "details": err.Error()})
	}

	return c.JSON(fiber.Map{"data": classes, "semester": service.GetCurrentSemester()})
}

// TakeClass (Req. 2)
func (h *KRSHandler) TakeClass(c *fiber.Ctx) error {
	mahasiswaID, err := getMahasiswaID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized", "details": err.Error()})
	}

	var req struct {
		ClassIDs []string `json:"class_ids"` // bisa 1 atau banyak
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Permintaan tidak valid"})
	}

	if len(req.ClassIDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Tidak ada class yang dikirim"})
	}

	var classUUIDs []uuid.UUID
	for _, idStr := range req.ClassIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Class ID tidak valid", "details": idStr})
		}
		classUUIDs = append(classUUIDs, id)
	}

	if err := h.Service.TakeClass(mahasiswaID, classUUIDs); err != nil {
		if strings.Contains(err.Error(), "sudah diverifikasi") || strings.Contains(err.Error(), "sudah Anda ambil") || strings.Contains(err.Error(), "kode yang sama") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil matakuliah", "details": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Matakuliah berhasil ditambahkan ke KRS"})
}

// DropClass (Req. 3)
func (h *KRSHandler) DropClass(c *fiber.Ctx) error {
	mahasiswaID, err := getMahasiswaID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized", "details": err.Error()})
	}

	classID, err := uuid.Parse(c.Params("classId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Class ID tidak valid"})
	}

	if err := h.Service.DropClass(mahasiswaID, classID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Matakuliah tidak ditemukan dalam KRS aktif Anda atau statusnya tidak aktif."})
		}
		if strings.Contains(err.Error(), "sudah diverifikasi") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menghapus matakuliah", "details": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// GetTakenClasses (Req. 4)
func (h *KRSHandler) GetTakenClasses(c *fiber.Ctx) error {
	mahasiswaID, err := getMahasiswaID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Unauthorized",
			"details": err.Error(),
		})
	}

	krs, err := h.Service.GetTakenClasses(mahasiswaID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "KRS untuk semester ini belum ditemukan",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Gagal mengambil KRS",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data":     krs,
		"semester": krs.Semester,
	})
}


// RequestCancellation (Req. 5)
func (h *KRSHandler) RequestCancellation(c *fiber.Ctx) error {
	mahasiswaID, err := getMahasiswaID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized", "details": err.Error()})
	}

	classID, err := uuid.Parse(c.Params("classId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Class ID tidak valid"})
	}

	if err := h.Service.RequestClassCancellation(mahasiswaID, classID); err != nil {
		if strings.Contains(err.Error(), "belum diverifikasi") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
		}
		if strings.Contains(err.Error(), "tidak ditemukan") || strings.Contains(err.Error(), "sudah dibatalkan") || strings.Contains(err.Error(), "sudah ada") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Gagal mengajukan pembatalan", "details": err.Error()})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengajukan pembatalan", "details": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Pengajuan pembatalan matakuliah berhasil. Menunggu persetujuan Dosen PA."})
}