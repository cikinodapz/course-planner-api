package service

import (
	"course-planner-api/internal/models"
	"course-planner-api/internal/repository"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	KRS_STATUS_DRAFT    = "DRAFT"
	KRS_STATUS_VERIFIED = "VERIFIED"

	KRS_ITEM_STATUS_ACTIVE               = "ACTIVE"
	KRS_ITEM_STATUS_CANCELLATION_REQUEST = "CANCELLATION_REQUEST"
	KRS_ITEM_STATUS_CANCELLED            = "CANCELLED"
)

type KRSService struct {
	Repo *repository.KRSRepository
}

func NewKRSService(repo *repository.KRSRepository) *KRSService {
	return &KRSService{Repo: repo}
}

// GetCurrentSemester otomatis berdasarkan bulan
func GetCurrentSemester() string {
	now := time.Now()
	month := now.Month()
	if month >= time.January && month <= time.June {
		return "ganjil"
	} else {
		return "genap"
	}
}

// CheckScheduleConflict memeriksa apakah ada jadwal bentrok antara kelas baru dan yang sudah ada
func CheckScheduleConflict(selected []models.Class, existing []models.KRSItem) error {
	for _, sel := range selected {
		for _, item := range existing {
			if item.Status == KRS_ITEM_STATUS_CANCELLED {
				continue
			}
			existClass := item.Class
			if sel.Hari != existClass.Hari {
				continue
			}
			// Cek tumpang tindih jam
			if sel.JamMulai.Before(existClass.JamSelesai) && sel.JamSelesai.After(existClass.JamMulai) {
				return fmt.Errorf("jadwal bentrok: %s dengan %s", sel.NamaKelas, existClass.NamaKelas)
			}
		}
	}
	return nil
}

// ListAvailableClasses: Menampilkan matkul yang tersedia (Req. 1)
func (s *KRSService) ListAvailableClasses(mahasiswaID uuid.UUID) ([]models.Class, error) {
	semester := GetCurrentSemester()
	krs, err := s.Repo.GetOrCreateKRS(mahasiswaID, semester, KRS_STATUS_DRAFT, KRS_STATUS_VERIFIED)
	if err != nil {
		return nil, err
	}

	var excludedClassIDs []uuid.UUID
	for _, item := range krs.Items {
		if item.Status != KRS_ITEM_STATUS_CANCELLED {
			excludedClassIDs = append(excludedClassIDs, item.ClassID)
		}
	}

	return s.Repo.ListAvailableClasses(semester, excludedClassIDs)
}

// TakeClass: Mahasiswa mengambil 1 atau banyak matkul sekaligus
func (s *KRSService) TakeClass(mahasiswaID uuid.UUID, classIDs []uuid.UUID) error {
	semester := GetCurrentSemester()

	krs, err := s.Repo.GetOrCreateKRS(mahasiswaID, semester, KRS_STATUS_DRAFT, KRS_STATUS_VERIFIED)
	if err != nil {
		return err
	}

	if krs.Status == KRS_STATUS_VERIFIED {
		return errors.New("KRS sudah diverifikasi oleh Dosen PA, tidak dapat menambah matakuliah.")
	}

	classes, err := s.Repo.GetClassesByIDs(classIDs)
	if err != nil {
		return err
	}
    
	seenCourseIDs := make(map[uuid.UUID]bool)
	for _, c := range classes {
		if _, ok := seenCourseIDs[c.CourseID]; ok {
			return fmt.Errorf("Anda mencoba mengambil mata kuliah yang sama (Kode: %s) lebih dari sekali dalam satu permintaan.", c.Course.Kode)
		}
		seenCourseIDs[c.CourseID] = true
	}

	if err := CheckScheduleConflict(classes, krs.Items); err != nil {
		return err
	}

	for _, c := range classes {
		for _, item := range krs.Items {
			if item.Status == KRS_ITEM_STATUS_CANCELLED {
				continue
			}

			if item.ClassID == c.ID {
				return errors.New("Kelas matakuliah ini sudah Anda ambil.")
			}

			if item.Class.CourseID == c.CourseID {
				return errors.New("Anda sudah mengambil mata kuliah dengan kode yang sama di kelas lain.")
			}
		}
	}

	return s.Repo.AddItemsBatch(krs.ID, classIDs, KRS_ITEM_STATUS_ACTIVE)
}

