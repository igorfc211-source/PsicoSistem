package auth

import (
	"context"
	"time"

	"api-on/internal/shared/database"
	sharederrors "api-on/internal/shared/errors"

	"github.com/google/uuid"
)

type PasswordResetToken struct {
	ID        uuid.UUID
	TenantID  uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}

type PasswordResetRepository interface {
	Create(ctx context.Context, item *PasswordResetToken) error
	ConsumeValid(ctx context.Context, tokenHash string, now time.Time) (*PasswordResetToken, error)
}

type JSONPasswordResetRepository struct {
	store *database.Store
}

func NewPasswordResetRepository(store *database.Store) *JSONPasswordResetRepository {
	return &JSONPasswordResetRepository{store: store}
}

func (r *JSONPasswordResetRepository) Create(_ context.Context, item *PasswordResetToken) error {
	return r.store.Update(func(state *database.State) error {
		for key, record := range state.PasswordResetTokens {
			if record.UserID != item.UserID.String() || record.UsedAt != nil {
				continue
			}

			usedAt := item.CreatedAt
			record.UsedAt = &usedAt
			state.PasswordResetTokens[key] = record
		}

		state.PasswordResetTokens[item.ID.String()] = database.PasswordResetRecord{
			ID:        item.ID.String(),
			TenantID:  item.TenantID.String(),
			UserID:    item.UserID.String(),
			TokenHash: item.TokenHash,
			ExpiresAt: item.ExpiresAt,
			UsedAt:    item.UsedAt,
			CreatedAt: item.CreatedAt,
		}

		return nil
	})
}

func (r *JSONPasswordResetRepository) ConsumeValid(_ context.Context, tokenHash string, now time.Time) (*PasswordResetToken, error) {
	var result *PasswordResetToken

	err := r.store.Update(func(state *database.State) error {
		for key, record := range state.PasswordResetTokens {
			if record.TokenHash != tokenHash || record.UsedAt != nil || !record.ExpiresAt.After(now) {
				continue
			}

			usedAt := now
			record.UsedAt = &usedAt
			state.PasswordResetTokens[key] = record

			item, err := passwordResetFromRecord(record)
			if err != nil {
				return err
			}
			result = item
			return nil
		}

		return invalidPasswordResetToken()
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func passwordResetFromRecord(record database.PasswordResetRecord) (*PasswordResetToken, error) {
	id, err := uuid.Parse(record.ID)
	if err != nil {
		return nil, sharederrors.Internal("invalid password reset record id")
	}
	tenantID, err := uuid.Parse(record.TenantID)
	if err != nil {
		return nil, sharederrors.Internal("invalid password reset tenant id")
	}
	userID, err := uuid.Parse(record.UserID)
	if err != nil {
		return nil, sharederrors.Internal("invalid password reset user id")
	}

	return &PasswordResetToken{
		ID:        id,
		TenantID:  tenantID,
		UserID:    userID,
		TokenHash: record.TokenHash,
		ExpiresAt: record.ExpiresAt,
		UsedAt:    record.UsedAt,
		CreatedAt: record.CreatedAt,
	}, nil
}

func invalidPasswordResetToken() error {
	return sharederrors.Invalid("INVALID_PASSWORD_RESET_TOKEN", "invalid or expired password reset token")
}
