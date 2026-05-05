package auth

import (
	"context"
	"path/filepath"
	"testing"

	"api-on/internal/shared/database"
	sharederrors "api-on/internal/shared/errors"
	"api-on/internal/shared/permissions"
	"api-on/internal/subscription"
	"api-on/internal/tenant"
	userdomain "api-on/internal/user"
	jwtpkg "api-on/pkg/jwt"
)

func TestRegisterCreatesTenantOwnerAndSubscription(t *testing.T) {
	t.Parallel()

	store := database.NewStore(filepath.Join(t.TempDir(), "state.json"))
	if err := store.Initialize(); err != nil {
		t.Fatalf("initialize store: %v", err)
	}

	tenantRepo := tenant.NewRepository(store)
	subscriptionRepo := subscription.NewRepository(store)
	userRepo := userdomain.NewRepository(store)
	usecase := NewUsecase(tenantRepo, subscriptionRepo, userRepo, jwtpkg.NewJWTService("secret", "issuer"))

	result, err := usecase.Register(context.Background(), RegisterInput{
		ClinicName:              "Clinica Aurora",
		Name:                    "Ana Souza",
		Email:                   "ana@clinica.com",
		Password:                "1234@",
		Phone:                   "19999999999",
		CPFOrCNPJ:               "11444777000161",
		PlanSlug:                "basico",
		PaymentSessionConfirmed: true,
	})
	if err != nil {
		t.Fatalf("register: %v", err)
	}

	if result.Tenant == nil || result.User == nil || result.Subscription == nil {
		t.Fatalf("expected tenant, user and subscription in payload")
	}

	if result.User.Role != userdomain.RoleOwner {
		t.Fatalf("expected owner role, got %s", result.User.Role)
	}

	if result.User.Permissions.UserDirectory != permissions.ScopeAll {
		t.Fatalf("expected owner to receive all user directory permissions")
	}

	if result.Subscription.Plan != "basico" {
		t.Fatalf("expected basico plan, got %s", result.Subscription.Plan)
	}
	if result.Subscription.NextAmountMonthly != 9700 {
		t.Fatalf("expected next amount to keep plan price, got %d", result.Subscription.NextAmountMonthly)
	}

	savedUser, err := userRepo.FindByEmail(context.Background(), "ana@clinica.com")
	if err != nil {
		t.Fatalf("find saved user: %v", err)
	}

	if savedUser.TenantID != result.Tenant.ID {
		t.Fatalf("expected saved user to belong to created tenant")
	}
}

func TestRegisterAppliesIntermediarioTrial(t *testing.T) {
	t.Parallel()

	store := database.NewStore(filepath.Join(t.TempDir(), "state.json"))
	if err := store.Initialize(); err != nil {
		t.Fatalf("initialize store: %v", err)
	}

	tenantRepo := tenant.NewRepository(store)
	subscriptionRepo := subscription.NewRepository(store)
	userRepo := userdomain.NewRepository(store)
	usecase := NewUsecase(tenantRepo, subscriptionRepo, userRepo, jwtpkg.NewJWTService("secret", "issuer"))

	result, err := usecase.Register(context.Background(), RegisterInput{
		ClinicName:              "Clinica Trial",
		Name:                    "Marina Trial",
		Email:                   "trial@clinica.com",
		Password:                "1234@",
		Phone:                   "11988887777",
		CPFOrCNPJ:               "52998224725",
		PlanSlug:                "intermediario",
		PaymentSessionConfirmed: true,
		AcceptTrialTerms:        true,
	})
	if err != nil {
		t.Fatalf("register with trial: %v", err)
	}

	if result.Subscription == nil {
		t.Fatalf("expected subscription in payload")
	}
	if result.Subscription.Status != subscription.StatusTrialing {
		t.Fatalf("expected trialing subscription, got %s", result.Subscription.Status)
	}
	if result.Subscription.AmountMonthly != 0 {
		t.Fatalf("expected trial to start with zero amount, got %d", result.Subscription.AmountMonthly)
	}
	if result.Subscription.NextAmountMonthly != 14700 {
		t.Fatalf("expected next amount to reflect intermediario price, got %d", result.Subscription.NextAmountMonthly)
	}
	if result.Subscription.TrialEndsAt == nil {
		t.Fatalf("expected trial end date to be informed")
	}
}

func TestRegisterRejectsDuplicatePhone(t *testing.T) {
	t.Parallel()

	store := database.NewStore(filepath.Join(t.TempDir(), "state.json"))
	if err := store.Initialize(); err != nil {
		t.Fatalf("initialize store: %v", err)
	}

	tenantRepo := tenant.NewRepository(store)
	subscriptionRepo := subscription.NewRepository(store)
	userRepo := userdomain.NewRepository(store)
	usecase := NewUsecase(tenantRepo, subscriptionRepo, userRepo, jwtpkg.NewJWTService("secret", "issuer"))

	_, err := usecase.Register(context.Background(), RegisterInput{
		ClinicName:              "Clinica Primeira",
		Name:                    "Dra Primeira",
		Email:                   "primeira@clinica.com",
		Password:                "1234@",
		Phone:                   "(11) 97777-6666",
		CPFOrCNPJ:               "11444777000161",
		PlanSlug:                "premium",
		PaymentSessionConfirmed: true,
	})
	if err != nil {
		t.Fatalf("register first clinic: %v", err)
	}

	_, err = usecase.Register(context.Background(), RegisterInput{
		ClinicName:              "Clinica Segunda",
		Name:                    "Dra Segunda",
		Email:                   "segunda@clinica.com",
		Password:                "1234@",
		Phone:                   "11977776666",
		CPFOrCNPJ:               "52998224725",
		PlanSlug:                "basico",
		PaymentSessionConfirmed: true,
	})
	if err == nil {
		t.Fatalf("expected duplicate phone conflict")
	}

	appErr := sharederrors.AsAppError(err)
	if appErr == nil || appErr.Code != "TENANT_PHONE_ALREADY_EXISTS" {
		t.Fatalf("expected TENANT_PHONE_ALREADY_EXISTS, got %#v", appErr)
	}
}
