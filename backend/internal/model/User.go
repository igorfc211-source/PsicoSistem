package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID    `json:"userId"`
	Name         string       `json:"name"`
	Email        string       `json:"email"`
	PasswordHash string       `json:"-"`
	TipoPlano    string       `json:"tipo_plano"`
	Status       string       `json:"status"`
	CreatedAt    time.Time    `json:"created_at"`
}






