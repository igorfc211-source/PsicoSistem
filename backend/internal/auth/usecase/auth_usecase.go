package usecase

import (
	"errors"
	"strings"
	"time"

	"api-on/internal/auth/model"
	"api-on/internal/organization/repository"
	"api-on/internal/auth/validation"
	"api-on/pkg/hash"
	jwtpkg "api-on/pkg/jwt"

	"github.com/google/uuid"
)

type AuthUsecase struct {
	userRepo repository.UserRepository
	jwtSvc   *jwtpkg.JWTService
}

func NewAuthUsecase(userRepo repository.UserRepository, jwtSvc *jwtpkg.JWTService) *AuthUsecase {
	return &AuthUsecase{
		userRepo: userRepo,
		jwtSvc:   jwtSvc,
	}
}

func (u *AuthUsecase) Register(input model.RegisterInput) (*model.UserResponse, string, error) {
	if err := validation.ValidateName(input.Name); err != nil {
		return nil, "", err
	}
	if err := validation.ValidateEmail(input.Email); err != nil {
		return nil, "", err
	}
	if err := validation.ValidatePassword(input.Password); err != nil {
		return nil, "", err
	}
	if err := validation.ValidateTipoPlano(input.TipoPlano); err != nil {
		return nil, "", err
	}

	_, err := u.userRepo.FindByEmail(input.Email)
	if err == nil {
		return nil, "", errors.New("email already registered")
	}

	passwordHash, err := hash.Generate(input.Password)
	if err != nil {
		return nil, "", err
	}

	user := &model.User{
		ID:           uuid.New(),
		Name:         strings.TrimSpace(input.Name),
		Email:        strings.TrimSpace(strings.ToLower(input.Email)),
		PasswordHash: passwordHash,
		TipoPlano:    input.TipoPlano,
		Status:       "active",
		CreatedAt:    time.Now(),
	}

	if err := u.userRepo.Create(user); err != nil {
		return nil, "", err
	}

	token, err := u.jwtSvc.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, "", err
	}

	return model.ToUserResponse(user), token, nil
}

func (u *AuthUsecase) Login(input model.LoginInput) (*model.UserResponse, string, error) {
	if err := validation.ValidateEmail(input.Email); err != nil {
		return nil, "", err
	}
	if strings.TrimSpace(input.Password) == "" {
		return nil, "", errors.New("password is required")
	}

	user, err := u.userRepo.FindByEmail(input.Email)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if err := hash.Compare(input.Password, user.PasswordHash); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := u.jwtSvc.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, "", err
	}

	return model.ToUserResponse(user), token, nil
}