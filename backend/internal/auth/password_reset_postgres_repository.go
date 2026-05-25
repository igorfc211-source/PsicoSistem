package auth

import (
	"context"
	"time"

	sharederrors "api-on/internal/shared/errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresPasswordResetRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresPasswordResetRepository(pool *pgxpool.Pool) *PostgresPasswordResetRepository {
	return &PostgresPasswordResetRepository{pool: pool}
}

func (r *PostgresPasswordResetRepository) Create(ctx context.Context, item *PasswordResetToken) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `
		UPDATE password_reset_tokens
		SET used_at = $1
		WHERE user_id = $2 AND used_at IS NULL
	`, item.CreatedAt, item.UserID); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO password_reset_tokens (
			id, tenant_id, user_id, token_hash, expires_at, used_at, created_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7)
	`,
		item.ID,
		item.TenantID,
		item.UserID,
		item.TokenHash,
		item.ExpiresAt,
		item.UsedAt,
		item.CreatedAt,
	); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *PostgresPasswordResetRepository) ConsumeValid(ctx context.Context, tokenHash string, now time.Time) (*PasswordResetToken, error) {
	row := r.pool.QueryRow(ctx, `
		UPDATE password_reset_tokens
		SET used_at = $1
		WHERE token_hash = $2
		  AND used_at IS NULL
		  AND expires_at > $1
		RETURNING id, tenant_id, user_id, token_hash, expires_at, used_at, created_at
	`, now, tokenHash)

	item, err := scanPasswordResetToken(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, invalidPasswordResetToken()
		}
		return nil, err
	}

	return item, nil
}

type passwordResetScanner interface {
	Scan(dest ...any) error
}

func scanPasswordResetToken(scanner passwordResetScanner) (*PasswordResetToken, error) {
	var item PasswordResetToken
	if err := scanner.Scan(
		&item.ID,
		&item.TenantID,
		&item.UserID,
		&item.TokenHash,
		&item.ExpiresAt,
		&item.UsedAt,
		&item.CreatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, err
		}
		return nil, sharederrors.Internal("could not read password reset token")
	}

	return &item, nil
}
