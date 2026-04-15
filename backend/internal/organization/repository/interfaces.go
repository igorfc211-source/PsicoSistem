package repository

import (
	"api-on/internal/auth/model"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByID(id uuid.UUID) (*model.User, error)
	List() ([]model.User, error)
	Update(user *model.User) error
	Delete(id uuid.UUID) error
}