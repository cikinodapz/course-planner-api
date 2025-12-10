package repository

import (
	"course-planner-api/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DosenRepository interface {
	FindByID(id uuid.UUID) (*models.User, error)
	FindAllDosen() ([]models.User, error)
	Update(user *models.User) error
	FindByNIDN(nidn string) (*models.User, error)
}

type dosenRepository struct {
	db *gorm.DB
}

func NewDosenRepository(db *gorm.DB) DosenRepository {
	return &dosenRepository{db: db}
}

func (r *dosenRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ? AND role = ?", id, "dosen").First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *dosenRepository) FindAllDosen() ([]models.User, error) {
	var users []models.User
	if err := r.db.Where("role = ?", "dosen").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *dosenRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *dosenRepository) FindByNIDN(nidn string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("nidn = ? AND role = ?", nidn, "dosen").First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
