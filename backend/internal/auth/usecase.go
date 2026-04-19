package auth

import (
	"context"
	"strings"
	"time"

	"api-on/internal/shared/security"
	sharederrors "api-on/internal/shared/errors"
	sharedvalidator "api-on/internal/shared/validator"
	"api-on/internal/subscription"
	"api-on/internal/tenant"
	"api-on/internal/user"
	"api-on/pkg/hash"
	jwtpkg "api-on/pkg/jwt"

	"github.com/google/uuid"
)

type Usecase struct {
	tenantRepo       tenant.Repository
	subscriptionRepo subscription.Repository
	userRepo         user.Repository
	jwtSvc           *jwtpkg.JWTService
}

func NewUsecase(
	tenantRepo tenant.Repository,
	subscriptionRepo subscription.Repository,
	userRepo user.Repository,
	jwtSvc *jwtpkg.JWTService,
) *Usecase {
	return &Usecase{
		tenantRepo:       tenantRepo,
		subscriptionRepo: subscriptionRepo,
		userRepo:         userRepo,
		jwtSvc:           jwtSvc,
	}
}

func (u *Usecase) Register(ctx context.Context, input RegisterInput) (*AuthPayload, error) {
	if err := sharedvalidator.ValidateClinicName(input.ClinicName); err != nil {
		return nil, err
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
	if err := sharedvalidator.ValidatePhone(input.Phone); err != nil {
		return nil, err
	}

	plan, err := u.subscriptionRepo.FindPlanBySlug(ctx, strings.ToLower(strings.TrimSpace(input.PlanSlug)))
	if err != nil {
		return nil, err
	}

	if _, err := u.userRepo.FindByEmail(ctx, input.Email); err == nil {
		return nil, sharederrors.Conflict("USER_EMAIL_ALREADY_EXISTS", "email already registered")
	} else if appErr := sharederrors.AsAppError(err); appErr != nil && appErr.Code != "USER_NOT_FOUND" {
		return nil, err
	}

	passwordHash, err := hash.Generate(input.Password)
	if err != nil {
		return nil, sharederrors.Internal("could not hash password")
	}

	now := time.Now()
	slug, err := u.nextTenantSlug(ctx, input.ClinicName)
	if err != nil {
		return nil, err
	}

	tenantItem := &tenant.Tenant{
		ID:        uuid.New(),
		Name:      strings.TrimSpace(input.ClinicName),
		Slug:      slug,
		Email:     sharedvalidator.NormalizeEmail(input.Email),
		Phone:     strings.TrimSpace(input.Phone),
		Status:    tenant.StatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}

	userItem := &user.User{
		ID:           uuid.New(),
		TenantID:     tenantItem.ID,
		Name:         strings.TrimSpace(input.Name),
		Email:        sharedvalidator.NormalizeEmail(input.Email),
		PasswordHash: passwordHash,
		Role:         user.RoleOwner,
		Status:       user.StatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	renewalAt := now.AddDate(0, 1, 0)
	subscriptionItem := &subscription.Subscription{
		ID:            uuid.New(),
		TenantID:      tenantItem.ID,
		PlanID:        plan.ID,
		Status:        subscription.StatusActive,
		BillingCycle:  "monthly",
		AmountMonthly: plan.PriceMonthlyCents,
		StartedAt:     now,
		RenewalAt:     &renewalAt,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := u.tenantRepo.Create(ctx, tenantItem); err != nil {
		return nil, err
	}

	if err := u.userRepo.Create(ctx, userItem); err != nil {
		_ = u.tenantRepo.Delete(ctx, tenantItem.ID)
		return nil, err
	}

	if err := u.subscriptionRepo.Create(ctx, subscriptionItem); err != nil {
		_ = u.userRepo.DeleteByID(ctx, userItem.ID)
		_ = u.tenantRepo.Delete(ctx, tenantItem.ID)
		return nil, err
	}

	token, err := u.issueToken(userItem)
	if err != nil {
		_ = u.subscriptionRepo.DeleteByTenantID(ctx, tenantItem.ID)
		_ = u.userRepo.DeleteByID(ctx, userItem.ID)
		_ = u.tenantRepo.Delete(ctx, tenantItem.ID)
		return nil, err
	}

	return &AuthPayload{
		Tenant: tenant.ToResponse(tenantItem),
		User:   user.ToResponse(userItem),
		Subscription: &subscription.SummaryResponse{
			Plan:              plan.Slug,
			Status:            subscriptionItem.Status,
			AmountMonthly:     subscriptionItem.AmountMonthly,
			RenewalAt:         subscriptionItem.RenewalAt,
			HasTestsLibrary:   plan.HasTestsLibrary,
			HasAI:             plan.HasAI,
			HasGuardianPortal: plan.HasGuardianPortal,
			MaxProfessionals:  plan.MaxProfessionals,
			MaxPatients:       plan.MaxPatients,
		},
		Token: token,
	}, nil
}

func (u *Usecase) Login(ctx context.Context, input LoginInput) (*AuthPayload, error) {
	if err := sharedvalidator.ValidateEmail(input.Email); err != nil {
		return nil, err
	}
	if strings.TrimSpace(input.Password) == "" {
		return nil, sharederrors.Invalid("INVALID_PASSWORD", "password is required")
	}

	userItem, err := u.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, sharederrors.Unauthorized("invalid credentials")
	}

	if userItem.Status != user.StatusActive {
		return nil, sharederrors.Forbidden("inactive users cannot authenticate")
	}

	if err := hash.Compare(input.Password, userItem.PasswordHash); err != nil {
		return nil, sharederrors.Unauthorized("invalid credentials")
	}

	tenantItem, err := u.tenantRepo.GetByID(ctx, userItem.TenantID)
	if err != nil {
		return nil, err
	}

	if err := u.userRepo.TouchLastLogin(ctx, userItem.ID, time.Now()); err != nil {
		return nil, err
	}

	refreshedUser, err := u.userRepo.GetByIDAndTenant(ctx, userItem.TenantID, userItem.ID)
	if err != nil {
		return nil, err
	}

	token, err := u.issueToken(refreshedUser)
	if err != nil {
		return nil, err
	}

	return &AuthPayload{
		Tenant: tenant.ToResponse(tenantItem),
		User:   user.ToResponse(refreshedUser),
		Token:  token,
	}, nil
}

func (u *Usecase) Refresh(ctx context.Context, actor security.Identity) (*AuthPayload, error) {
	userItem, err := u.userRepo.GetByIDAndTenant(ctx, actor.TenantID, actor.UserID)
	if err != nil {
		return nil, err
	}

	if userItem.Status != user.StatusActive {
		return nil, sharederrors.Forbidden("inactive users cannot refresh session")
	}

	tenantItem, err := u.tenantRepo.GetByID(ctx, actor.TenantID)
	if err != nil {
		return nil, err
	}

	token, err := u.issueToken(userItem)
	if err != nil {
		return nil, err
	}

	return &AuthPayload{
		Tenant: tenant.ToResponse(tenantItem),
		User:   user.ToResponse(userItem),
		Token:  token,
	}, nil
}

func (u *Usecase) nextTenantSlug(ctx context.Context, clinicName string) (string, error) {
	base := sharedvalidator.Slugify(clinicName)
	candidate := base

	for {
		exists, err := u.tenantRepo.ExistsBySlug(ctx, candidate)
		if err != nil {
			return "", err
		}
		if !exists {
			return candidate, nil
		}
		candidate = base + "-" + uuid.NewString()[:8]
	}
}

func (u *Usecase) issueToken(userItem *user.User) (string, error) {
	token, err := u.jwtSvc.GenerateToken(jwtpkg.TokenInput{
		UserID:   userItem.ID.String(),
		TenantID: userItem.TenantID.String(),
		Role:     userItem.Role,
		Email:    userItem.Email,
		Type:     security.UserTypeInternal,
	})
	if err != nil {
		return "", sharederrors.Internal("could not generate auth token")
	}
	return token, nil
}
