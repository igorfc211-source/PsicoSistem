package security

import (
	"api-on/internal/shared/permissions"

	"github.com/google/uuid"
)

const (
	UserTypeInternal = "internal"
	UserTypeGuardian = "guardian"
)

// Identity representa o ator autenticado da requisição.
type Identity struct {
	UserID      uuid.UUID
	TenantID    uuid.UUID
	Role        string
	Email       string
	Type        string
	Permissions permissions.AccountPermissions
}

func (i Identity) HasRole(roles ...string) bool {
	for _, role := range roles {
		if i.Role == role {
			return true
		}
	}

	return false
}

func (i Identity) IsInternal() bool {
	return i.Type == UserTypeInternal
}
