package auth

import (
	"context"
	"net/url"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"api-on/internal/shared/database"
	sharedemail "api-on/internal/shared/email"
	sharederrors "api-on/internal/shared/errors"
	"api-on/internal/shared/permissions"
	"api-on/internal/subscription"
	"api-on/internal/tenant"
	userdomain "api-on/internal/user"
	jwtpkg "api-on/pkg/jwt"
)

type fakeEmailSender struct {
	messages []sharedemail.Message
}

func (s *fakeEmailSender) Send(_ context.Context, message sharedemail.Message) error {
	s.messages = append(s.messages, message)
	return nil
}

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

func TestPasswordResetSendsEmailAndChangesPassword(t *testing.T) {
	t.Parallel()

	store := database.NewStore(filepath.Join(t.TempDir(), "state.json"))
	if err := store.Initialize(); err != nil {
		t.Fatalf("initialize store: %v", err)
	}

	tenantRepo := tenant.NewRepository(store)
	subscriptionRepo := subscription.NewRepository(store)
	userRepo := userdomain.NewRepository(store)
	resetRepo := NewPasswordResetRepository(store)
	emailSender := &fakeEmailSender{}
	usecase := NewUsecaseWithPasswordReset(
		tenantRepo,
		subscriptionRepo,
		userRepo,
		resetRepo,
		jwtpkg.NewJWTService("secret", "issuer"),
		emailSender,
		"http://localhost:3000",
		30*time.Minute,
	)

	_, err := usecase.Register(context.Background(), RegisterInput{
		ClinicName:              "Clinica Reset",
		Name:                    "Reset Owner",
		Email:                   "reset@clinica.com",
		Password:                "1234@",
		Phone:                   "11999991111",
		CPFOrCNPJ:               "11444777000161",
		PlanSlug:                "premium",
		PaymentSessionConfirmed: true,
	})
	if err != nil {
		t.Fatalf("register: %v", err)
	}

	if _, err := usecase.RequestPasswordReset(context.Background(), ForgotPasswordInput{Email: "reset@clinica.com"}); err != nil {
		t.Fatalf("request password reset: %v", err)
	}
	if len(emailSender.messages) != 1 {
		t.Fatalf("expected one password reset email, got %d", len(emailSender.messages))
	}

	token := extractResetToken(t, emailSender.messages[0].TextBody)
	if _, err := usecase.ResetPassword(context.Background(), ResetPasswordInput{
		Token:    token,
		Password: "nova@senha",
	}); err != nil {
		t.Fatalf("reset password: %v", err)
	}

	if _, err := usecase.Login(context.Background(), LoginInput{Email: "reset@clinica.com", Password: "1234@"}); err == nil {
		t.Fatalf("expected old password to stop working")
	}
	if _, err := usecase.Login(context.Background(), LoginInput{Email: "reset@clinica.com", Password: "nova@senha"}); err != nil {
		t.Fatalf("expected new password to work: %v", err)
	}

	_, err = usecase.ResetPassword(context.Background(), ResetPasswordInput{
		Token:    token,
		Password: "outra@senha",
	})
	if err == nil {
		t.Fatalf("expected consumed token to be rejected")
	}
	appErr := sharederrors.AsAppError(err)
	if appErr == nil || appErr.Code != "INVALID_PASSWORD_RESET_TOKEN" {
		t.Fatalf("expected INVALID_PASSWORD_RESET_TOKEN, got %#v", appErr)
	}
}

func TestPasswordResetDoesNotRevealUnknownEmail(t *testing.T) {
	t.Parallel()

	store := database.NewStore(filepath.Join(t.TempDir(), "state.json"))
	if err := store.Initialize(); err != nil {
		t.Fatalf("initialize store: %v", err)
	}

	tenantRepo := tenant.NewRepository(store)
	subscriptionRepo := subscription.NewRepository(store)
	userRepo := userdomain.NewRepository(store)
	resetRepo := NewPasswordResetRepository(store)
	emailSender := &fakeEmailSender{}
	usecase := NewUsecaseWithPasswordReset(
		tenantRepo,
		subscriptionRepo,
		userRepo,
		resetRepo,
		jwtpkg.NewJWTService("secret", "issuer"),
		emailSender,
		"http://localhost:3000",
		30*time.Minute,
	)

	result, err := usecase.RequestPasswordReset(context.Background(), ForgotPasswordInput{Email: "ausente@clinica.com"})
	if err != nil {
		t.Fatalf("request unknown password reset: %v", err)
	}
	if result == nil || result.Message == "" {
		t.Fatalf("expected generic success response")
	}
	if len(emailSender.messages) != 0 {
		t.Fatalf("expected no email for unknown account")
	}
}

func extractResetToken(t *testing.T, body string) string {
	t.Helper()

	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		if !strings.Contains(line, "/reset-password?") {
			continue
		}

		parsed, err := url.Parse(line)
		if err != nil {
			t.Fatalf("parse reset url: %v", err)
		}
		token := parsed.Query().Get("token")
		if token == "" {
			t.Fatalf("reset url does not include token")
		}
		return token
	}

	t.Fatalf("reset email does not include reset url: %s", body)
	return ""
}
