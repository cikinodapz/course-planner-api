package repository

import (
	"course-planner-api/internal/models"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KRSRepository struct {
	DB *gorm.DB
}

func NewKRSRepository(db *gorm.DB) *KRSRepository {
	return &KRSRepository{DB: db}
}

func (r *KRSRepository) GetClassesByIDs(classIDs []uuid.UUID) ([]models.Class, error) {
	var classes []models.Class
	err := r.DB.Where("id IN ?", classIDs).
		Preload("Course").
		Preload("Dosen").
		Find(&classes).Error
	if err != nil {
		return nil, err
	}
	if len(classes) == 0 {
		return nil, errors.New("kelas tidak ditemukan")
	}
	return classes, nil
}

func (r *KRSRepository) AddItemsBatch(krsID uuid.UUID, classIDs []uuid.UUID, status string) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		for _, classID := range classIDs {
			item := models.KRSItem{
				KRSID:     krsID,
				ClassID:   classID,
				Status:    status,
				CreatedAt: time.Now(),
			}
			if err := tx.Create(&item).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// Modifikasi GetOrCreateKRS supaya preload Mahasiswa juga
func (r *KRSRepository) GetOrCreateKRS(mahasiswaID uuid.UUID, semester string, krsStatusDraft string, krsStatusVerified string) (*models.KRS, error) {
	var krs models.KRS

	err := r.DB.Where("mahasiswa_id = ? AND semester = ? AND status IN (?, ?)",
		mahasiswaID, semester, krsStatusDraft, krsStatusVerified).
		Preload("Mahasiswa").      // tambahkan ini
		Preload("Items").
		Preload("Items.Class").
		Preload("Items.KRS.Mahasiswa").
		Preload("Items.Class.Course").
		Preload("Items.Class.Dosen").
		Preload("Items.Class.Room").
		First(&krs).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		newKRS := models.KRS{
			MahasiswaID: mahasiswaID,
			Semester:    semester,
			Status:      krsStatusDraft,
			CreatedAt:   time.Now(),
		}
		if createErr := r.DB.Create(&newKRS).Error; createErr != nil {
			return nil, createErr
		}
		// preload Mahasiswa setelah membuat KRS baru
		r.DB.Preload("Mahasiswa").First(&newKRS, "id = ?", newKRS.ID)
		return &newKRS, nil
	} else if err != nil {
		return nil, err
	}

	return &krs, nil
}


func (r *KRSRepository) GetKRSByMahasiswaID(mahasiswaID uuid.UUID, semester string) (*models.KRS, error) {
	var krs models.KRS

	err := r.DB.
		Where("mahasiswa_id = ? AND semester = ?", mahasiswaID, semester).
		Preload("Mahasiswa").                    // preload user/mahasiswa
		Preload("Items").                        // preload KRS items
		Preload("Items.Class").                  // preload class di setiap item
		Preload("Items.Class.Course").           // preload course di class
		Preload("Items.Class.Dosen").            // preload dosen di class
		Preload("Items.Class.Room").             // preload room di class
		First(&krs).Error

	if err != nil {
		return nil, err
	}

	return &krs, nil
}

// ListAvailableClasses menampilkan semua kelas yang ditawarkan di semester ini dan belum diambil oleh mahasiswa (Req. 1)
func (r *KRSRepository) ListAvailableClasses(semester string, excludedClassIDs []uuid.UUID) ([]models.Class, error) {
	var classes []models.Class
	query := r.DB.
		Where("semester_penawaran = ?", semester).
		Preload("Course").
		Preload("Dosen")

	if len(excludedClassIDs) > 0 {
		query = query.Where("id NOT IN (?)", excludedClassIDs)
	}

	if err := query.Find(&classes).Error; err != nil {
		return nil, err
	}
	return classes, nil
}

// AddItem menambahkan KRS Item baru (Req. 2)
func (r *KRSRepository) AddItem(krsID uuid.UUID, classID uuid.UUID, krsItemStatusActive string) error {
	item := models.KRSItem{
		KRSID:     krsID,
		ClassID:   classID,
		Status:    krsItemStatusActive,
		CreatedAt: time.Now(),
	}
	return r.DB.Create(&item).Error
}

// RemoveItem menghapus KRS Item (drop/hapus matkul sebelum diverifikasi dosen PA) (Req. 3)
func (r *KRSRepository) RemoveItem(krsID uuid.UUID, classID uuid.UUID, krsItemStatusActive string) error {
	result := r.DB.Where("krs_id = ? AND class_id = ? AND status = ?",
		krsID, classID, krsItemStatusActive).
		Delete(&models.KRSItem{})
		
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound 
	}

	return nil
}

// RequestCancellation mengubah status KRS Item menjadi CANCELLATION_REQUEST (Req. 5)
func (r *KRSRepository) RequestCancellation(krsID uuid.UUID, classID uuid.UUID, krsItemStatusActive string, krsItemStatusCancellationRequest string) error {
	result := r.DB.Model(&models.KRSItem{}).
		Where("krs_id = ? AND class_id = ? AND status = ?", krsID, classID, krsItemStatusActive).
		Updates(map[string]interface{}{
			"status":            krsItemStatusCancellationRequest,
			"diajukan_batal_at": time.Now(),
		})
	
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return errors.New("Matakuliah tidak ditemukan dalam KRS aktif Anda atau sudah diajukan pembatalan.")
	}
	
	return nil
}

// ListMahasiswaByDosenPA mengembalikan seluruh mahasiswa bimbingan dosen PA tertentu
func (r *KRSRepository) ListMahasiswaByDosenPA(dosenID uuid.UUID) ([]models.User, error) {
	var students []models.User
	err := r.DB.Where("role = ? AND dosen_pa_id = ?", "mahasiswa", dosenID).Find(&students).Error
	return students, err
}

// UpdateKRSItemStatus memperbarui status item KRS (digunakan dosen PA saat verifikasi atau pembatalan)
func (r *KRSRepository) UpdateKRSItemStatus(krsID uuid.UUID, classID uuid.UUID, newStatus string, extra map[string]interface{}) error {
	updates := map[string]interface{}{
		"status": newStatus,
	}
	for k, v := range extra {
		updates[k] = v
	}

	result := r.DB.Model(&models.KRSItem{}).
		Where("krs_id = ? AND class_id = ?", krsID, classID).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// UpdateKRSStatus memperbarui status KRS dan verified_at jika diperlukan
func (r *KRSRepository) UpdateKRSStatus(krsID uuid.UUID, newStatus string, verifiedAt *time.Time) error {
	updates := map[string]interface{}{
		"status": newStatus,
	}
	if verifiedAt != nil {
		updates["verified_at"] = *verifiedAt
	}

	return r.DB.Model(&models.KRS{}).
		Where("id = ?", krsID).
		Updates(updates).Error
}

// UpdateKRSItemClass memindahkan item KRS ke class lain (dipakai dosen PA untuk edit)
func (r *KRSRepository) UpdateKRSItemClass(krsID uuid.UUID, oldClassID uuid.UUID, newClassID uuid.UUID, newStatus string) error {
	updates := map[string]interface{}{
		"class_id": newClassID,
	}
	if newStatus != "" {
		updates["status"] = newStatus
	}

	result := r.DB.Model(&models.KRSItem{}).
		Where("krs_id = ? AND class_id = ?", krsID, oldClassID).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
