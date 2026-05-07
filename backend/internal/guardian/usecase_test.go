package guardian_test

import (
	"context"
	"path/filepath"
	"testing"

	"api-on/internal/auth"
	guardiandomain "api-on/internal/guardian"
	learnerdomain "api-on/internal/learner"
	"api-on/internal/shared/database"
	sharederrors "api-on/internal/shared/errors"
	"api-on/internal/shared/security"
	"api-on/internal/subscription"
	"api-on/internal/tenant"
	"api-on/internal/user"
	jwtpkg "api-on/pkg/jwt"

	"github.com/google/uuid"
)

func TestCreateGuardianPersistsContactData(t *testing.T) {
	t.Parallel()

	guardianUsecase, _, owner := buildGuardianUsecases(t, "clinica-responsavel@teste.com", "19666666666", "52998224725")

	created, err := guardianUsecase.Create(context.Background(), owner, guardiandomain.CreateInput{
		Name:    "Marina Andrade",
		Phone:   "(11) 98888-7777",
		Address: "Rua das Flores, 123",
		CPF:     "529.982.247-25",
	})
	if err != nil {
		t.Fatalf("create guardian: %v", err)
	}
	if created.Phone != "11988887777" {
		t.Fatalf("expected normalized phone, got %s", created.Phone)
	}
	if created.CPF != "52998224725" {
		t.Fatalf("expected normalized cpf, got %s", created.CPF)
	}

	found, err := guardianUsecase.Get(context.Background(), owner, created.ID)
	if err != nil {
		t.Fatalf("get guardian: %v", err)
	}
	if found.Name != "Marina Andrade" || found.Address != "Rua das Flores, 123" {
		t.Fatalf("unexpected guardian data: %#v", found)
	}
}

func TestDeleteGuardianRejectsLinkedLearner(t *testing.T) {
	t.Parallel()

	guardianUsecase, learnerUsecase, owner := buildGuardianUsecases(t, "clinica-responsavel-vinculo@teste.com", "19555555555", "52998224725")
	guardian, err := guardianUsecase.Create(context.Background(), owner, guardiandomain.CreateInput{
		Name:    "Paulo Andrade",
		Phone:   "11977776666",
		Address: "Avenida Central, 40",
	})
	if err != nil {
		t.Fatalf("create guardian: %v", err)
	}

	createdLearner, err := learnerUsecase.Create(context.Background(), owner, learnerdomain.CreateInput{
		Name:              "Lia Andrade",
		GuardianIDs:       []uuid.UUID{guardian.ID},
		Status:            learnerdomain.StatusActive,
		SessionPriceCents: 15000,
	})
	if err != nil {
		t.Fatalf("create learner: %v", err)
	}

	foundGuardian, err := guardianUsecase.Get(context.Background(), owner, guardian.ID)
	if err != nil {
		t.Fatalf("get guardian: %v", err)
	}
	if len(foundGuardian.LearnerIDs) != 1 || foundGuardian.LearnerIDs[0] != createdLearner.ID {
		t.Fatalf("expected linked learner %s, got %#v", createdLearner.ID, foundGuardian.LearnerIDs)
	}

	err = guardianUsecase.Delete(context.Background(), owner, guardian.ID)
	if err == nil {
		t.Fatalf("expected conflict deleting linked guardian")
	}

	appErr := sharederrors.AsAppError(err)
	if appErr == nil || appErr.Code != "GUARDIAN_HAS_LEARNERS" {
		t.Fatalf("expected GUARDIAN_HAS_LEARNERS error, got %#v", appErr)
	}
}

func buildGuardianUsecases(t *testing.T, email string, phone string, taxID string) (*guardiandomain.Usecase, *learnerdomain.Usecase, security.Identity) {
	t.Helper()

	store := database.NewStore(filepath.Join(t.TempDir(), "state.json"))
	if err := store.Initialize(); err != nil {
		t.Fatalf("initialize store: %v", err)
	}

	tenantRepo := tenant.NewRepository(store)
	subscriptionRepo := subscription.NewRepository(store)
	userRepo := user.NewRepository(store)
	learnerRepo := learnerdomain.NewRepository(store)
	guardianRepo := guardiandomain.NewRepository(store)
	authUsecase := auth.NewUsecase(tenantRepo, subscriptionRepo, userRepo, jwtpkg.NewJWTService("secret", "issuer"))
	guardianUsecase := guardiandomain.NewUsecase(guardianRepo)
	learnerUsecase := learnerdomain.NewUsecase(learnerRepo, subscriptionRepo, guardianRepo)

	registerResult, err := authUsecase.Register(context.Background(), auth.RegisterInput{
		ClinicName:              "Clinica Responsaveis",
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

	return guardianUsecase, learnerUsecase, owner
}
