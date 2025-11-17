package repository

import (
	"course-planner-api/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClassRepository interface {
	Create(class *models.Class) error
	FindByID(id uuid.UUID) (*models.Class, error)
	FindAll() ([]models.Class, error)
	Update(class *models.Class) error
	Delete(id uuid.UUID) error
	HasTimeConflict(roomID uuid.UUID, hari string, start, end time.Time, excludeID *uuid.UUID) (bool, error)
}

type classRepository struct {
	db *gorm.DB
}

func NewClassRepository(db *gorm.DB) ClassRepository {
	return &classRepository{db: db}
}

func (r *classRepository) Create(class *models.Class) error {
	return r.db.Create(class).Error
}

func (r *classRepository) FindByID(id uuid.UUID) (*models.Class, error) {
	var class models.Class
	if err := r.db.Preload("Course").Preload("Dosen").Preload("Room").First(&class, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &class, nil
}

func (r *classRepository) FindAll() ([]models.Class, error) {
	var classes []models.Class
	if err := r.db.Preload("Course").Preload("Dosen").Preload("Room").Find(&classes).Error; err != nil {
		return nil, err
	}
	return classes, nil
}

func (r *classRepository) HasTimeConflict(roomID uuid.UUID, hari string, start, end time.Time, excludeID *uuid.UUID) (bool, error) {
	var count int64
	query := r.db.Model(&models.Class{}).
		Where("room_id = ? AND hari = ? AND NOT (jam_selesai <= ? OR jam_mulai >= ?)", roomID, hari, start, end)

	if excludeID != nil {
		query = query.Where("id <> ?", *excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}


func (r *classRepository) Update(class *models.Class) error {
	return r.db.Save(class).Error
}

func (r *classRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Class{}, "id = ?", id).Error
}
