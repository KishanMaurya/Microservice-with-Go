package services

import (
	"time"

	"go-backend/internal/models"
	"go-backend/internal/repositories"

	"github.com/google/uuid"
)

type UserService struct {
	Repo *repositories.UserRepository
}

// CreateUser creates a new user with the given request.
// It returns the created user and an error if any occurs.
func (s *UserService) CreateUser(req models.CreateUserRequest) (*models.User, error) {

	user := models.User{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := s.Repo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}