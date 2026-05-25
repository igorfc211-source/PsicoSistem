package user_test

import (
	"context"
	"path/filepath"
	"testing"

	"api-on/internal/auth"
	"api-on/internal/shared/database"
	sharederrors "api-on/internal/shared/errors"
	"api-on/internal/shared/permissions"
	"api-on/internal/shared/security"
	"api-on/internal/subscription"
	"api-on/internal/tenant"
	userdomain "api-on/internal/user"
	jwtpkg "api-on/pkg/jwt"
)

func TestProfessionalCannotManageUsers(t *testing.T) {
	t.Parallel()

	store := database.NewStore(filepath.Join(t.TempDir(), "state.json"))
	if err := store.Initialize(); err != nil {
		t.Fatalf("initialize store: %v", err)
	}

	tenantRepo := tenant.NewRepository(store)
	subscriptionRepo := subscription.NewRepository(store)
	userRepo := userdomain.NewRepository(store)
	authUsecase := auth.NewUsecase(tenantRepo, subscriptionRepo, userRepo, jwtpkg.NewJWTService("secret", "issuer"))
	userUsecase := userdomain.NewUsecase(userRepo, subscriptionRepo)

	registerResult, err := authUsecase.Register(context.Background(), auth.RegisterInput{
		ClinicName:              "Clinica Horizonte",
		Name:                    "Owner User",
		Email:                   "owner@clinica.com",
		Password:                "1234@",
		Phone:                   "19999999999",
		CPFOrCNPJ:               "11444777000161",
		PlanSlug:                "premium",
		PaymentSessionConfirmed: true,
	})
	if err != nil {
		t.Fatalf("register owner: %v", err)
	}

	owner := security.Identity{
		UserID:      registerResult.User.ID,
		TenantID:    registerResult.Tenant.ID,
		Role:        userdomain.RoleOwner,
		Type:        security.UserTypeInternal,
		Permissions: registerResult.User.Permissions,
	}

	created, err := userUsecase.Create(context.Background(), owner, userdomain.CreateInput{
		Name:     "Bruna Lima",
		Email:    "bruna@clinica.com",
		Password: "1234@",
		Role:     userdomain.RoleProfessional,
		Status:   userdomain.StatusActive,
	})
	if err != nil {
		t.Fatalf("create professional: %v", err)
	}

	professional := security.Identity{
		UserID:      created.ID,
		TenantID:    registerResult.Tenant.ID,
		Role:        userdomain.RoleProfessional,
		Type:        security.UserTypeInternal,
		Permissions: permissions.DefaultForRole(userdomain.RoleProfessional),
	}

	_, _, err = userUsecase.List(context.Background(), professional, userdomain.ParseListInput("1", "20", "", "", ""))
	if err == nil {
		t.Fatalf("expected forbidden error")
	}

	appErr := sharederrors.AsAppError(err)
	if appErr == nil || appErr.Code != "FORBIDDEN" {
		t.Fatalf("expected FORBIDDEN error, got %#v", appErr)
	}
}

func TestProfessionalReceivesOwnScopedDefaults(t *testing.T) {
	t.Parallel()

	store := database.NewStore(filepath.Join(t.TempDir(), "state.json"))
	if err := store.Initialize(); err != nil {
		t.Fatalf("initialize store: %v", err)
	}

	tenantRepo := tenant.NewRepository(store)
	subscriptionRepo := subscription.NewRepository(store)
	userRepo := userdomain.NewRepository(store)
	authUsecase := auth.NewUsecase(tenantRepo, subscriptionRepo, userRepo, jwtpkg.NewJWTService("secret", "issuer"))
	userUsecase := userdomain.NewUsecase(userRepo, subscriptionRepo)

	registerResult, err := authUsecase.Register(context.Background(), auth.RegisterInput{
		ClinicName:              "Clinica Prisma",
		Name:                    "Owner User",
		Email:                   "owner2@clinica.com",
		Password:                "1234@",
		Phone:                   "19999999999",
		CPFOrCNPJ:               "52998224725",
		PlanSlug:                "premium",
		PaymentSessionConfirmed: true,
	})
	if err != nil {
		t.Fatalf("register owner: %v", err)
	}

	owner := security.Identity{
		UserID:      registerResult.User.ID,
		TenantID:    registerResult.Tenant.ID,
		Role:        userdomain.RoleOwner,
		Type:        security.UserTypeInternal,
		Permissions: registerResult.User.Permissions,
	}

	created, err := userUsecase.Create(context.Background(), owner, userdomain.CreateInput{
		Name:     "Carlos Melo",
		Email:    "carlos@clinica.com",
		Password: "1234@",
		Role:     userdomain.RoleProfessional,
		Status:   userdomain.StatusActive,
	})
	if err != nil {
		t.Fatalf("create professional: %v", err)
	}

	if created.Permissions.Patients != permissions.ScopeOwn ||
		created.Permissions.Services != permissions.ScopeOwn ||
		created.Permissions.Calendar != permissions.ScopeOwn ||
		created.Permissions.Finance != permissions.ScopeOwn ||
		created.Permissions.AIHistory != permissions.ScopeOwn ||
		created.Permissions.Plans != permissions.ScopeOwn {
		t.Fatalf("expected professional to receive own-scoped defaults, got %#v", created.Permissions)
	}
}

