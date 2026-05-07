package guardian

import (
	"context"

	sharederrors "api-on/internal/shared/errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(ctx context.Context, item *Guardian) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO guardians (
			id, tenant_id, name, phone, address, cpf, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`,
		item.ID,
		item.TenantID,
		item.Name,
		item.Phone,
		item.Address,
		item.CPF,
		item.CreatedAt,
		item.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository) GetByIDAndTenant(ctx context.Context, tenantID uuid.UUID, guardianID uuid.UUID) (*Guardian, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, name, phone, address, cpf, created_at, updated_at
		FROM guardians
		WHERE id = $1 AND tenant_id = $2
	`, guardianID, tenantID)

	item, err := scanGuardian(row)
	if err != nil {
		return nil, err
	}

	learnerIDs, err := r.loadLearnerIDsByGuardian(ctx, tenantID, guardianID)
	if err != nil {
		return nil, err
	}
	item.LearnerIDs = learnerIDs

	return item, nil
}

func (r *PostgresRepository) ListByTenant(ctx context.Context, tenantID uuid.UUID, input ListInput) ([]Guardian, int, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, tenant_id, name, phone, address, cpf, created_at, updated_at
		FROM guardians
		WHERE tenant_id = $1
		  AND (
		  	$2 = ''
		  	OR LOWER(name) LIKE '%' || LOWER($2) || '%'
		  	OR LOWER(phone) LIKE '%' || LOWER($2) || '%'
		  	OR LOWER(address) LIKE '%' || LOWER($2) || '%'
		  	OR LOWER(cpf) LIKE '%' || LOWER($2) || '%'
		  )
		ORDER BY created_at ASC
		LIMIT $3 OFFSET $4
	`, tenantID, input.Search, input.PerPage, input.Offset())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]Guardian, 0)
	for rows.Next() {
		item, err := scanGuardian(rows)
		if err != nil {
			return nil, 0, err
		}

		learnerIDs, err := r.loadLearnerIDsByGuardian(ctx, tenantID, item.ID)
		if err != nil {
			return nil, 0, err
		}
		item.LearnerIDs = learnerIDs
		items = append(items, *item)
	}

	countRow := r.pool.QueryRow(ctx, `
		SELECT COUNT(1)
		FROM guardians
		WHERE tenant_id = $1
		  AND (
		  	$2 = ''
		  	OR LOWER(name) LIKE '%' || LOWER($2) || '%'
		  	OR LOWER(phone) LIKE '%' || LOWER($2) || '%'
		  	OR LOWER(address) LIKE '%' || LOWER($2) || '%'
		  	OR LOWER(cpf) LIKE '%' || LOWER($2) || '%'
		  )
	`, tenantID, input.Search)

	var total int
	if err := countRow.Scan(&total); err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *PostgresRepository) Update(ctx context.Context, item *Guardian) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE guardians
		SET name = $1, phone = $2, address = $3, cpf = $4, updated_at = $5
		WHERE id = $6 AND tenant_id = $7
	`,
		item.Name,
		item.Phone,
		item.Address,
		item.CPF,
		item.UpdatedAt,
		item.ID,
		item.TenantID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return sharederrors.NotFound("GUARDIAN_NOT_FOUND", "guardian not found")
	}

	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, tenantID uuid.UUID, guardianID uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `
		DELETE FROM guardians
		WHERE id = $1 AND tenant_id = $2
	`, guardianID, tenantID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return sharederrors.NotFound("GUARDIAN_NOT_FOUND", "guardian not found")
	}

	return nil
}

func (r *PostgresRepository) EnsureByIDs(ctx context.Context, tenantID uuid.UUID, guardianIDs []uuid.UUID) error {
	for _, guardianID := range guardianIDs {
		var exists bool
		if err := r.pool.QueryRow(ctx, `
			SELECT EXISTS(
				SELECT 1 FROM guardians WHERE id = $1 AND tenant_id = $2
			)
		`, guardianID, tenantID).Scan(&exists); err != nil {
			return err
		}
		if !exists {
			return sharederrors.NotFound("GUARDIAN_NOT_FOUND", "guardian not found")
		}
	}

	return nil
}

func (r *PostgresRepository) ReplaceLearnerGuardians(ctx context.Context, tenantID uuid.UUID, learnerID uuid.UUID, guardianIDs []uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var learnerExists bool
	if err := tx.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM learners WHERE id = $1 AND tenant_id = $2
		)
	`, learnerID, tenantID).Scan(&learnerExists); err != nil {
		return err
	}
	if !learnerExists {
		return sharederrors.NotFound("LEARNER_NOT_FOUND", "learner not found")
	}

	for _, guardianID := range guardianIDs {
		var guardianExists bool
		if err := tx.QueryRow(ctx, `
			SELECT EXISTS(
				SELECT 1 FROM guardians WHERE id = $1 AND tenant_id = $2
			)
		`, guardianID, tenantID).Scan(&guardianExists); err != nil {
			return err
		}
		if !guardianExists {
			return sharederrors.NotFound("GUARDIAN_NOT_FOUND", "guardian not found")
		}
	}

	if _, err := tx.Exec(ctx, `
		DELETE FROM learner_guardians
		WHERE tenant_id = $1 AND learner_id = $2
	`, tenantID, learnerID); err != nil {
		return err
	}

	for _, guardianID := range guardianIDs {
		if _, err := tx.Exec(ctx, `
			INSERT INTO learner_guardians (
				tenant_id, learner_id, guardian_id, created_at
			) VALUES ($1,$2,$3,NOW())
		`, tenantID, learnerID, guardianID); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *PostgresRepository) CountLearnersByGuardian(ctx context.Context, tenantID uuid.UUID, guardianID uuid.UUID) (int, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT COUNT(1)
		FROM learner_guardians
		WHERE tenant_id = $1 AND guardian_id = $2
	`, tenantID, guardianID)

	var total int
	if err := row.Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

type guardianScanner interface {
	Scan(dest ...any) error
}

func scanGuardian(scanner guardianScanner) (*Guardian, error) {
	var item Guardian
	if err := scanner.Scan(
		&item.ID,
		&item.TenantID,
		&item.Name,
		&item.Phone,
		&item.Address,
		&item.CPF,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, sharederrors.NotFound("GUARDIAN_NOT_FOUND", "guardian not found")
		}
		return nil, err
	}

	return &item, nil
}

func (r *PostgresRepository) loadLearnerIDsByGuardian(ctx context.Context, tenantID uuid.UUID, guardianID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT learner_id
		FROM learner_guardians
		WHERE tenant_id = $1 AND guardian_id = $2
		ORDER BY learner_id ASC
	`, tenantID, guardianID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	learnerIDs := make([]uuid.UUID, 0)
	for rows.Next() {
		var learnerID uuid.UUID
		if err := rows.Scan(&learnerID); err != nil {
			return nil, err
		}
		learnerIDs = append(learnerIDs, learnerID)
	}

	return learnerIDs, nil
}
