package user

import (
	"context"
	"strings"
	"time"

	sharederrors "api-on/internal/shared/errors"
	"api-on/internal/shared/permissions"
	"api-on/internal/shared/security"
	sharedvalidator "api-on/internal/shared/validator"
	"api-on/internal/subscription"
	"api-on/pkg/hash"

	"github.com/google/uuid"
)

type Usecase struct {
	repo             Repository
	subscriptionRepo subscription.Repository
}

func NewUsecase(repo Repository, subscriptionRepo subscription.Repository) *Usecase {
	return &Usecase{
		repo:             repo,
		subscriptionRepo: subscriptionRepo,
	}
}

// GetMe retorna a conta autenticada apenas pela identidade resolvida do JWT.
// A rota nao aceita userID externo para evitar IDOR/BOLA entre contas e clinicas.
func (u *Usecase) GetMe(ctx context.Context, actor security.Identity) (*Response, error) {
	item, err := u.repo.GetByIDAndTenant(ctx, actor.TenantID, actor.UserID)
	if err != nil {
		return nil, err
	}
	return ToResponse(item), nil
}

func (u *Usecase) List(ctx context.Context, actor security.Identity, input ListInput) ([]Response, ListMeta, error) {
	if !canManageUsers(actor) {
		return nil, ListMeta{}, sharederrors.Forbidden("only owner or admin can list clinic users")
	}

	items, total, err := u.repo.ListByTenant(ctx, actor.TenantID, input)
	if err != nil {
		return nil, ListMeta{}, err
	}

	resp := make([]Response, 0, len(items))
	for i := range items {
		resp = append(resp, *ToResponse(&items[i]))
	}

	return resp, BuildListMeta(input, total), nil
}

func (u *Usecase) Create(ctx context.Context, actor security.Identity, input CreateInput) (*Response, error) {
	if !canManageUsers(actor) {
		return nil, sharederrors.Forbidden("only owner or admin can create users")
	}

	if err := sharedvalidator.ValidateName(input.Name); err != nil {
		return nil, err
	}
	if err := sharedvalidator.ValidateEmail(input.Email); err != nil {
		return nil, err
	}
	if err := sharedvalidator.ValidatePassword(input.Password); err != nil {
		return nil, err
	}
	if err := ValidateRole(input.Role); err != nil {
		return nil, err
	}
	if err := ValidateStatus(input.Status); err != nil {
		return nil, err
	}

	if _, err := u.repo.FindByEmail(ctx, input.Email); err == nil {
		return nil, sharederrors.Conflict("USER_EMAIL_ALREADY_EXISTS", "email already registered")
	} else if appErr := sharederrors.AsAppError(err); appErr != nil && appErr.Code != "USER_NOT_FOUND" {
		return nil, err
	}

	if err := u.ensurePlanCapacity(ctx, actor.TenantID); err != nil {
		return nil, err
	}

	passwordHash, err := hash.Generate(input.Password)
	if err != nil {
		return nil, sharederrors.Internal("could not hash password")
	}

	now := time.Now()
	item := &User{
		ID:           uuid.New(),
		TenantID:     actor.TenantID,
		Name:         strings.TrimSpace(input.Name),
		Email:        sharedvalidator.NormalizeEmail(input.Email),
		PasswordHash: passwordHash,
		Role:         strings.TrimSpace(input.Role),
		Status:       strings.TrimSpace(input.Status),
		Permissions:  permissions.DefaultForRole(strings.TrimSpace(input.Role)),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := u.repo.Create(ctx, item); err != nil {
		return nil, err
	}

	return ToResponse(item), nil
}

func (u *Usecase) Update(ctx context.Context, actor security.Identity, userID uuid.UUID, input UpdateInput) (*Response, error) {
	if !canManageUsers(actor) {
		return nil, sharederrors.Forbidden("only owner or admin can update users")
	}

	if err := sharedvalidator.ValidateName(input.Name); err != nil {
		return nil, err
	}
	if err := sharedvalidator.ValidateEmail(input.Email); err != nil {
		return nil, err
	}
	if err := ValidateRole(input.Role); err != nil {
		return nil, err
	}
	if err := ValidateStatus(input.Status); err != nil {
		return nil, err
	}

	item, err := u.repo.GetByIDAndTenant(ctx, actor.TenantID, userID)
	if err != nil {
		return nil, err
	}

	if item.Role == RoleOwner && actor.Role != RoleOwner {
		return nil, sharederrors.Forbidden("only owner can change owner account")
	}

	if actor.UserID == item.ID && (input.Role != item.Role || input.Status != item.Status) {
		return nil, sharederrors.Invalid("SELF_UPDATE_NOT_ALLOWED", "you cannot change your own role or status")
	}

	item.Name = strings.TrimSpace(input.Name)
	item.Email = sharedvalidator.NormalizeEmail(input.Email)
	if item.Role != strings.TrimSpace(input.Role) {
		item.Permissions = permissions.DefaultForRole(strings.TrimSpace(input.Role))
	}
	item.Role = strings.TrimSpace(input.Role)
	item.Status = strings.TrimSpace(input.Status)
	item.UpdatedAt = time.Now()

	if err := u.repo.Update(ctx, item); err != nil {
		return nil, err
	}

	return ToResponse(item), nil
}

func (u *Usecase) Deactivate(ctx context.Context, actor security.Identity, userID uuid.UUID) error {
	if !canManageUsers(actor) {
		return sharederrors.Forbidden("only owner or admin can deactivate users")
	}

	item, err := u.repo.GetByIDAndTenant(ctx, actor.TenantID, userID)
	if err != nil {
		return err
	}

	if actor.UserID == item.ID {
		return sharederrors.Invalid("SELF_DEACTIVATION_NOT_ALLOWED", "you cannot deactivate your own account")
	}

	if item.Role == RoleOwner && actor.Role != RoleOwner {
		return sharederrors.Forbidden("only owner can deactivate owner account")
	}

	return u.repo.Deactivate(ctx, actor.TenantID, userID)
}

func (u *Usecase) ensurePlanCapacity(ctx context.Context, tenantID uuid.UUID) error {
	_, plan, err := u.subscriptionRepo.GetByTenantID(ctx, tenantID)
	if err != nil {
		return err
	}

	total, err := u.repo.CountActiveByTenant(ctx, tenantID)
	if err != nil {
		return err
	}

	if total >= plan.MaxProfessionals {
		return sharederrors.Conflict("PLAN_LIMIT_REACHED", "current plan has reached the maximum number of active professionals")
	}

	return nil
}

func canManageUsers(actor security.Identity) bool {
	return actor.IsInternal() &&
		actor.HasRole(RoleOwner, RoleAdmin) &&
		actor.Permissions.CanViewAllUsers()
}
