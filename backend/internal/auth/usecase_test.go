package auth

import (
	"context"
	"path/filepath"
	"testing"

	"api-on/internal/shared/database"
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
		ClinicName: "Clinica Aurora",
		Name:       "Ana Souza",
		Email:      "ana@clinica.com",
		Password:   "1234@",
		Phone:      "19999999999",
		PlanSlug:   "basico",
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

	if result.Subscription.Plan != "basico" {
		t.Fatalf("expected basico plan, got %s", result.Subscription.Plan)
	}

	savedUser, err := userRepo.FindByEmail(context.Background(), "ana@clinica.com")
	if err != nil {
		t.Fatalf("find saved user: %v", err)
	}

	if savedUser.TenantID != result.Tenant.ID {
		t.Fatalf("expected saved user to belong to created tenant")
	}
}
