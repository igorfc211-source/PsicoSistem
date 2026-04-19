package auth

import (
	"api-on/internal/subscription"
	"api-on/internal/tenant"
	"api-on/internal/user"
)

type RegisterInput struct {
	ClinicName string `json:"clinic_name"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Phone      string `json:"phone"`
	PlanSlug   string `json:"plan_slug"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthPayload struct {
	Tenant       *tenant.Response              `json:"tenant"`
	User         *user.Response                `json:"user"`
	Subscription *subscription.SummaryResponse `json:"subscription,omitempty"`
	Token        string                        `json:"token"`
}
