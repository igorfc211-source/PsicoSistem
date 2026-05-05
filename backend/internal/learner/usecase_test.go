package learner_test

import (
	"context"
	"path/filepath"
	"testing"

	"api-on/internal/auth"
	learnerdomain "api-on/internal/learner"
	"api-on/internal/shared/database"
	sharederrors "api-on/internal/shared/errors"
	"api-on/internal/shared/security"
	"api-on/internal/subscription"
	"api-on/internal/tenant"
	"api-on/internal/user"
	jwtpkg "api-on/pkg/jwt"
)

func TestCreateLearnerPersistsSessionPrice(t *testing.T) {
	t.Parallel()

	learnerUsecase, owner := buildLearnerUsecase(t, "clinica-preco@teste.com", "19999999999", "11444777000161")

	created, err := learnerUsecase.Create(context.Background(), owner, learnerdomain.CreateInput{
		Name:              "Lia Andrade",
		Guardian:          "Marina Andrade",
		Status:            learnerdomain.StatusActive,
		VisitCount:        8,
		SessionPriceCents: 15000,
	})
	if err != nil {
		t.Fatalf("create learner: %v", err)
	}
	if created.SessionPriceCents != 15000 {
		t.Fatalf("expected session price 15000, got %d", created.SessionPriceCents)
	}

	found, err := learnerUsecase.Get(context.Background(), owner, created.ID)
	if err != nil {
		t.Fatalf("get learner: %v", err)
	}
	if found.SessionPriceCents != 15000 {
		t.Fatalf("expected stored session price 15000, got %d", found.SessionPriceCents)
	}
}

func TestCreateLearnerRejectsNegativeSessionPrice(t *testing.T) {
	t.Parallel()

	learnerUsecase, owner := buildLearnerUsecase(t, "clinica-preco-negativo@teste.com", "19888888888", "52998224725")

	_, err := learnerUsecase.Create(context.Background(), owner, learnerdomain.CreateInput{
		Name:              "Nina Costa",
		Status:            learnerdomain.StatusActive,
		SessionPriceCents: -1,
	})
	if err == nil {
		t.Fatalf("expected validation error")
	}

	appErr := sharederrors.AsAppError(err)
	if appErr == nil || appErr.Code != "INVALID_SESSION_PRICE" {
		t.Fatalf("expected INVALID_SESSION_PRICE error, got %#v", appErr)
	}
}

func buildLearnerUsecase(t *testing.T, email string, phone string, taxID string) (*learnerdomain.Usecase, security.Identity) {
	t.Helper()

	store := database.NewStore(filepath.Join(t.TempDir(), "state.json"))
	if err := store.Initialize(); err != nil {
		t.Fatalf("initialize store: %v", err)
	}

	tenantRepo := tenant.NewRepository(store)
	subscriptionRepo := subscription.NewRepository(store)
	userRepo := user.NewRepository(store)
	learnerRepo := learnerdomain.NewRepository(store)
	authUsecase := auth.NewUsecase(tenantRepo, subscriptionRepo, userRepo, jwtpkg.NewJWTService("secret", "issuer"))
	learnerUsecase := learnerdomain.NewUsecase(learnerRepo, subscriptionRepo)

	registerResult, err := authUsecase.Register(context.Background(), auth.RegisterInput{
		ClinicName:              "Clinica Aprendentes",
		Name:                    "Owner User",
		Email:                   email,
		Password:                "1234@",
		Phone:                   phone,
		CPFOrCNPJ:               taxID,
		PlanSlug:                "basico",
		PaymentSessionConfirmed: true,
	})
	if err != nil {
		t.Fatalf("register owner: %v", err)
	}

	owner := security.Identity{
		UserID:      registerResult.User.ID,
		TenantID:    registerResult.Tenant.ID,
		Role:        user.RoleOwner,
		Type:        security.UserTypeInternal,
		Permissions: registerResult.User.Permissions,
	}

	return learnerUsecase, owner
}