// DropClass: Mahasiswa menghapus matkul (Req. 3)
func (s *KRSService) DropClass(mahasiswaID uuid.UUID, classID uuid.UUID) error {
	semester := GetCurrentSemester()
	krs, err := s.Repo.GetOrCreateKRS(mahasiswaID, semester, KRS_STATUS_DRAFT, KRS_STATUS_VERIFIED)
	if err != nil {
		return err
	}

	if krs.Status == KRS_STATUS_VERIFIED {
		return errors.New("KRS sudah diverifikasi oleh Dosen PA. Untuk menghapus matakuliah, silakan ajukan pembatalan.")
	}

	return s.Repo.RemoveItem(krs.ID, classID, KRS_ITEM_STATUS_ACTIVE)
}

// GetTakenClasses: Mahasiswa melihat matkul yang sudah diambil (Req. 4)
func (s *KRSService) GetTakenClasses(mahasiswaID uuid.UUID) (*models.KRS, error) {
	semester := GetCurrentSemester()
	return s.Repo.GetKRSByMahasiswaID(mahasiswaID, semester)
}


// RequestClassCancellation: Mahasiswa mengajukan pembatalan matkul (Req. 5)
func (s *KRSService) RequestClassCancellation(mahasiswaID uuid.UUID, classID uuid.UUID) error {
	semester := GetCurrentSemester()
	krs, err := s.Repo.GetOrCreateKRS(mahasiswaID, semester, KRS_STATUS_DRAFT, KRS_STATUS_VERIFIED)
	if err != nil {
		return err
	}

	if krs.Status != KRS_STATUS_VERIFIED {
		return errors.New("KRS belum diverifikasi oleh Dosen PA. Silakan drop matakuliah jika belum diverifikasi.")
	}

	foundItem := false
	for _, item := range krs.Items {
		if item.ClassID == classID {
			foundItem = true
			if item.Status == KRS_ITEM_STATUS_CANCELLATION_REQUEST {
				return errors.New("Pengajuan pembatalan untuk matakuliah ini sudah ada.")
			}
			if item.Status == KRS_ITEM_STATUS_CANCELLED {
				return errors.New("Matakuliah ini sudah dibatalkan sebelumnya.")
			}
			break
		}
	}

	if !foundItem {
		return errors.New(fmt.Sprintf("Matakuliah dengan ID %s tidak ditemukan dalam KRS Anda.", classID))
	}

	return s.Repo.RequestCancellation(krs.ID, classID, KRS_ITEM_STATUS_ACTIVE, KRS_ITEM_STATUS_CANCELLATION_REQUEST)
}

// GET mahasiswa KRS
func (s *KRSService) GetStudentClasses(MahasiswaID string) ([]models.Class, error) {
    return s.Repo.GetStudentClasses(MahasiswaID)
}
//DELETE class by dosen pa
func (s *KRSService) DropClassByDosen(mahasiswaId string, classId string) error {
    return s.Repo.DropClassByDosen(mahasiswaId, classId)
}


// PATCH Verify KRS Item oleh dosen PA //
func (s *KRSService) VerifyKRSItem(
	itemId uuid.UUID,
	mahasiswaId uuid.UUID,
	dosenId uuid.UUID,
	status string,
	catatan string,
) error {

	// 1. Ambil item + relasi
	item, err := s.Repo.GetKRSItemWithRelations(itemId)
	if err != nil {
		return errors.New("krs item not found")
	}

	// 2. Validasi item milik mahasiswa
	if item.KRS.MahasiswaID != mahasiswaId {
		return errors.New("item does not belong to this student")
	}

	// 3. Validasi dosen PA
	if item.KRS.Mahasiswa.DosenPAID == nil ||
		item.KRS.Mahasiswa.DosenPAID.String() != dosenId.String() {
		return errors.New("forbidden: you are not this student's advisor")
	}

	// 4. Update item
	item.Status = status
	if err := s.Repo.UpdateKRSItem(item); err != nil {
		return errors.New("failed to update item")
	}

	// 5. Jika approved -> update KRS
	if status == "approved" {
		now := time.Now()
		item.KRS.Status = "approved"
		item.KRS.VerifiedAt = &now
		item.KRS.CatatanDosen = catatan
		return s.Repo.UpdateKRS(&item.KRS)
	}

	// rejected -> hanya item yang update
	return nil
}
