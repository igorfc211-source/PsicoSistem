package tenant

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

func (u *Usecase) GetCurrent(ctx context.Context, actor security.Identity) (*Response, error) {
	item, err := u.repo.GetByID(ctx, actor.TenantID)
	if err != nil {
		return nil, err
	}

	return ToResponse(item), nil
}
