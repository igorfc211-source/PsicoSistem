package subscription

import (
	"context"
	"time"

	"api-on/internal/shared/security"
)

type Usecase struct {
	repo Repository
}

func NewUsecase(repo Repository) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) GetCurrent(ctx context.Context, actor security.Identity) (*SummaryResponse, error) {
	subscription, plan, err := u.repo.GetByTenantID(ctx, actor.TenantID)
	if err != nil {
		return nil, err
	}

	var trialEndsAt *time.Time
	if subscription.Status == StatusTrialing {
		trialEndsAt = subscription.RenewalAt
	}

	return &SummaryResponse{
		Plan:              plan.Slug,
		Status:            subscription.Status,
		AmountMonthly:     subscription.AmountMonthly,
		NextAmountMonthly: plan.PriceMonthlyCents,
		RenewalAt:         subscription.RenewalAt,
		TrialEndsAt:       trialEndsAt,
		HasTestsLibrary:   plan.HasTestsLibrary,
		HasAI:             plan.HasAI,
		HasGuardianPortal: plan.HasGuardianPortal,
		MaxProfessionals:  plan.MaxProfessionals,
		MaxPatients:       plan.MaxPatients,
	}, nil
}
