package user

import (
	"context"
	"sort"
	"strings"
	"time"

	"api-on/internal/shared/database"
	sharederrors "api-on/internal/shared/errors"
	sharedvalidator "api-on/internal/shared/validator"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, item *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	GetByIDAndTenant(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) (*User, error)
	ListByTenant(ctx context.Context, tenantID uuid.UUID, input ListInput) ([]User, int, error)
	Update(ctx context.Context, item *User) error
	Deactivate(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) error
	DeleteByID(ctx context.Context, userID uuid.UUID) error
	CountActiveByTenant(ctx context.Context, tenantID uuid.UUID) (int, error)
	TouchLastLogin(ctx context.Context, userID uuid.UUID, at time.Time) error
}

type JSONRepository struct {
	store *database.Store
}

func NewRepository(store *database.Store) *JSONRepository {
	return &JSONRepository{store: store}
}

func (r *JSONRepository) Create(_ context.Context, item *User) error {
	return r.store.Update(func(state *database.State) error {
		for _, record := range state.Users {
			if strings.EqualFold(record.Email, item.Email) {
				return sharederrors.Conflict("USER_EMAIL_ALREADY_EXISTS", "email already registered")
			}
		}

		state.Users[item.ID.String()] = database.UserRecord{
			ID:           item.ID.String(),
			TenantID:     item.TenantID.String(),
			Name:         item.Name,
			Email:        item.Email,
			PasswordHash: item.PasswordHash,
			Role:         item.Role,
			Status:       item.Status,
			Permissions:  item.Permissions,
			LastLoginAt:  item.LastLoginAt,
			CreatedAt:    item.CreatedAt,
			UpdatedAt:    item.UpdatedAt,
		}
		return nil
	})
}

func (r *JSONRepository) FindByEmail(_ context.Context, email string) (*User, error) {
	var result *User
	email = sharedvalidator.NormalizeEmail(email)

	err := r.store.View(func(state database.State) error {
		for _, record := range state.Users {
			if !strings.EqualFold(record.Email, email) {
				continue
			}

			user, err := fromRecord(record)
			if err != nil {
				return err
			}
			result = user
			return nil
		}

		return sharederrors.NotFound("USER_NOT_FOUND", "user not found")
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *JSONRepository) GetByIDAndTenant(_ context.Context, tenantID uuid.UUID, userID uuid.UUID) (*User, error) {
	var result *User

	err := r.store.View(func(state database.State) error {
		record, exists := state.Users[userID.String()]
		if !exists || record.TenantID != tenantID.String() {
			return sharederrors.NotFound("USER_NOT_FOUND", "user not found")
		}

		user, err := fromRecord(record)
		if err != nil {
			return err
		}
		result = user
		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *JSONRepository) ListByTenant(_ context.Context, tenantID uuid.UUID, input ListInput) ([]User, int, error) {
	var result []User

	err := r.store.View(func(state database.State) error {
		for _, record := range state.Users {
			if record.TenantID != tenantID.String() {
				continue
			}
			if input.Role != "" && record.Role != input.Role {
				continue
			}
			if input.Status != "" && record.Status != input.Status {
				continue
			}
			if input.Search != "" {
				search := strings.ToLower(input.Search)
				if !strings.Contains(strings.ToLower(record.Name), search) && !strings.Contains(strings.ToLower(record.Email), search) {
					continue
				}
			}

			user, err := fromRecord(record)
			if err != nil {
				return err
			}
			result = append(result, *user)
		}
		return nil
	})
	if err != nil {
		return nil, 0, err
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})

	total := len(result)
	start := input.Offset()
	if start >= total {
		return []User{}, total, nil
	}

	end := start + input.PerPage
	if end > total {
		end = total
	}

	return result[start:end], total, nil
}

func (r *JSONRepository) Update(_ context.Context, item *User) error {
	return r.store.Update(func(state *database.State) error {
		record, exists := state.Users[item.ID.String()]
		if !exists || record.TenantID != item.TenantID.String() {
			return sharederrors.NotFound("USER_NOT_FOUND", "user not found")
		}

		for _, existing := range state.Users {
			if existing.ID == item.ID.String() {
				continue
			}
			if strings.EqualFold(existing.Email, item.Email) {
				return sharederrors.Conflict("USER_EMAIL_ALREADY_EXISTS", "email already registered")
			}
		}

		state.Users[item.ID.String()] = database.UserRecord{
			ID:           item.ID.String(),
			TenantID:     item.TenantID.String(),
			Name:         item.Name,
			Email:        item.Email,
			PasswordHash: item.PasswordHash,
			Role:         item.Role,
			Status:       item.Status,
			Permissions:  item.Permissions,
			LastLoginAt:  item.LastLoginAt,
			CreatedAt:    item.CreatedAt,
			UpdatedAt:    item.UpdatedAt,
		}
		return nil
	})
}

func (r *JSONRepository) Deactivate(_ context.Context, tenantID uuid.UUID, userID uuid.UUID) error {
	return r.store.Update(func(state *database.State) error {
		record, exists := state.Users[userID.String()]
		if !exists || record.TenantID != tenantID.String() {
			return sharederrors.NotFound("USER_NOT_FOUND", "user not found")
		}

		record.Status = StatusInactive
		record.UpdatedAt = time.Now()
		state.Users[userID.String()] = record
		return nil
	})
}

func (r *JSONRepository) DeleteByID(_ context.Context, userID uuid.UUID) error {
	return r.store.Update(func(state *database.State) error {
		delete(state.Users, userID.String())
		return nil
	})
}

func (r *JSONRepository) CountActiveByTenant(_ context.Context, tenantID uuid.UUID) (int, error) {
	total := 0

	err := r.store.View(func(state database.State) error {
		for _, record := range state.Users {
			if record.TenantID == tenantID.String() && record.Status == StatusActive {
				total++
			}
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *JSONRepository) TouchLastLogin(_ context.Context, userID uuid.UUID, at time.Time) error {
	return r.store.Update(func(state *database.State) error {
		record, exists := state.Users[userID.String()]
		if !exists {
			return sharederrors.NotFound("USER_NOT_FOUND", "user not found")
		}

		record.LastLoginAt = &at
		record.UpdatedAt = at
		state.Users[userID.String()] = record
		return nil
	})
}

func fromRecord(record database.UserRecord) (*User, error) {
	id, err := uuid.Parse(record.ID)
	if err != nil {
		return nil, sharederrors.Internal("invalid user record id")
	}
	tenantID, err := uuid.Parse(record.TenantID)
	if err != nil {
		return nil, sharederrors.Internal("invalid user record tenant id")
	}

	return &User{
		ID:           id,
		TenantID:     tenantID,
		Name:         record.Name,
		Email:        record.Email,
		PasswordHash: record.PasswordHash,
		Role:         record.Role,
		Status:       record.Status,
		Permissions:  record.Permissions,
		LastLoginAt:  record.LastLoginAt,
		CreatedAt:    record.CreatedAt,
		UpdatedAt:    record.UpdatedAt,
	}, nil
}
