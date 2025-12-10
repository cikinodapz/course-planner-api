package repository

import (
	"course-planner-api/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomRepository interface {
	Create(room *models.Room) error
	FindByID(id uuid.UUID) (*models.Room, error)
	FindAll() ([]models.Room, error)
	Update(room *models.Room) error
	Delete(id uuid.UUID) error
	FindByNama(nama string) (*models.Room, error)
}

type roomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) Create(room *models.Room) error {
	return r.db.Create(room).Error
}

func (r *roomRepository) FindByID(id uuid.UUID) (*models.Room, error) {
	var room models.Room
	if err := r.db.First(&room, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *roomRepository) FindAll() ([]models.Room, error) {
	var rooms []models.Room
	if err := r.db.Find(&rooms).Error; err != nil {
		return nil, err
	}
	return rooms, nil
}

func (r *roomRepository) Update(room *models.Room) error {
	return r.db.Save(room).Error
}

func (r *roomRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Room{}, "id = ?", id).Error
}

func (r *roomRepository) FindByNama(nama string) (*models.Room, error) {
	var room models.Room
	if err := r.db.Where("nama = ?", nama).First(&room).Error; err != nil {
		return nil, err
	}
	return &room, nil
}
