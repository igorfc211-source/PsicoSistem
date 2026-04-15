package model

import (
	"time"

	"github.com/google/uuid"
)

type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	TipoPlano string `json:"tipo_plano"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	
}

type CreateUserInput struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	TipoPlano string `json:"tipo_plano"`
	Status    string `json:"status"`
}

type UpdateUserInput struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	TipoPlano string `json:"tipo_plano"`
	Status    string `json:"status"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"userId"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	TipoPlano string    `json:"tipo_plano"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func ToUserResponse(user *User) *UserResponse {
	if user == nil {
		return nil
	}

	return &UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		TipoPlano: user.TipoPlano,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
	}
}