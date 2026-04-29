package permission

import (
	"context"

	sharederrors "api-on/internal/shared/errors"
	sharedpermissions "api-on/internal/shared/permissions"
	"api-on/internal/shared/security"
	"api-on/internal/user"

	"github.com/google/uuid"
)

type Usecase struct {
	userRepo user.Repository
}

func NewUsecase(userRepo user.Repository) *Usecase {
	return &Usecase{userRepo: userRepo}
}

func (u *Usecase) GetMe(ctx context.Context, actor security.Identity) (*Response, error) {
	item, err := u.userRepo.GetByIDAndTenant(ctx, actor.TenantID, actor.UserID)
	if err != nil {
		return nil, err
	}

	return toResponse(item), nil
}

func (u *Usecase) GetByUserID(ctx context.Context, actor security.Identity, userID uuid.UUID) (*Response, error) {
	if !canManagePermissions(actor) && actor.UserID != userID {
		return nil, sharederrors.Forbidden("you do not have permission to inspect this account permissions")
	}

	item, err := u.userRepo.GetByIDAndTenant(ctx, actor.TenantID, userID)
	if err != nil {
		return nil, err
	}

	return toResponse(item), nil
}

func (u *Usecase) UpdateByUserID(
	ctx context.Context,
	actor security.Identity,
	userID uuid.UUID,
	input UpdateInput,
) (*Response, error) {
	if !canManagePermissions(actor) {
		return nil, sharederrors.Forbidden("only owner or admin can change account permissions")
	}

	item, err := u.userRepo.GetByIDAndTenant(ctx, actor.TenantID, userID)
	if err != nil {
		return nil, err
	}

	if item.Role == user.RoleOwner && actor.Role != user.RoleOwner {
		return nil, sharederrors.Forbidden("only owner can change owner permissions")
	}

	normalized, err := sharedpermissions.Normalize(item.Role, &input.Permissions)
	if err != nil {
		return nil, err
	}

	item.Permissions = normalized
	if err := u.userRepo.Update(ctx, item); err != nil {
		return nil, err
	}

	return toResponse(item), nil
}

func toResponse(item *user.User) *Response {
	return &Response{
		UserID:         item.ID,
		Role:           item.Role,
		Permissions:    item.Permissions,
		CanManageUsers: item.Permissions.CanViewAllUsers() && (item.Role == user.RoleOwner || item.Role == user.RoleAdmin),
	}
}

func canManagePermissions(actor security.Identity) bool {
	return actor.IsInternal() &&
		actor.Permissions.CanViewAllUsers() &&
		actor.HasRole(user.RoleOwner, user.RoleAdmin)
}
