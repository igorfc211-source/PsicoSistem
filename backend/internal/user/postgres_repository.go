package user

import (
	"context"
	"encoding/json"
	"time"

	sharederrors "api-on/internal/shared/errors"
	"api-on/internal/shared/permissions"
	sharedvalidator "api-on/internal/shared/validator"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(ctx context.Context, item *User) error {
	permissionsPayload, err := json.Marshal(item.Permissions)
	if err != nil {
		return sharederrors.Internal("could not encode account permissions")
	}

	_, err = r.pool.Exec(ctx, `
		INSERT INTO users (
			id, tenant_id, name, email, password_hash, role, status, permissions,
			last_login_at, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`,
		item.ID,
		item.TenantID,
		item.Name,
		item.Email,
		item.PasswordHash,
		item.Role,
		item.Status,
		permissionsPayload,
		item.LastLoginAt,
		item.CreatedAt,
		item.UpdatedAt,
	)
	if err != nil {
		return sharederrors.Conflict("USER_EMAIL_ALREADY_EXISTS", "email already registered")
	}

	return nil
}

func (r *PostgresRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, name, email, password_hash, role, status, permissions,
		       last_login_at, created_at, updated_at
		FROM users
		WHERE email = $1
	`, sharedvalidator.NormalizeEmail(email))

	return scanUser(row)
}

func (r *PostgresRepository) GetByIDAndTenant(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) (*User, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, name, email, password_hash, role, status, permissions,
		       last_login_at, created_at, updated_at
		FROM users
		WHERE id = $1 AND tenant_id = $2
	`, userID, tenantID)

	return scanUser(row)
}

func (r *PostgresRepository) ListByTenant(ctx context.Context, tenantID uuid.UUID, input ListInput) ([]User, int, error) {
	filteredRows, err := r.pool.Query(ctx, `
		SELECT id, tenant_id, name, email, password_hash, role, status, permissions,
		       last_login_at, created_at, updated_at
		FROM users
		WHERE tenant_id = $1
		  AND ($2 = '' OR role = $2)
		  AND ($3 = '' OR status = $3)
		  AND ($4 = '' OR LOWER(name) LIKE '%' || LOWER($4) || '%' OR LOWER(email) LIKE '%' || LOWER($4) || '%')
		ORDER BY created_at ASC
		LIMIT $5 OFFSET $6
	`, tenantID, input.Role, input.Status, input.Search, input.PerPage, input.Offset())
	if err != nil {
		return nil, 0, err
	}
	defer filteredRows.Close()

	items := make([]User, 0)
	for filteredRows.Next() {
		item, err := scanUser(filteredRows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, *item)
	}

	countRow := r.pool.QueryRow(ctx, `
		SELECT COUNT(1)
		FROM users
		WHERE tenant_id = $1
		  AND ($2 = '' OR role = $2)
		  AND ($3 = '' OR status = $3)
		  AND ($4 = '' OR LOWER(name) LIKE '%' || LOWER($4) || '%' OR LOWER(email) LIKE '%' || LOWER($4) || '%')
	`, tenantID, input.Role, input.Status, input.Search)

	var total int
	if err := countRow.Scan(&total); err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *PostgresRepository) Update(ctx context.Context, item *User) error {
	permissionsPayload, err := json.Marshal(item.Permissions)
	if err != nil {
		return sharederrors.Internal("could not encode account permissions")
	}

	tag, err := r.pool.Exec(ctx, `
		UPDATE users
		SET name = $1, email = $2, password_hash = $3, role = $4, status = $5,
		    permissions = $6, last_login_at = $7, created_at = $8, updated_at = $9
		WHERE id = $10 AND tenant_id = $11
	`,
		item.Name,
		item.Email,
		item.PasswordHash,
		item.Role,
		item.Status,
		permissionsPayload,
		item.LastLoginAt,
		item.CreatedAt,
		item.UpdatedAt,
		item.ID,
		item.TenantID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return sharederrors.NotFound("USER_NOT_FOUND", "user not found")
	}

	return nil
}

func (r *PostgresRepository) Deactivate(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE users
		SET status = $1, updated_at = NOW()
		WHERE id = $2 AND tenant_id = $3
	`, StatusInactive, userID, tenantID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return sharederrors.NotFound("USER_NOT_FOUND", "user not found")
	}

	return nil
}

func (r *PostgresRepository) DeleteByID(ctx context.Context, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
	return err
}

func (r *PostgresRepository) CountActiveByTenant(ctx context.Context, tenantID uuid.UUID) (int, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT COUNT(1)
		FROM users
		WHERE tenant_id = $1 AND status = $2
	`, tenantID, StatusActive)

	var total int
	if err := row.Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

func (r *PostgresRepository) TouchLastLogin(ctx context.Context, userID uuid.UUID, at time.Time) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE users
		SET last_login_at = $1, updated_at = $1
		WHERE id = $2
	`, at, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return sharederrors.NotFound("USER_NOT_FOUND", "user not found")
	}

	return nil
}

type userScanner interface {
	Scan(dest ...any) error
}

func scanUser(scanner userScanner) (*User, error) {
	var item User
	var permissionsPayload []byte
	if err := scanner.Scan(
		&item.ID,
		&item.TenantID,
		&item.Name,
		&item.Email,
		&item.PasswordHash,
		&item.Role,
		&item.Status,
		&permissionsPayload,
		&item.LastLoginAt,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		return nil, sharederrors.NotFound("USER_NOT_FOUND", "user not found")
	}

	if len(permissionsPayload) == 0 {
		item.Permissions = permissions.DefaultForRole(item.Role)
	} else if err := json.Unmarshal(permissionsPayload, &item.Permissions); err != nil {
		return nil, sharederrors.Internal("could not decode account permissions")
	}

	return &item, nil
}
