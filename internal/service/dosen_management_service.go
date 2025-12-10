package service

import (
	"course-planner-api/internal/models"
	"course-planner-api/internal/repository"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DosenManagementService interface {
	GetDosenByID(id uuid.UUID) (*models.User, error)
	GetAllDosen() ([]models.User, error)
	UpdateDosen(id uuid.UUID, name, nidn *string) (*models.User, error)
}

type dosenManagementService struct {
	repo repository.DosenRepository
}

func NewDosenManagementService(repo repository.DosenRepository) DosenManagementService {
	return &dosenManagementService{repo: repo}
}

func (s *dosenManagementService) GetDosenByID(id uuid.UUID) (*models.User, error) {
	return s.repo.FindByID(id)
}

func (s *dosenManagementService) GetAllDosen() ([]models.User, error) {
	return s.repo.FindAllDosen()
}

func (s *dosenManagementService) UpdateDosen(id uuid.UUID, name, nidn *string) (*models.User, error) {
	dosen, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("dosen not found")
		}
		return nil, err
	}

	if name != nil {
		dosen.Name = *name
	}

	if nidn != nil {
		// Check if new nidn already exists
		existing, err := s.repo.FindByNIDN(*nidn)
		if err == nil && existing != nil && existing.ID != id {
			return nil, errors.New("dosen dengan NIDN tersebut sudah ada")
		}
		dosen.NIDN = *nidn
	}

	if err := s.repo.Update(dosen); err != nil {
		return nil, err
	}

	return dosen, nil
}
