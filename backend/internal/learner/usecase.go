package learner

import (
	"context"
	"strings"
	"time"

	sharederrors "api-on/internal/shared/errors"
	"api-on/internal/shared/permissions"
	"api-on/internal/shared/security"
	sharedvalidator "api-on/internal/shared/validator"
	"api-on/internal/subscription"

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

func (u *Usecase) List(ctx context.Context, actor security.Identity, input ListInput) ([]Response, ListMeta, error) {
	if !canAccessLearners(actor) {
		return nil, ListMeta{}, sharederrors.Forbidden("you do not have permission to access learners")
	}
	if input.Status != "" {
		if err := ValidateStatus(input.Status); err != nil {
			return nil, ListMeta{}, err
		}
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

func (u *Usecase) Get(ctx context.Context, actor security.Identity, learnerID uuid.UUID) (*Response, error) {
	if !canAccessLearners(actor) {
		return nil, sharederrors.Forbidden("you do not have permission to access learners")
	}

	item, err := u.repo.GetByIDAndTenant(ctx, actor.TenantID, learnerID)
	if err != nil {
		return nil, err
	}

	return ToResponse(item), nil
}

func (u *Usecase) Create(ctx context.Context, actor security.Identity, input CreateInput) (*Response, error) {
	if !canAccessLearners(actor) {
		return nil, sharederrors.Forbidden("you do not have permission to create learners")
	}

	status := normalizeStatus(input.Status)
	if err := validateLearnerInput(input.Name, status, input.VisitCount, input.SessionPriceCents); err != nil {
		return nil, err
	}
	if status == StatusActive {
		if err := u.ensurePatientCapacity(ctx, actor.TenantID); err != nil {
			return nil, err
		}
	}

	now := time.Now()
	item := &Learner{
		ID:                uuid.New(),
		TenantID:          actor.TenantID,
		Name:              strings.TrimSpace(input.Name),
		PhotoURL:          strings.TrimSpace(input.PhotoURL),
		Gender:            strings.TrimSpace(input.Gender),
		Guardian:          strings.TrimSpace(input.Guardian),
		Age:               strings.TrimSpace(input.Age),
		Status:            status,
		StartDate:         strings.TrimSpace(input.StartDate),
		EndDate:           strings.TrimSpace(input.EndDate),
		VisitCount:        input.VisitCount,
		SessionPriceCents: input.SessionPriceCents,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if err := u.repo.Create(ctx, item); err != nil {
		return nil, err
	}

	return ToResponse(item), nil
}

func (u *Usecase) Update(ctx context.Context, actor security.Identity, learnerID uuid.UUID, input UpdateInput) (*Response, error) {
	if !canAccessLearners(actor) {
		return nil, sharederrors.Forbidden("you do not have permission to update learners")
	}

	status := normalizeStatus(input.Status)
	if err := validateLearnerInput(input.Name, status, input.VisitCount, input.SessionPriceCents); err != nil {
		return nil, err
	}

	item, err := u.repo.GetByIDAndTenant(ctx, actor.TenantID, learnerID)
	if err != nil {
		return nil, err
	}
	if item.Status != StatusActive && status == StatusActive {
		if err := u.ensurePatientCapacity(ctx, actor.TenantID); err != nil {
			return nil, err
		}
	}

	item.Name = strings.TrimSpace(input.Name)
	item.PhotoURL = strings.TrimSpace(input.PhotoURL)
	item.Gender = strings.TrimSpace(input.Gender)
	item.Guardian = strings.TrimSpace(input.Guardian)
	item.Age = strings.TrimSpace(input.Age)
	item.Status = status
	item.StartDate = strings.TrimSpace(input.StartDate)
	item.EndDate = strings.TrimSpace(input.EndDate)
	item.VisitCount = input.VisitCount
	item.SessionPriceCents = input.SessionPriceCents
	item.UpdatedAt = time.Now()

	if err := u.repo.Update(ctx, item); err != nil {
		return nil, err
	}

	return ToResponse(item), nil
}

func (u *Usecase) Deactivate(ctx context.Context, actor security.Identity, learnerID uuid.UUID) error {
	if !canAccessLearners(actor) {
		return sharederrors.Forbidden("you do not have permission to deactivate learners")
	}

	return u.repo.Deactivate(ctx, actor.TenantID, learnerID)
}

func (u *Usecase) ensurePatientCapacity(ctx context.Context, tenantID uuid.UUID) error {
	_, plan, err := u.subscriptionRepo.GetByTenantID(ctx, tenantID)
	if err != nil {
		return err
	}

	total, err := u.repo.CountActiveByTenant(ctx, tenantID)
	if err != nil {
		return err
	}

	if total >= plan.MaxPatients {
		return sharederrors.Conflict("PLAN_PATIENT_LIMIT_REACHED", "current plan has reached the maximum number of active learners")
	}

	return nil
}

func validateLearnerInput(name string, status string, visitCount int, sessionPriceCents int64) error {
	if err := sharedvalidator.ValidateName(name); err != nil {
		return err
	}
	if err := ValidateStatus(status); err != nil {
		return err
	}
	if visitCount < 0 {
		return sharederrors.Invalid("INVALID_VISIT_COUNT", "visit_count cannot be negative")
	}
	if sessionPriceCents < 0 {
		return sharederrors.Invalid("INVALID_SESSION_PRICE", "session_price_cents cannot be negative")
	}

	return nil
}

func normalizeStatus(status string) string {
	status = strings.TrimSpace(status)
	if status == "" {
		return StatusActive
	}
	return status
}

func canAccessLearners(actor security.Identity) bool {
	return actor.IsInternal() &&
		(actor.Permissions.Patients != permissions.ScopeNone ||
			actor.Permissions.Finance != permissions.ScopeNone)
}
