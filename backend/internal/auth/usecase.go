package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"html"
	"net/url"
	"strings"
	"time"

	sharedemail "api-on/internal/shared/email"
	sharederrors "api-on/internal/shared/errors"
	"api-on/internal/shared/permissions"
	"api-on/internal/shared/security"
	sharedvalidator "api-on/internal/shared/validator"
	"api-on/internal/subscription"
	"api-on/internal/tenant"
	"api-on/internal/user"
	"api-on/pkg/hash"
	jwtpkg "api-on/pkg/jwt"

	"github.com/google/uuid"
)

type Usecase struct {
	tenantRepo        tenant.Repository
	subscriptionRepo  subscription.Repository
	userRepo          user.Repository
	passwordResetRepo PasswordResetRepository
	jwtSvc            *jwtpkg.JWTService
	emailSender       sharedemail.Sender
	frontendURL       string
	passwordResetTTL  time.Duration
}

func NewUsecase(
	tenantRepo tenant.Repository,
	subscriptionRepo subscription.Repository,
	userRepo user.Repository,
	jwtSvc *jwtpkg.JWTService,
) *Usecase {
	return NewUsecaseWithPasswordReset(
		tenantRepo,
		subscriptionRepo,
		userRepo,
		nil,
		jwtSvc,
		sharedemail.NoopSender{},
		"http://localhost:3000",
		30*time.Minute,
	)
}

