package usecase

import (
	"errors"
	"strings"
	"time"

	"api-on/internal/auth/model"
	"api-on/internal/organization/repository"
	"api-on/internal/auth/validation"
	"api-on/pkg/hash"

	"github.com/google/uuid"
)

type UserUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) *UserUsecase {
	return &UserUsecase{userRepo: userRepo}
}

func (u *UserUsecase) GetByID(id uuid.UUID) (*model.UserResponse, error) {
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return model.ToUserResponse(user), nil
}

func (u *UserUsecase) List() ([]model.UserResponse, error) {
	users, err := u.userRepo.List()
	if err != nil {
		return nil, err
	}

	resp := make([]model.UserResponse, 0, len(users))
	for i := range users {
		resp = append(resp, *model.ToUserResponse(&users[i]))
	}

	return resp, nil
}

func (u *UserUsecase) Create(input model.CreateUserInput) (*model.UserResponse, error) {
	if err := validation.ValidateName(input.Name); err != nil {
		return nil, err
	}
	if err := validation.ValidateEmail(input.Email); err != nil {
		return nil, err
	}
	if err := validation.ValidatePassword(input.Password); err != nil {
		return nil, err
	}
	if err := validation.ValidateTipoPlano(input.TipoPlano); err != nil {
		return nil, err
	}
	if strings.TrimSpace(input.Status) == "" {
		return nil, errors.New("status is required")
	}

	_, err := u.userRepo.FindByEmail(input.Email)
	if err == nil {
		return nil, errors.New("email already registered")
	}

	passwordHash, err := hash.Generate(input.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:           uuid.New(),
		Name:         strings.TrimSpace(input.Name),
		Email:        strings.TrimSpace(strings.ToLower(input.Email)),
		PasswordHash: passwordHash,
		TipoPlano:    input.TipoPlano,
		Status:       input.Status,
		CreatedAt:    time.Now(),
	}

	if err := u.userRepo.Create(user); err != nil {
		return nil, err
	}

	return model.ToUserResponse(user), nil
}

func (u *UserUsecase) Update(id uuid.UUID, input model.UpdateUserInput) (*model.UserResponse, error) {
	if err := validation.ValidateName(input.Name); err != nil {
		return nil, err
	}
	if err := validation.ValidateEmail(input.Email); err != nil {
		return nil, err
	}
	if err := validation.ValidateTipoPlano(input.TipoPlano); err != nil {
		return nil, err
	}
	if strings.TrimSpace(input.Status) == "" {
		return nil, errors.New("status is required")
	}

	user, err := u.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	user.Name = strings.TrimSpace(input.Name)
	user.Email = strings.TrimSpace(strings.ToLower(input.Email))
	user.TipoPlano = input.TipoPlano
	user.Status = input.Status

	if err := u.userRepo.Update(user); err != nil {
		return nil, err
	}

	return model.ToUserResponse(user), nil
}

func (u *UserUsecase) Delete(id uuid.UUID) error {
	return u.userRepo.Delete(id)
}