package subscription

import (
	"context"

	sharederrors "api-on/internal/shared/errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) FindPlanBySlug(ctx context.Context, slug string) (*Plan, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, slug, name, price_monthly_cents, has_tests_library, has_ai,
		       has_guardian_portal, max_professionals, max_patients, status, created_at
		FROM plans
		WHERE slug = $1
	`, slug)

	var item Plan
	if err := row.Scan(
		&item.ID,
		&item.Slug,
		&item.Name,
		&item.PriceMonthlyCents,
		&item.HasTestsLibrary,
		&item.HasAI,
		&item.HasGuardianPortal,
		&item.MaxProfessionals,
		&item.MaxPatients,
		&item.Status,
		&item.CreatedAt,
	); err != nil {
		return nil, sharederrors.NotFound("PLAN_NOT_FOUND", "plan not found")
	}

	return &item, nil
}

func (r *PostgresRepository) FindPlanByID(ctx context.Context, id uuid.UUID) (*Plan, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, slug, name, price_monthly_cents, has_tests_library, has_ai,
		       has_guardian_portal, max_professionals, max_patients, status, created_at
		FROM plans
		WHERE id = $1
	`, id)

	var item Plan
	if err := row.Scan(
		&item.ID,
		&item.Slug,
		&item.Name,
		&item.PriceMonthlyCents,
		&item.HasTestsLibrary,
		&item.HasAI,
		&item.HasGuardianPortal,
		&item.MaxProfessionals,
		&item.MaxPatients,
		&item.Status,
		&item.CreatedAt,
	); err != nil {
		return nil, sharederrors.NotFound("PLAN_NOT_FOUND", "plan not found")
	}

	return &item, nil
}

func (r *PostgresRepository) Create(ctx context.Context, item *Subscription) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO subscriptions (
			id, tenant_id, plan_id, status, billing_cycle, amount_monthly,
			started_at, renewal_at, canceled_at, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`,
		item.ID,
		item.TenantID,
		item.PlanID,
		item.Status,
		item.BillingCycle,
		item.AmountMonthly,
		item.StartedAt,
		item.RenewalAt,
		item.CanceledAt,
		item.CreatedAt,
		item.UpdatedAt,
	)
	return err
}

func (r *PostgresRepository) GetByTenantID(ctx context.Context, tenantID uuid.UUID) (*Subscription, *Plan, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT s.id, s.tenant_id, s.plan_id, s.status, s.billing_cycle, s.amount_monthly,
		       s.started_at, s.renewal_at, s.canceled_at, s.created_at, s.updated_at,
		       p.id, p.slug, p.name, p.price_monthly_cents, p.has_tests_library, p.has_ai,
		       p.has_guardian_portal, p.max_professionals, p.max_patients, p.status, p.created_at
		FROM subscriptions s
		INNER JOIN plans p ON p.id = s.plan_id
		WHERE s.tenant_id = $1
		LIMIT 1
	`, tenantID)

	var subscriptionItem Subscription
	var planItem Plan
	if err := row.Scan(
		&subscriptionItem.ID,
		&subscriptionItem.TenantID,
		&subscriptionItem.PlanID,
		&subscriptionItem.Status,
		&subscriptionItem.BillingCycle,
		&subscriptionItem.AmountMonthly,
		&subscriptionItem.StartedAt,
		&subscriptionItem.RenewalAt,
		&subscriptionItem.CanceledAt,
		&subscriptionItem.CreatedAt,
		&subscriptionItem.UpdatedAt,
		&planItem.ID,
		&planItem.Slug,
		&planItem.Name,
		&planItem.PriceMonthlyCents,
		&planItem.HasTestsLibrary,
		&planItem.HasAI,
		&planItem.HasGuardianPortal,
		&planItem.MaxProfessionals,
		&planItem.MaxPatients,
		&planItem.Status,
		&planItem.CreatedAt,
	); err != nil {
		return nil, nil, sharederrors.NotFound("SUBSCRIPTION_NOT_FOUND", "subscription not found")
	}

	return &subscriptionItem, &planItem, nil
}

func (r *PostgresRepository) DeleteByTenantID(ctx context.Context, tenantID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM subscriptions WHERE tenant_id = $1`, tenantID)
	return err
}
