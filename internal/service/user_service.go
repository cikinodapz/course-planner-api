package service

import (
	"course-planner-api/internal/models"
	"course-planner-api/internal/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo}
}

func (s *UserService) ListStudents() ([]models.User, error) {
	return s.repo.FindAllStudents()
}