func NewUsecaseWithPasswordReset(
	tenantRepo tenant.Repository,
	subscriptionRepo subscription.Repository,
	userRepo user.Repository,
	passwordResetRepo PasswordResetRepository,
	jwtSvc *jwtpkg.JWTService,
	emailSender sharedemail.Sender,
	frontendURL string,
	passwordResetTTL time.Duration,
) *Usecase {
	if emailSender == nil {
		emailSender = sharedemail.NoopSender{}
	}
	frontendURL = strings.TrimRight(strings.TrimSpace(frontendURL), "/")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}
	if passwordResetTTL <= 0 {
		passwordResetTTL = 30 * time.Minute
	}

	return &Usecase{
		tenantRepo:        tenantRepo,
		subscriptionRepo:  subscriptionRepo,
		userRepo:          userRepo,
		passwordResetRepo: passwordResetRepo,
		jwtSvc:            jwtSvc,
		emailSender:       emailSender,
		frontendURL:       frontendURL,
		passwordResetTTL:  passwordResetTTL,
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
	if err := sharedvalidator.ValidateCPFOrCNPJ(input.CPFOrCNPJ); err != nil {
		return nil, err
	}
	if !input.PaymentSessionConfirmed {
		return nil, sharederrors.Invalid("PAYMENT_SESSION_REQUIRED", "payment session must be confirmed before registration")
	}

	normalizedEmail := sharedvalidator.NormalizeEmail(input.Email)
	normalizedPhone := sharedvalidator.NormalizePhone(input.Phone)
	normalizedDocument := sharedvalidator.NormalizeCPFOrCNPJ(input.CPFOrCNPJ)

	plan, err := u.subscriptionRepo.FindPlanBySlug(ctx, strings.ToLower(strings.TrimSpace(input.PlanSlug)))
	if err != nil {
		return nil, err
	}
	if plan.Slug == "intermediario" && !input.AcceptTrialTerms {
		return nil, sharederrors.Invalid("TRIAL_TERMS_REQUIRED", "accept trial terms before creating the account")
	}

	if _, err := u.userRepo.FindByEmail(ctx, normalizedEmail); err == nil {
		return nil, sharederrors.Conflict("USER_EMAIL_ALREADY_EXISTS", "email already registered")
	} else if appErr := sharederrors.AsAppError(err); appErr != nil && appErr.Code != "USER_NOT_FOUND" {
		return nil, err
	}

	if exists, err := u.tenantRepo.ExistsByCNPJ(ctx, normalizedDocument); err != nil {
		return nil, err
	} else if exists {
		return nil, sharederrors.Conflict("TENANT_DOCUMENT_ALREADY_EXISTS", "cpf_cnpj already registered")
	}

	if exists, err := u.tenantRepo.ExistsByPhone(ctx, normalizedPhone); err != nil {
		return nil, err
	} else if exists {
		return nil, sharederrors.Conflict("TENANT_PHONE_ALREADY_EXISTS", "phone already registered")
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
		CNPJ:      normalizedDocument,
		Email:     normalizedEmail,
		Phone:     normalizedPhone,
		Status:    tenant.StatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}

	userItem := &user.User{
		ID:           uuid.New(),
		TenantID:     tenantItem.ID,
		Name:         strings.TrimSpace(input.Name),
		Email:        normalizedEmail,
		PasswordHash: passwordHash,
		Role:         user.RoleOwner,
		Status:       user.StatusActive,
		Permissions:  permissions.DefaultForRole(user.RoleOwner),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	renewalAt := now.AddDate(0, 1, 0)
	subscriptionStatus := subscription.StatusActive
	amountMonthly := plan.PriceMonthlyCents
	if plan.Slug == "intermediario" {
		subscriptionStatus = subscription.StatusTrialing
		amountMonthly = 0
	}

	subscriptionItem := &subscription.Subscription{
		ID:            uuid.New(),
		TenantID:      tenantItem.ID,
		PlanID:        plan.ID,
		Status:        subscriptionStatus,
		BillingCycle:  "monthly",
		AmountMonthly: amountMonthly,
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
		Tenant:       tenant.ToResponse(tenantItem),
		User:         user.ToResponse(userItem),
		Subscription: u.buildSubscriptionSummary(subscriptionItem, plan),
		Token:        token,
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

	subscriptionItem, plan, err := u.subscriptionRepo.GetByTenantID(ctx, userItem.TenantID)
	if err != nil {
		return nil, err
	}

	return &AuthPayload{
		Tenant:       tenant.ToResponse(tenantItem),
		User:         user.ToResponse(refreshedUser),
		Subscription: u.buildSubscriptionSummary(subscriptionItem, plan),
		Token:        token,
	}, nil
}

func (u *Usecase) RequestPasswordReset(ctx context.Context, input ForgotPasswordInput) (*MessageResponse, error) {
	if err := sharedvalidator.ValidateEmail(input.Email); err != nil {
		return nil, err
	}

	response := passwordResetRequestedResponse()
	if u.passwordResetRepo == nil {
		return response, nil
	}

	userItem, err := u.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		appErr := sharederrors.AsAppError(err)
		if appErr != nil && appErr.Code == "USER_NOT_FOUND" {
			return response, nil
		}
		return nil, err
	}

	if userItem.Status != user.StatusActive {
		return response, nil
	}

	rawToken, err := newPasswordResetToken()
	if err != nil {
		return nil, sharederrors.Internal("could not generate password reset token")
	}

	now := time.Now()
	resetToken := &PasswordResetToken{
		ID:        uuid.New(),
		TenantID:  userItem.TenantID,
		UserID:    userItem.ID,
		TokenHash: hashPasswordResetToken(rawToken),
		ExpiresAt: now.Add(u.passwordResetTTL),
		CreatedAt: now,
	}

	if err := u.passwordResetRepo.Create(ctx, resetToken); err != nil {
		return nil, err
	}

	_ = u.emailSender.Send(ctx, sharedemail.Message{
		To:       userItem.Email,
		Subject:  "Recuperacao de senha do PsicoSistem",
		TextBody: u.buildPasswordResetText(rawToken),
		HTMLBody: u.buildPasswordResetHTML(rawToken),
	})

	return response, nil
}

func (u *Usecase) ResetPassword(ctx context.Context, input ResetPasswordInput) (*MessageResponse, error) {
	token := strings.TrimSpace(input.Token)
	if token == "" {
		return nil, invalidPasswordResetToken()
	}
	if err := sharedvalidator.ValidatePassword(input.Password); err != nil {
		return nil, err
	}
	if u.passwordResetRepo == nil {
		return nil, invalidPasswordResetToken()
	}

	now := time.Now()
	resetToken, err := u.passwordResetRepo.ConsumeValid(ctx, hashPasswordResetToken(token), now)
	if err != nil {
		return nil, err
	}

	userItem, err := u.userRepo.GetByIDAndTenant(ctx, resetToken.TenantID, resetToken.UserID)
	if err != nil {
		return nil, invalidPasswordResetToken()
	}
	if userItem.Status != user.StatusActive {
		return nil, invalidPasswordResetToken()
	}

	passwordHash, err := hash.Generate(input.Password)
	if err != nil {
		return nil, sharederrors.Internal("could not hash password")
	}

	userItem.PasswordHash = passwordHash
	userItem.UpdatedAt = now
	if err := u.userRepo.Update(ctx, userItem); err != nil {
		return nil, err
	}

	return &MessageResponse{Message: "Senha atualizada com sucesso."}, nil
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

	subscriptionItem, plan, err := u.subscriptionRepo.GetByTenantID(ctx, actor.TenantID)
	if err != nil {
		return nil, err
	}

	return &AuthPayload{
		Tenant:       tenant.ToResponse(tenantItem),
		User:         user.ToResponse(userItem),
		Subscription: u.buildSubscriptionSummary(subscriptionItem, plan),
		Token:        token,
	}, nil
}

func passwordResetRequestedResponse() *MessageResponse {
	return &MessageResponse{
		Message: "Se o e-mail informado estiver cadastrado, enviaremos um link de recuperacao.",
	}
}

func newPasswordResetToken() (string, error) {
	buffer := make([]byte, 32)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buffer), nil
}

