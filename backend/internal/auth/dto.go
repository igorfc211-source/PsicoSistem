package auth

import (
	"api-on/internal/subscription"
	"api-on/internal/tenant"
	"api-on/internal/user"
)

type RegisterInput struct {
	ClinicName              string `json:"clinic_name"`
	Name                    string `json:"name"`
	Email                   string `json:"email"`
	Password                string `json:"password"`
	Phone                   string `json:"phone"`
	CPFOrCNPJ               string `json:"cpf_cnpj"`
	PlanSlug                string `json:"plan_slug"`
	PaymentSessionConfirmed bool   `json:"payment_session_confirmed"`
	AcceptTrialTerms        bool   `json:"accept_trial_terms"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ForgotPasswordInput struct {
	Email string `json:"email"`
}

type ResetPasswordInput struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type AuthPayload struct {
	Tenant       *tenant.Response              `json:"tenant"`
	User         *user.Response                `json:"user"`
	Subscription *subscription.SummaryResponse `json:"subscription,omitempty"`
	Token        string                        `json:"token"`
}
