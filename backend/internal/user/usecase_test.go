package user_test

import (
	"context"
	"path/filepath"
	"testing"

	"api-on/internal/auth"
	"api-on/internal/shared/database"
	sharederrors "api-on/internal/shared/errors"
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
		ClinicName: "Clinica Horizonte",
		Name:       "Owner User",
		Email:      "owner@clinica.com",
		Password:   "1234@",
		Phone:      "19999999999",
		PlanSlug:   "premium",
	})
	if err != nil {
		t.Fatalf("register owner: %v", err)
	}

	owner := security.Identity{
		UserID:   registerResult.User.ID,
		TenantID: registerResult.Tenant.ID,
		Role:     userdomain.RoleOwner,
		Type:     security.UserTypeInternal,
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
		UserID:   created.ID,
		TenantID: registerResult.Tenant.ID,
		Role:     userdomain.RoleProfessional,
		Type:     security.UserTypeInternal,
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
