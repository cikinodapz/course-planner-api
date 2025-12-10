package service

import (
	"course-planner-api/internal/models"
	"course-planner-api/internal/repository"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CourseService interface {
	CreateCourse(kode, nama string, sks int) (*models.Course, error)
	GetCourseByID(id uuid.UUID) (*models.Course, error)
	GetAllCourses() ([]models.Course, error)
	UpdateCourse(id uuid.UUID, kode, nama *string, sks *int) (*models.Course, error)
	DeleteCourse(id uuid.UUID) error
}

type courseService struct {
	repo repository.CourseRepository
}

func NewCourseService(repo repository.CourseRepository) CourseService {
	return &courseService{repo: repo}
}

func (s *courseService) CreateCourse(kode, nama string, sks int) (*models.Course, error) {
	// Check if kode already exists
	existing, err := s.repo.FindByKode(kode)
	if err == nil && existing != nil {
		return nil, errors.New("course dengan kode tersebut sudah ada")
	}

	course := &models.Course{
		Kode: kode,
		Nama: nama,
		SKS:  sks,
	}

	if err := s.repo.Create(course); err != nil {
		return nil, err
	}

	return course, nil
}

func (s *courseService) GetCourseByID(id uuid.UUID) (*models.Course, error) {
	return s.repo.FindByID(id)
}

func (s *courseService) GetAllCourses() ([]models.Course, error) {
	return s.repo.FindAll()
}

func (s *courseService) UpdateCourse(id uuid.UUID, kode, nama *string, sks *int) (*models.Course, error) {
	course, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("course not found")
		}
		return nil, err
	}

	if kode != nil {
		// Check if new kode already exists
		existing, err := s.repo.FindByKode(*kode)
		if err == nil && existing != nil && existing.ID != id {
			return nil, errors.New("course dengan kode tersebut sudah ada")
		}
		course.Kode = *kode
	}

	if nama != nil {
		course.Nama = *nama
	}

	if sks != nil {
		course.SKS = *sks
	}

	if err := s.repo.Update(course); err != nil {
		return nil, err
	}

	return course, nil
}

func (s *courseService) DeleteCourse(id uuid.UUID) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("course not found")
		}
		return err
	}

	return s.repo.Delete(id)
}
