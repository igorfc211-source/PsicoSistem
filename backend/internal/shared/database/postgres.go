package database

import (
	"context"
	"fmt"

	"api-on/internal/shared/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	Pool *pgxpool.Pool
}

func NewPostgres(ctx context.Context, cfg *config.Config) (*Postgres, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse DATABASE_URL: %w", err)
	}

	poolConfig.MaxConns = cfg.DatabaseMaxConns
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	db := &Postgres{Pool: pool}
	if cfg.DatabaseAutoMigrate {
		if err := db.EnsureSchema(ctx); err != nil {
			pool.Close()
			return nil, err
		}
	}

	return db, nil
}

func (p *Postgres) Close() {
	if p != nil && p.Pool != nil {
		p.Pool.Close()
	}
}

func (p *Postgres) EnsureSchema(ctx context.Context) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS tenants (
			id UUID PRIMARY KEY,
			name TEXT NOT NULL,
			slug TEXT NOT NULL UNIQUE,
			cnpj TEXT NOT NULL DEFAULT '',
			email TEXT NOT NULL,
			phone TEXT NOT NULL,
			status TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS plans (
			id UUID PRIMARY KEY,
			slug TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL,
			price_monthly_cents BIGINT NOT NULL,
			has_tests_library BOOLEAN NOT NULL,
			has_ai BOOLEAN NOT NULL,
			has_guardian_portal BOOLEAN NOT NULL,
			max_professionals INTEGER NOT NULL,
			max_patients INTEGER NOT NULL,
			status TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			role TEXT NOT NULL,
			status TEXT NOT NULL,
			permissions JSONB NOT NULL,
			last_login_at TIMESTAMPTZ NULL,
			created_at TIMESTAMPTZ NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS password_reset_tokens (
			id UUID PRIMARY KEY,
			tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			token_hash TEXT NOT NULL UNIQUE,
			expires_at TIMESTAMPTZ NOT NULL,
			used_at TIMESTAMPTZ NULL,
			created_at TIMESTAMPTZ NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS subscriptions (
			id UUID PRIMARY KEY,
			tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			plan_id UUID NOT NULL REFERENCES plans(id),
			status TEXT NOT NULL,
			billing_cycle TEXT NOT NULL,
			amount_monthly BIGINT NOT NULL,
			started_at TIMESTAMPTZ NOT NULL,
			renewal_at TIMESTAMPTZ NULL,
			canceled_at TIMESTAMPTZ NULL,
			created_at TIMESTAMPTZ NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS learners (
			id UUID PRIMARY KEY,
			tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			name TEXT NOT NULL,
			photo_url TEXT NOT NULL DEFAULT '',
			gender TEXT NOT NULL DEFAULT '',
			guardian TEXT NOT NULL DEFAULT '',
			age TEXT NOT NULL DEFAULT '',
			status TEXT NOT NULL,
			start_date TEXT NOT NULL DEFAULT '',
			end_date TEXT NOT NULL DEFAULT '',
			visit_count INTEGER NOT NULL DEFAULT 0,
			session_price_cents BIGINT NOT NULL DEFAULT 0,
			general_value_cents BIGINT NOT NULL DEFAULT 0,
			created_at TIMESTAMPTZ NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL
		)`,
		`ALTER TABLE learners ADD COLUMN IF NOT EXISTS general_value_cents BIGINT NOT NULL DEFAULT 0`,
		`CREATE TABLE IF NOT EXISTS guardians (
			id UUID PRIMARY KEY,
			tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			name TEXT NOT NULL,
			relationship TEXT NOT NULL DEFAULT '',
			phone TEXT NOT NULL,
			address TEXT NOT NULL,
			cpf TEXT NOT NULL DEFAULT '',
			created_at TIMESTAMPTZ NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL
		)`,
		`ALTER TABLE guardians ADD COLUMN IF NOT EXISTS relationship TEXT NOT NULL DEFAULT ''`,
		`CREATE TABLE IF NOT EXISTS learner_guardians (
			tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			learner_id UUID NOT NULL REFERENCES learners(id) ON DELETE CASCADE,
			guardian_id UUID NOT NULL REFERENCES guardians(id) ON DELETE CASCADE,
			created_at TIMESTAMPTZ NOT NULL,
			PRIMARY KEY (learner_id, guardian_id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_users_tenant_id ON users (tenant_id)`,
		`CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_user_id ON password_reset_tokens (user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_token_hash ON password_reset_tokens (token_hash)`,
		`CREATE INDEX IF NOT EXISTS idx_subscriptions_tenant_id ON subscriptions (tenant_id)`,
		`CREATE INDEX IF NOT EXISTS idx_learners_tenant_id ON learners (tenant_id)`,
		`CREATE INDEX IF NOT EXISTS idx_guardians_tenant_id ON guardians (tenant_id)`,
		`CREATE INDEX IF NOT EXISTS idx_learner_guardians_learner_id ON learner_guardians (learner_id)`,
		`CREATE INDEX IF NOT EXISTS idx_learner_guardians_guardian_id ON learner_guardians (guardian_id)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_tenants_phone_unique ON tenants (phone)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_tenants_cnpj_unique ON tenants (cnpj) WHERE cnpj <> ''`,
	}

	for _, statement := range statements {
		if _, err := p.Pool.Exec(ctx, statement); err != nil {
			return fmt.Errorf("ensure postgres schema: %w", err)
		}
	}

	for _, plan := range seedState().Plans {
		if _, err := p.Pool.Exec(ctx, `
			INSERT INTO plans (
				id, slug, name, price_monthly_cents, has_tests_library, has_ai,
				has_guardian_portal, max_professionals, max_patients, status, created_at
			) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
			ON CONFLICT (slug) DO UPDATE SET
				name = EXCLUDED.name,
				price_monthly_cents = EXCLUDED.price_monthly_cents,
				has_tests_library = EXCLUDED.has_tests_library,
				has_ai = EXCLUDED.has_ai,
				has_guardian_portal = EXCLUDED.has_guardian_portal,
				max_professionals = EXCLUDED.max_professionals,
				max_patients = EXCLUDED.max_patients,
				status = EXCLUDED.status
		`,
			plan.ID,
			plan.Slug,
			plan.Name,
			plan.PriceMonthlyCents,
			plan.HasTestsLibrary,
			plan.HasAI,
			plan.HasGuardianPortal,
			plan.MaxProfessionals,
			plan.MaxPatients,
			plan.Status,
			plan.CreatedAt,
		); err != nil {
			return fmt.Errorf("seed plans in postgres: %w", err)
		}
	}

	return nil
}
