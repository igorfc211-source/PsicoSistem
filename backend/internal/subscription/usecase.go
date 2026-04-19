package subscription

import (
	"context"

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

	return &SummaryResponse{
		Plan:              plan.Slug,
		Status:            subscription.Status,
		AmountMonthly:     subscription.AmountMonthly,
		RenewalAt:         subscription.RenewalAt,
		HasTestsLibrary:   plan.HasTestsLibrary,
		HasAI:             plan.HasAI,
		HasGuardianPortal: plan.HasGuardianPortal,
		MaxProfessionals:  plan.MaxProfessionals,
		MaxPatients:       plan.MaxPatients,
	}, nil
}
