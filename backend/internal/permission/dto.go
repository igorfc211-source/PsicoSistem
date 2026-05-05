package permission

import (
	"api-on/internal/shared/permissions"

	"github.com/google/uuid"
)

type UpdateInput struct {
	Permissions permissions.AccountPermissions `json:"permissions"`
}

type Response struct {
	UserID         uuid.UUID                      `json:"user_id"`
	Role           string                         `json:"role"`
	Permissions    permissions.AccountPermissions `json:"permissions"`
	CanManageUsers bool                           `json:"can_manage_users"`
}
