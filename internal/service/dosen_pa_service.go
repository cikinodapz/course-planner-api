package service

import (
	"course-planner-api/internal/models"
	"course-planner-api/internal/repository"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DosenPAService struct {
	Repo *repository.KRSRepository
}

func NewDosenPAService(repo *repository.KRSRepository) *DosenPAService {
	return &DosenPAService{Repo: repo}
}

// ListMahasiswa mengembalikan daftar mahasiswa bimbingan dosen PA
func (s *DosenPAService) ListMahasiswa(dosenID uuid.UUID) ([]models.User, error) {
	return s.Repo.ListMahasiswaByDosenPA(dosenID)
}

// GetMahasiswaKRS mengambil KRS mahasiswa bimbingan di semester berjalan
func (s *DosenPAService) GetMahasiswaKRS(dosenID uuid.UUID, mahasiswaID uuid.UUID) (*models.KRS, error) {
	return s.getKRSForAdvisee(dosenID, mahasiswaID)
}

// RemoveMahasiswaClass menghapus matakuliah yang dipilih mahasiswa dari KRS
func (s *DosenPAService) RemoveMahasiswaClass(dosenID uuid.UUID, mahasiswaID uuid.UUID, classID uuid.UUID) error {
	krs, err := s.getKRSForAdvisee(dosenID, mahasiswaID)
	if err != nil {
		return err
	}

	return s.Repo.DB.Transaction(func(tx *gorm.DB) error {
		res := tx.Where("krs_id = ? AND class_id = ?", krs.ID, classID).
			Delete(&models.KRSItem{})

		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return nil
	})
}

// UpdateMahasiswaClass memindahkan matakuliah yang sudah dipilih ke kelas lain
func (s *DosenPAService) UpdateMahasiswaClass(dosenID uuid.UUID, mahasiswaID uuid.UUID, classID uuid.UUID, newClassID uuid.UUID) error {
	krs, err := s.getKRSForAdvisee(dosenID, mahasiswaID)
	if err != nil {
		return err
	}

	var targetItem *models.KRSItem
	for i := range krs.Items {
		if krs.Items[i].ClassID == classID {
			targetItem = &krs.Items[i]
			break
		}
	}

	if targetItem == nil {
		return gorm.ErrRecordNotFound
	}

	if targetItem.Status == KRS_ITEM_STATUS_CANCELLED {
		return errors.New("Matakuliah ini sudah dibatalkan.")
	}

	classes, err := s.Repo.GetClassesByIDs([]uuid.UUID{newClassID})
	if err != nil {
		return err
	}
	newClass := classes[0]

	var otherItems []models.KRSItem
	for _, item := range krs.Items {
		if item.ClassID == classID {
			continue
		}
		if item.Status == KRS_ITEM_STATUS_CANCELLED || item.Status == KRS_ITEM_STATUS_REJECTED {
			continue
		}
		otherItems = append(otherItems, item)
	}

	if err := CheckScheduleConflict([]models.Class{newClass}, otherItems); err != nil {
		return err
	}

	for _, item := range otherItems {
		if item.ClassID == newClassID {
			return errors.New("Kelas ini sudah ada di KRS mahasiswa.")
		}
		if item.Class.CourseID == newClass.CourseID {
			return fmt.Errorf("Mahasiswa sudah mengambil mata kuliah dengan kode %s di kelas lain.", newClass.Course.Kode)
		}
	}

	return s.Repo.UpdateKRSItemClass(krs.ID, classID, newClassID, KRS_ITEM_STATUS_ACTIVE)
}

// ApproveMahasiswaClass menandai matakuliah sebagai disetujui oleh dosen PA
func (s *DosenPAService) ApproveMahasiswaClass(dosenID uuid.UUID, mahasiswaID uuid.UUID, classID uuid.UUID) error {
	return s.setItemStatusAndVerify(dosenID, mahasiswaID, classID, KRS_ITEM_STATUS_APPROVED)
}

// RejectMahasiswaClass menandai matakuliah sebagai ditolak oleh dosen PA
func (s *DosenPAService) RejectMahasiswaClass(dosenID uuid.UUID, mahasiswaID uuid.UUID, classID uuid.UUID) error {
	return s.setItemStatusAndVerify(dosenID, mahasiswaID, classID, KRS_ITEM_STATUS_REJECTED)
}

func (s *DosenPAService) setItemStatusAndVerify(dosenID uuid.UUID, mahasiswaID uuid.UUID, classID uuid.UUID, status string) error {
	krs, err := s.getKRSForAdvisee(dosenID, mahasiswaID)
	if err != nil {
		return err
	}

	var targetItem *models.KRSItem
	for i := range krs.Items {
		if krs.Items[i].ClassID == classID {
			targetItem = &krs.Items[i]
			break
		}
	}

	if targetItem == nil {
		return gorm.ErrRecordNotFound
	}

	if targetItem.Status == KRS_ITEM_STATUS_CANCELLED {
		return errors.New("Matakuliah ini sudah dibatalkan.")
	}

	if targetItem.Status == status {
		return fmt.Errorf("Status matakuliah sudah %s.", status)
	}

	now := time.Now()
	return s.Repo.DB.Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&models.KRSItem{}).
			Where("krs_id = ? AND class_id = ?", krs.ID, classID).
			Update("status", status)

		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		if err := tx.Model(&models.KRS{}).
			Where("id = ?", krs.ID).
			Updates(map[string]interface{}{
				"status":      KRS_STATUS_VERIFIED,
				"verified_at": now,
			}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *DosenPAService) getKRSForAdvisee(dosenID uuid.UUID, mahasiswaID uuid.UUID) (*models.KRS, error) {
	semester := GetCurrentSemester()
	krs, err := s.Repo.GetOrCreateKRS(mahasiswaID, semester, KRS_STATUS_DRAFT, KRS_STATUS_VERIFIED)
	if err != nil {
		return nil, err
	}

	if krs.Mahasiswa.DosenPAID == nil || *krs.Mahasiswa.DosenPAID != dosenID {
		return nil, errors.New("Mahasiswa ini bukan bimbingan Anda.")
	}

	return krs, nil
}