package guardian

import (
	"context"
	"strings"
	"time"

	sharederrors "api-on/internal/shared/errors"
	"api-on/internal/shared/permissions"
	"api-on/internal/shared/security"
	sharedvalidator "api-on/internal/shared/validator"

	"github.com/google/uuid"
)

type Usecase struct {
	repo Repository
}

func NewUsecase(repo Repository) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) List(ctx context.Context, actor security.Identity, input ListInput) ([]Response, ListMeta, error) {
	if !canAccessGuardians(actor) {
		return nil, ListMeta{}, sharederrors.Forbidden("you do not have permission to access guardians")
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

func (u *Usecase) Get(ctx context.Context, actor security.Identity, guardianID uuid.UUID) (*Response, error) {
	if !canAccessGuardians(actor) {
		return nil, sharederrors.Forbidden("you do not have permission to access guardians")
	}

	item, err := u.repo.GetByIDAndTenant(ctx, actor.TenantID, guardianID)
	if err != nil {
		return nil, err
	}

	return ToResponse(item), nil
}

func (u *Usecase) Create(ctx context.Context, actor security.Identity, input CreateInput) (*Response, error) {
	if !canAccessGuardians(actor) {
		return nil, sharederrors.Forbidden("you do not have permission to create guardians")
	}

	normalized, err := normalizeInput(input.Name, input.Phone, input.Address, input.CPF)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	item := &Guardian{
		ID:        uuid.New(),
		TenantID:  actor.TenantID,
		Name:      normalized.name,
		Phone:     normalized.phone,
		Address:   normalized.address,
		CPF:       normalized.cpf,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := u.repo.Create(ctx, item); err != nil {
		return nil, err
	}

	return ToResponse(item), nil
}

func (u *Usecase) Update(ctx context.Context, actor security.Identity, guardianID uuid.UUID, input UpdateInput) (*Response, error) {
	if !canAccessGuardians(actor) {
		return nil, sharederrors.Forbidden("you do not have permission to update guardians")
	}

	normalized, err := normalizeInput(input.Name, input.Phone, input.Address, input.CPF)
	if err != nil {
		return nil, err
	}

	item, err := u.repo.GetByIDAndTenant(ctx, actor.TenantID, guardianID)
	if err != nil {
		return nil, err
	}

	item.Name = normalized.name
	item.Phone = normalized.phone
	item.Address = normalized.address
	item.CPF = normalized.cpf
	item.UpdatedAt = time.Now()

	if err := u.repo.Update(ctx, item); err != nil {
		return nil, err
	}

	return ToResponse(item), nil
}

func (u *Usecase) Delete(ctx context.Context, actor security.Identity, guardianID uuid.UUID) error {
	if !canAccessGuardians(actor) {
		return sharederrors.Forbidden("you do not have permission to delete guardians")
	}

	total, err := u.repo.CountLearnersByGuardian(ctx, actor.TenantID, guardianID)
	if err != nil {
		return err
	}
	if total > 0 {
		return sharederrors.Conflict("GUARDIAN_HAS_LEARNERS", "guardian is linked to one or more learners")
	}

	return u.repo.Delete(ctx, actor.TenantID, guardianID)
}

type normalizedInput struct {
	name    string
	phone   string
	address string
	cpf     string
}

func normalizeInput(name string, phone string, address string, cpf string) (normalizedInput, error) {
	if err := sharedvalidator.ValidateName(name); err != nil {
		return normalizedInput{}, err
	}
	if err := sharedvalidator.ValidatePhone(phone); err != nil {
		return normalizedInput{}, err
	}
	if err := sharedvalidator.ValidateRequired("address", address); err != nil {
		return normalizedInput{}, err
	}

	normalizedCPF := sharedvalidator.NormalizeCPFOrCNPJ(cpf)
	if normalizedCPF != "" {
		if len(normalizedCPF) != 11 {
			return normalizedInput{}, sharederrors.Invalid("INVALID_CPF", "cpf must have 11 digits")
		}
		if err := sharedvalidator.ValidateCPFOrCNPJ(normalizedCPF); err != nil {
			return normalizedInput{}, err
		}
	}

	return normalizedInput{
		name:    strings.TrimSpace(name),
		phone:   sharedvalidator.NormalizePhone(phone),
		address: strings.TrimSpace(address),
		cpf:     normalizedCPF,
	}, nil
}

func canAccessGuardians(actor security.Identity) bool {
	return actor.IsInternal() &&
		(actor.Permissions.Patients != permissions.ScopeNone ||
			actor.Permissions.Finance != permissions.ScopeNone)
}
