package tenant

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

func (r *PostgresRepository) Create(ctx context.Context, item *Tenant) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO tenants (id, name, slug, cnpj, email, phone, status, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`,
		item.ID,
		item.Name,
		item.Slug,
		item.CNPJ,
		item.Email,
		item.Phone,
		item.Status,
		item.CreatedAt,
		item.UpdatedAt,
	)
	if err != nil {
		return sharederrors.Conflict("TENANT_SLUG_ALREADY_EXISTS", "clinic slug already exists")
	}

	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*Tenant, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, name, slug, cnpj, email, phone, status, created_at, updated_at
		FROM tenants
		WHERE id = $1
	`, id)

	var item Tenant
	if err := row.Scan(
		&item.ID,
		&item.Name,
		&item.Slug,
		&item.CNPJ,
		&item.Email,
		&item.Phone,
		&item.Status,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		return nil, sharederrors.NotFound("TENANT_NOT_FOUND", "tenant not found")
	}

	return &item, nil
}

func (r *PostgresRepository) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	row := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM tenants WHERE slug = $1)`, slug)

	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (r *PostgresRepository) ExistsByCNPJ(ctx context.Context, cnpj string) (bool, error) {
	row := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM tenants WHERE cnpj = $1)`, cnpj)

	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (r *PostgresRepository) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	row := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM tenants WHERE phone = $1)`, phone)

	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM tenants WHERE id = $1`, id)
	return err
}
