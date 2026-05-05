package user

import (
	"context"

	sharederrors "api-on/internal/shared/errors"
	"api-on/internal/shared/security"

	"github.com/google/uuid"
)

type IdentityResolver struct {
	repo Repository
}

func NewIdentityResolver(repo Repository) *IdentityResolver {
	return &IdentityResolver{repo: repo}
}

func (r *IdentityResolver) ResolveInternalIdentity(
	ctx context.Context,
	tenantID uuid.UUID,
	userID uuid.UUID,
) (security.Identity, error) {
	item, err := r.repo.GetByIDAndTenant(ctx, tenantID, userID)
	if err != nil {
		return security.Identity{}, sharederrors.Unauthorized("user session is no longer valid")
	}

	if item.Status != StatusActive {
		return security.Identity{}, sharederrors.Forbidden("inactive users cannot access protected routes")
	}

	return security.Identity{
		UserID:      item.ID,
		TenantID:    item.TenantID,
		Role:        item.Role,
		Email:       item.Email,
		Type:        security.UserTypeInternal,
		Permissions: item.Permissions,
	}, nil
}
