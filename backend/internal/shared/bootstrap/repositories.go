package bootstrap

import (
	"context"

	"api-on/internal/learner"
	"api-on/internal/shared/config"
	"api-on/internal/shared/database"
	"api-on/internal/subscription"
	"api-on/internal/tenant"
	"api-on/internal/user"
)

type Repositories struct {
	TenantRepo       tenant.Repository
	SubscriptionRepo subscription.Repository
	UserRepo         user.Repository
	LearnerRepo      learner.Repository
	Close            func()
}

func BuildRepositories(ctx context.Context, cfg *config.Config) (*Repositories, error) {
	switch cfg.StorageDriver {
	case "postgres":
		postgresDB, err := database.NewPostgres(ctx, cfg)
		if err != nil {
			return nil, err
		}

		return &Repositories{
			TenantRepo:       tenant.NewPostgresRepository(postgresDB.Pool),
			SubscriptionRepo: subscription.NewPostgresRepository(postgresDB.Pool),
			UserRepo:         user.NewPostgresRepository(postgresDB.Pool),
			LearnerRepo:      learner.NewPostgresRepository(postgresDB.Pool),
			Close:            postgresDB.Close,
		}, nil
	default:
		store := database.NewStore(cfg.DataFile)
		if err := store.Initialize(); err != nil {
			return nil, err
		}

		return &Repositories{
			TenantRepo:       tenant.NewRepository(store),
			SubscriptionRepo: subscription.NewRepository(store),
			UserRepo:         user.NewRepository(store),
			LearnerRepo:      learner.NewRepository(store),
			Close:            func() {},
		}, nil
	}
}
