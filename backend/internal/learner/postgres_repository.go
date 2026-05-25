package learner

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

func (r *PostgresRepository) Create(ctx context.Context, item *Learner) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO learners (
			id, tenant_id, name, photo_url, gender, guardian, age, status,
			start_date, end_date, visit_count, session_price_cents, general_value_cents,
			created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
	`,
		item.ID,
		item.TenantID,
		item.Name,
		item.PhotoURL,
		item.Gender,
		item.Guardian,
		item.Age,
		item.Status,
		item.StartDate,
		item.EndDate,
		item.VisitCount,
		item.SessionPriceCents,
		item.GeneralValueCents,
		item.CreatedAt,
		item.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository) GetByIDAndTenant(ctx context.Context, tenantID uuid.UUID, learnerID uuid.UUID) (*Learner, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, name, photo_url, gender, guardian, age, status,
		       start_date, end_date, visit_count, session_price_cents, general_value_cents,
		       created_at, updated_at
		FROM learners
		WHERE id = $1 AND tenant_id = $2
	`, learnerID, tenantID)

	item, err := scanLearner(row)
	if err != nil {
		return nil, err
	}

	guardianIDs, err := r.loadGuardianIDsByLearner(ctx, tenantID, learnerID)
	if err != nil {
		return nil, err
	}
	item.GuardianIDs = guardianIDs

	return item, nil
}

func (r *PostgresRepository) ExistsByIDAndTenant(ctx context.Context, tenantID uuid.UUID, learnerID uuid.UUID) (bool, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM learners WHERE id = $1 AND tenant_id = $2
		)
	`, learnerID, tenantID)

	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (r *PostgresRepository) ListByTenant(ctx context.Context, tenantID uuid.UUID, input ListInput) ([]Learner, int, error) {
	filteredRows, err := r.pool.Query(ctx, `
		SELECT id, tenant_id, name, photo_url, gender, guardian, age, status,
		       start_date, end_date, visit_count, session_price_cents, general_value_cents,
		       created_at, updated_at
		FROM learners
		WHERE tenant_id = $1
		  AND ($2 = '' OR status = $2)
		  AND (
		  	$3 = ''
		  	OR LOWER(name) LIKE '%' || LOWER($3) || '%'
		  	OR LOWER(guardian) LIKE '%' || LOWER($3) || '%'
		  	OR LOWER(gender) LIKE '%' || LOWER($3) || '%'
		  	OR LOWER(age) LIKE '%' || LOWER($3) || '%'
		  )
		ORDER BY created_at ASC
		LIMIT $4 OFFSET $5
	`, tenantID, input.Status, input.Search, input.PerPage, input.Offset())
	if err != nil {
		return nil, 0, err
	}
	defer filteredRows.Close()

	items := make([]Learner, 0)
	for filteredRows.Next() {
		item, err := scanLearner(filteredRows)
		if err != nil {
			return nil, 0, err
		}

		guardianIDs, err := r.loadGuardianIDsByLearner(ctx, tenantID, item.ID)
		if err != nil {
			return nil, 0, err
		}
		item.GuardianIDs = guardianIDs
		items = append(items, *item)
	}

	countRow := r.pool.QueryRow(ctx, `
		SELECT COUNT(1)
		FROM learners
		WHERE tenant_id = $1
		  AND ($2 = '' OR status = $2)
		  AND (
		  	$3 = ''
		  	OR LOWER(name) LIKE '%' || LOWER($3) || '%'
		  	OR LOWER(guardian) LIKE '%' || LOWER($3) || '%'
		  	OR LOWER(gender) LIKE '%' || LOWER($3) || '%'
		  	OR LOWER(age) LIKE '%' || LOWER($3) || '%'
		  )
	`, tenantID, input.Status, input.Search)

	var total int
	if err := countRow.Scan(&total); err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *PostgresRepository) Update(ctx context.Context, item *Learner) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE learners
		SET name = $1, photo_url = $2, gender = $3, guardian = $4, age = $5,
		    status = $6, start_date = $7, end_date = $8, visit_count = $9,
		    session_price_cents = $10, general_value_cents = $11, created_at = $12,
		    updated_at = $13
		WHERE id = $14 AND tenant_id = $15
	`,
		item.Name,
		item.PhotoURL,
		item.Gender,
		item.Guardian,
		item.Age,
		item.Status,
		item.StartDate,
		item.EndDate,
		item.VisitCount,
		item.SessionPriceCents,
		item.GeneralValueCents,
		item.CreatedAt,
		item.UpdatedAt,
		item.ID,
		item.TenantID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return sharederrors.NotFound("LEARNER_NOT_FOUND", "learner not found")
	}

	return nil
}

func (r *PostgresRepository) Deactivate(ctx context.Context, tenantID uuid.UUID, learnerID uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE learners
		SET status = $1, updated_at = NOW()
		WHERE id = $2 AND tenant_id = $3
	`, StatusInactive, learnerID, tenantID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return sharederrors.NotFound("LEARNER_NOT_FOUND", "learner not found")
	}

	return nil
}

func (r *PostgresRepository) CountActiveByTenant(ctx context.Context, tenantID uuid.UUID) (int, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT COUNT(1)
		FROM learners
		WHERE tenant_id = $1 AND status = $2
	`, tenantID, StatusActive)

	var total int
	if err := row.Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

type learnerScanner interface {
	Scan(dest ...any) error
}

func scanLearner(scanner learnerScanner) (*Learner, error) {
	var item Learner
	if err := scanner.Scan(
		&item.ID,
		&item.TenantID,
		&item.Name,
		&item.PhotoURL,
		&item.Gender,
		&item.Guardian,
		&item.Age,
		&item.Status,
		&item.StartDate,
		&item.EndDate,
		&item.VisitCount,
		&item.SessionPriceCents,
		&item.GeneralValueCents,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		return nil, sharederrors.NotFound("LEARNER_NOT_FOUND", "learner not found")
	}

	return &item, nil
}

func (r *PostgresRepository) loadGuardianIDsByLearner(ctx context.Context, tenantID uuid.UUID, learnerID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT guardian_id
		FROM learner_guardians
		WHERE tenant_id = $1 AND learner_id = $2
		ORDER BY guardian_id ASC
	`, tenantID, learnerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	guardianIDs := make([]uuid.UUID, 0)
	for rows.Next() {
		var guardianID uuid.UUID
		if err := rows.Scan(&guardianID); err != nil {
			return nil, err
		}
		guardianIDs = append(guardianIDs, guardianID)
	}

	return guardianIDs, nil
}