func hashPasswordResetToken(token string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(token)))
	return hex.EncodeToString(sum[:])
}

func (u *Usecase) buildPasswordResetURL(token string) string {
	resetURL, err := url.Parse(u.frontendURL + "/reset-password")
	if err != nil {
		return u.frontendURL + "/reset-password?token=" + url.QueryEscape(token)
	}

	query := resetURL.Query()
	query.Set("token", token)
	resetURL.RawQuery = query.Encode()
	return resetURL.String()
}

func (u *Usecase) buildPasswordResetText(token string) string {
	link := u.buildPasswordResetURL(token)
	return "Recebemos uma solicitacao para redefinir sua senha no PsicoSistem.\n\n" +
		"Acesse o link abaixo para criar uma nova senha:\n" + link + "\n\n" +
		"Este link expira em " + u.passwordResetTTL.String() + ". Se voce nao solicitou, ignore este e-mail."
}

func (u *Usecase) buildPasswordResetHTML(token string) string {
	link := html.EscapeString(u.buildPasswordResetURL(token))
	return "<p>Recebemos uma solicitacao para redefinir sua senha no PsicoSistem.</p>" +
		"<p><a href=\"" + link + "\">Criar nova senha</a></p>" +
		"<p>Este link expira em " + html.EscapeString(u.passwordResetTTL.String()) + ".</p>" +
		"<p>Se voce nao solicitou, ignore este e-mail.</p>"
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

func (u *Usecase) buildSubscriptionSummary(item *subscription.Subscription, plan *subscription.Plan) *subscription.SummaryResponse {
	if item == nil || plan == nil {
		return nil
	}

	var trialEndsAt *time.Time
	if item.Status == subscription.StatusTrialing {
		trialEndsAt = item.RenewalAt
	}

	return &subscription.SummaryResponse{
		Plan:              plan.Slug,
		Status:            item.Status,
		AmountMonthly:     item.AmountMonthly,
		NextAmountMonthly: plan.PriceMonthlyCents,
		RenewalAt:         item.RenewalAt,
		TrialEndsAt:       trialEndsAt,
		HasTestsLibrary:   plan.HasTestsLibrary,
		HasAI:             plan.HasAI,
		HasGuardianPortal: plan.HasGuardianPortal,
		MaxProfessionals:  plan.MaxProfessionals,
		MaxPatients:       plan.MaxPatients,
	}
}
