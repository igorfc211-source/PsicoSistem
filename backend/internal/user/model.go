package user

import (
	"time"

	"github.com/google/uuid"
)

const (
	RoleOwner        = "owner"
	RoleAdmin        = "admin"
	RoleCoordinator  = "coordinator"
	RoleProfessional = "professional"
	RoleFinancial    = "financial"

	StatusActive   = "active"
	StatusInactive = "inactive"
)

type User struct {
	ID           uuid.UUID  `json:"user_id"`
	TenantID     uuid.UUID  `json:"tenant_id"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	Role         string     `json:"role"`
	Status       string     `json:"status"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
