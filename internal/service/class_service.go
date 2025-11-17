package service

import (
	"course-planner-api/internal/models"
	"course-planner-api/internal/repository"
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrClassTimeConflict = errors.New("room is already used at the given time")

type CreateClassInput struct {
	CourseID          uuid.UUID
	DosenID           uuid.UUID
	NamaKelas         string
	Hari              string
	JamMulai          time.Time
	JamSelesai        time.Time
	RoomID            uuid.UUID
	Kuota             int
	SemesterPenawaran string
}

type UpdateClassInput struct {
	CourseID          *uuid.UUID
	DosenID           *uuid.UUID
	NamaKelas         *string
	Hari              *string
	JamMulai          *time.Time
	JamSelesai        *time.Time
	RoomID            *uuid.UUID
	Kuota             *int
	SemesterPenawaran *string
}

type ClassService interface {
	CreateClass(input CreateClassInput) (*models.Class, error)
	GetClass(id uuid.UUID) (*models.Class, error)
	ListClasses() ([]models.Class, error)
	UpdateClass(id uuid.UUID, input UpdateClassInput) (*models.Class, error)
	DeleteClass(id uuid.UUID) error
}

type classService struct {
	classRepo repository.ClassRepository
}

func NewClassService(classRepo repository.ClassRepository) ClassService {
	return &classService{classRepo: classRepo}
}

func (s *classService) CreateClass(input CreateClassInput) (*models.Class, error) {
	conflict, err := s.classRepo.HasTimeConflict(input.RoomID, input.Hari, input.JamMulai, input.JamSelesai, nil)
	if err != nil {
		return nil, err
	}
	if conflict {
		return nil, ErrClassTimeConflict
	}

	class := &models.Class{
		CourseID:          input.CourseID,
		DosenID:           input.DosenID,
		NamaKelas:         input.NamaKelas,
		Hari:              input.Hari,
		JamMulai:          input.JamMulai,
		JamSelesai:        input.JamSelesai,
		RoomID:            input.RoomID,
		Kuota:             input.Kuota,
		SemesterPenawaran: input.SemesterPenawaran,
	}

	if err := s.classRepo.Create(class); err != nil {
		return nil, err
	}

	return s.classRepo.FindByID(class.ID)
}

func (s *classService) GetClass(id uuid.UUID) (*models.Class, error) {
	return s.classRepo.FindByID(id)
}

func (s *classService) ListClasses() ([]models.Class, error) {
	return s.classRepo.FindAll()
}

func (s *classService) UpdateClass(id uuid.UUID, input UpdateClassInput) (*models.Class, error) {
	class, err := s.classRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	finalRoomID := class.RoomID
	finalHari := class.Hari
	finalJamMulai := class.JamMulai
	finalJamSelesai := class.JamSelesai

	if input.CourseID != nil {
		class.CourseID = *input.CourseID
	}
	if input.DosenID != nil {
		class.DosenID = *input.DosenID
	}
	if input.NamaKelas != nil {
		class.NamaKelas = *input.NamaKelas
	}
	if input.Hari != nil {
		class.Hari = *input.Hari
		finalHari = *input.Hari
	}
	if input.JamMulai != nil {
		class.JamMulai = *input.JamMulai
		finalJamMulai = *input.JamMulai
	}
	if input.JamSelesai != nil {
		class.JamSelesai = *input.JamSelesai
		finalJamSelesai = *input.JamSelesai
	}
	if input.RoomID != nil {
		class.RoomID = *input.RoomID
		finalRoomID = *input.RoomID
	}
	if input.Kuota != nil {
		class.Kuota = *input.Kuota
	}
	if input.SemesterPenawaran != nil {
		class.SemesterPenawaran = *input.SemesterPenawaran
	}

	conflict, err := s.classRepo.HasTimeConflict(finalRoomID, finalHari, finalJamMulai, finalJamSelesai, &id)
	if err != nil {
		return nil, err
	}
	if conflict {
		return nil, ErrClassTimeConflict
	}

	if err := s.classRepo.Update(class); err != nil {
		return nil, err
	}

	return s.classRepo.FindByID(id)
}

func (s *classService) DeleteClass(id uuid.UUID) error {
	return s.classRepo.Delete(id)
}
