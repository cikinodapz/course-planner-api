package repository

import (
	"course-planner-api/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CourseRepository interface {
	Create(course *models.Course) error
	FindByID(id uuid.UUID) (*models.Course, error)
	FindAll() ([]models.Course, error)
	Update(course *models.Course) error
	Delete(id uuid.UUID) error
	FindByKode(kode string) (*models.Course, error)
}

type courseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) CourseRepository {
	return &courseRepository{db: db}
}

func (r *courseRepository) Create(course *models.Course) error {
	return r.db.Create(course).Error
}

func (r *courseRepository) FindByID(id uuid.UUID) (*models.Course, error) {
	var course models.Course
	if err := r.db.First(&course, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &course, nil
}

func (r *courseRepository) FindAll() ([]models.Course, error) {
	var courses []models.Course
	if err := r.db.Find(&courses).Error; err != nil {
		return nil, err
	}
	return courses, nil
}

func (r *courseRepository) Update(course *models.Course) error {
	return r.db.Save(course).Error
}

func (r *courseRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Course{}, "id = ?", id).Error
}

func (r *courseRepository) FindByKode(kode string) (*models.Course, error) {
	var course models.Course
	if err := r.db.Where("kode = ?", kode).First(&course).Error; err != nil {
		return nil, err
	}
	return &course, nil
}
