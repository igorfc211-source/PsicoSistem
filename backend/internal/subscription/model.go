package subscription

import (
	"time"

	"github.com/google/uuid"
)

const (
	StatusTrialing  = "trialing"
	StatusActive    = "active"
	StatusPastDue   = "past_due"
	StatusCanceled  = "canceled"
	StatusSuspended = "suspended"
)

type Plan struct {
	ID                uuid.UUID `json:"id"`
	Slug              string    `json:"slug"`
	Name              string    `json:"name"`
	PriceMonthlyCents int64     `json:"price_monthly_cents"`
	HasTestsLibrary   bool      `json:"has_tests_library"`
	HasAI             bool      `json:"has_ai"`
	HasGuardianPortal bool      `json:"has_guardian_portal"`
	MaxProfessionals  int       `json:"max_professionals"`
	MaxPatients       int       `json:"max_patients"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
}

type Subscription struct {
	ID            uuid.UUID  `json:"id"`
	TenantID      uuid.UUID  `json:"tenant_id"`
	PlanID        uuid.UUID  `json:"plan_id"`
	Status        string     `json:"status"`
	BillingCycle  string     `json:"billing_cycle"`
	AmountMonthly int64      `json:"amount_monthly"`
	StartedAt     time.Time  `json:"started_at"`
	RenewalAt     *time.Time `json:"renewal_at,omitempty"`
	CanceledAt    *time.Time `json:"canceled_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type SummaryResponse struct {
	Plan              string     `json:"plan"`
	Status            string     `json:"status"`
	AmountMonthly     int64      `json:"amount_monthly"`
	RenewalAt         *time.Time `json:"renewal_at,omitempty"`
	HasTestsLibrary   bool       `json:"has_tests_library"`
	HasAI             bool       `json:"has_ai"`
	HasGuardianPortal bool       `json:"has_guardian_portal"`
	MaxProfessionals  int        `json:"max_professionals"`
	MaxPatients       int        `json:"max_patients"`
}