func TestGetMeUsesAuthenticatedTenantScope(t *testing.T) {
	t.Parallel()

	store := database.NewStore(filepath.Join(t.TempDir(), "state.json"))
	if err := store.Initialize(); err != nil {
		t.Fatalf("initialize store: %v", err)
	}

	tenantRepo := tenant.NewRepository(store)
	subscriptionRepo := subscription.NewRepository(store)
	userRepo := userdomain.NewRepository(store)
	authUsecase := auth.NewUsecase(tenantRepo, subscriptionRepo, userRepo, jwtpkg.NewJWTService("secret", "issuer"))
	userUsecase := userdomain.NewUsecase(userRepo, subscriptionRepo)

	firstAccount, err := authUsecase.Register(context.Background(), auth.RegisterInput{
		ClinicName:              "Clinica Alfa",
		Name:                    "Owner Alfa",
		Email:                   "owner-alfa@clinica.com",
		Password:                "1234@",
		Phone:                   "11999999991",
		CPFOrCNPJ:               "11444777000161",
		PlanSlug:                "premium",
		PaymentSessionConfirmed: true,
	})
	if err != nil {
		t.Fatalf("register first owner: %v", err)
	}

	secondAccount, err := authUsecase.Register(context.Background(), auth.RegisterInput{
		ClinicName:              "Clinica Beta",
		Name:                    "Owner Beta",
		Email:                   "owner-beta@clinica.com",
		Password:                "1234@",
		Phone:                   "11999999992",
		CPFOrCNPJ:               "52998224725",
		PlanSlug:                "premium",
		PaymentSessionConfirmed: true,
	})
	if err != nil {
		t.Fatalf("register second owner: %v", err)
	}

	firstActor := security.Identity{
		UserID:      firstAccount.User.ID,
		TenantID:    firstAccount.Tenant.ID,
		Role:        userdomain.RoleOwner,
		Type:        security.UserTypeInternal,
		Permissions: firstAccount.User.Permissions,
	}

	currentUser, err := userUsecase.GetMe(context.Background(), firstActor)
	if err != nil {
		t.Fatalf("get current user: %v", err)
	}
	if currentUser.ID != firstAccount.User.ID || currentUser.TenantID != firstAccount.Tenant.ID {
		t.Fatalf("expected authenticated account, got user=%s tenant=%s", currentUser.ID, currentUser.TenantID)
	}

	crossTenantActor := security.Identity{
		UserID:      secondAccount.User.ID,
		TenantID:    firstAccount.Tenant.ID,
		Role:        userdomain.RoleOwner,
		Type:        security.UserTypeInternal,
		Permissions: secondAccount.User.Permissions,
	}

	_, err = userUsecase.GetMe(context.Background(), crossTenantActor)
	if err == nil {
		t.Fatalf("expected cross-tenant account lookup to be blocked")
	}

	appErr := sharederrors.AsAppError(err)
	if appErr == nil || appErr.Code != "USER_NOT_FOUND" {
		t.Fatalf("expected USER_NOT_FOUND for cross-tenant lookup, got %#v", appErr)
	}
}
