package permission

import (
	"context"
	"path/filepath"
	"testing"

	"api-on/internal/auth"
	"api-on/internal/shared/database"
	"api-on/internal/shared/permissions"
	"api-on/internal/shared/security"
	"api-on/internal/subscription"
	"api-on/internal/tenant"
	"api-on/internal/user"
	jwtpkg "api-on/pkg/jwt"
)

func TestAdminCanUpdateUserPermissions(t *testing.T) {
	t.Parallel()

	store := database.NewStore(filepath.Join(t.TempDir(), "state.json"))
	if err := store.Initialize(); err != nil {
		t.Fatalf("initialize store: %v", err)
	}

	tenantRepo := tenant.NewRepository(store)
	subscriptionRepo := subscription.NewRepository(store)
	userRepo := user.NewRepository(store)
	authUsecase := auth.NewUsecase(tenantRepo, subscriptionRepo, userRepo, jwtpkg.NewJWTService("secret", "issuer"))
	userUsecase := user.NewUsecase(userRepo, subscriptionRepo)
	permissionUsecase := NewUsecase(userRepo)

	registerResult, err := authUsecase.Register(context.Background(), auth.RegisterInput{
		ClinicName:              "Clinica Atlas",
		Name:                    "Owner User",
		Email:                   "owner-atlas@clinica.com",
		Password:                "1234@",
		Phone:                   "19999999999",
		CPFOrCNPJ:               "11444777000161",
		PlanSlug:                "premium",
		PaymentSessionConfirmed: true,
	})
	if err != nil {
		t.Fatalf("register owner: %v", err)
	}

	adminActor := security.Identity{
		UserID:      registerResult.User.ID,
		TenantID:    registerResult.Tenant.ID,
		Role:        user.RoleOwner,
		Type:        security.UserTypeInternal,
		Permissions: registerResult.User.Permissions,
	}

	created, err := userUsecase.Create(context.Background(), adminActor, user.CreateInput{
		Name:     "Bruna Lima",
		Email:    "bruna.atlas@clinica.com",
		Password: "1234@",
		Role:     user.RoleProfessional,
		Status:   user.StatusActive,
	})
	if err != nil {
		t.Fatalf("create professional: %v", err)
	}

	updated, err := permissionUsecase.UpdateByUserID(context.Background(), adminActor, created.ID, UpdateInput{
		Permissions: permissions.AccountPermissions{
			UserDirectory: permissions.ScopeOwn,
			Patients:      permissions.ScopeAll,
			Services:      permissions.ScopeAll,
			Calendar:      permissions.ScopeOwn,
			Finance:       permissions.ScopeNone,
			AIHistory:     permissions.ScopeOwn,
			Plans:         permissions.ScopeAll,
		},
	})
	if err != nil {
		t.Fatalf("update permissions: %v", err)
	}

	if updated.Permissions.Patients != permissions.ScopeAll || updated.Permissions.Finance != permissions.ScopeNone {
		t.Fatalf("expected permissions update to persist, got %#v", updated.Permissions)
	}
}
