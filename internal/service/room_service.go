package service

import (
	"course-planner-api/internal/models"
	"course-planner-api/internal/repository"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomService interface {
	CreateRoom(nama string) (*models.Room, error)
	GetRoomByID(id uuid.UUID) (*models.Room, error)
	GetAllRooms() ([]models.Room, error)
	UpdateRoom(id uuid.UUID, nama *string) (*models.Room, error)
	DeleteRoom(id uuid.UUID) error
}

type roomService struct {
	repo repository.RoomRepository
}

func NewRoomService(repo repository.RoomRepository) RoomService {
	return &roomService{repo: repo}
}

func (s *roomService) CreateRoom(nama string) (*models.Room, error) {
	// Check if nama already exists
	existing, err := s.repo.FindByNama(nama)
	if err == nil && existing != nil {
		return nil, errors.New("room dengan nama tersebut sudah ada")
	}

	room := &models.Room{
		Nama: nama,
	}

	if err := s.repo.Create(room); err != nil {
		return nil, err
	}

	return room, nil
}

func (s *roomService) GetRoomByID(id uuid.UUID) (*models.Room, error) {
	return s.repo.FindByID(id)
}

func (s *roomService) GetAllRooms() ([]models.Room, error) {
	return s.repo.FindAll()
}

func (s *roomService) UpdateRoom(id uuid.UUID, nama *string) (*models.Room, error) {
	room, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("room not found")
		}
		return nil, err
	}

	if nama != nil {
		// Check if new nama already exists
		existing, err := s.repo.FindByNama(*nama)
		if err == nil && existing != nil && existing.ID != id {
			return nil, errors.New("room dengan nama tersebut sudah ada")
		}
		room.Nama = *nama
	}

	if err := s.repo.Update(room); err != nil {
		return nil, err
	}

	return room, nil
}

func (s *roomService) DeleteRoom(id uuid.UUID) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("room not found")
		}
		return err
	}

	return s.repo.Delete(id)
}
