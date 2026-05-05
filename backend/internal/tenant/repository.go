package tenant

import (
	"context"
	"strings"

	"api-on/internal/shared/database"
	sharederrors "api-on/internal/shared/errors"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, item *Tenant) error
	GetByID(ctx context.Context, id uuid.UUID) (*Tenant, error)
	ExistsBySlug(ctx context.Context, slug string) (bool, error)
	ExistsByCNPJ(ctx context.Context, cnpj string) (bool, error)
	ExistsByPhone(ctx context.Context, phone string) (bool, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type JSONRepository struct {
	store *database.Store
}

func NewRepository(store *database.Store) *JSONRepository {
	return &JSONRepository{store: store}
}

func (r *JSONRepository) Create(_ context.Context, item *Tenant) error {
	return r.store.Update(func(state *database.State) error {
		for _, record := range state.Tenants {
			if strings.EqualFold(record.Slug, item.Slug) {
				return sharederrors.Conflict("TENANT_SLUG_ALREADY_EXISTS", "clinic slug already exists")
			}
			if item.CNPJ != "" && record.CNPJ == item.CNPJ {
				return sharederrors.Conflict("TENANT_DOCUMENT_ALREADY_EXISTS", "cpf_cnpj already registered")
			}
			if record.Phone == item.Phone {
				return sharederrors.Conflict("TENANT_PHONE_ALREADY_EXISTS", "phone already registered")
			}
		}

		state.Tenants[item.ID.String()] = database.TenantRecord{
			ID:        item.ID.String(),
			Name:      item.Name,
			Slug:      item.Slug,
			CNPJ:      item.CNPJ,
			Email:     item.Email,
			Phone:     item.Phone,
			Status:    item.Status,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}
		return nil
	})
}

func (r *JSONRepository) GetByID(_ context.Context, id uuid.UUID) (*Tenant, error) {
	var result *Tenant

	err := r.store.View(func(state database.State) error {
		record, exists := state.Tenants[id.String()]
		if !exists {
			return sharederrors.NotFound("TENANT_NOT_FOUND", "tenant not found")
		}

		item, err := fromRecord(record)
		if err != nil {
			return err
		}
		result = item
		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *JSONRepository) ExistsBySlug(_ context.Context, slug string) (bool, error) {
	exists := false

	err := r.store.View(func(state database.State) error {
		for _, record := range state.Tenants {
			if strings.EqualFold(record.Slug, slug) {
				exists = true
				break
			}
		}
		return nil
	})

	return exists, err
}

func (r *JSONRepository) ExistsByCNPJ(_ context.Context, cnpj string) (bool, error) {
	exists := false

	err := r.store.View(func(state database.State) error {
		for _, record := range state.Tenants {
			if record.CNPJ == cnpj {
				exists = true
				break
			}
		}
		return nil
	})

	return exists, err
}

func (r *JSONRepository) ExistsByPhone(_ context.Context, phone string) (bool, error) {
	exists := false

	err := r.store.View(func(state database.State) error {
		for _, record := range state.Tenants {
			if record.Phone == phone {
				exists = true
				break
			}
		}
		return nil
	})

	return exists, err
}

func (r *JSONRepository) Delete(_ context.Context, id uuid.UUID) error {
	return r.store.Update(func(state *database.State) error {
		delete(state.Tenants, id.String())
		return nil
	})
}

func fromRecord(record database.TenantRecord) (*Tenant, error) {
	id, err := uuid.Parse(record.ID)
	if err != nil {
		return nil, sharederrors.Internal("invalid tenant record id")
	}

	return &Tenant{
		ID:        id,
		Name:      record.Name,
		Slug:      record.Slug,
		CNPJ:      record.CNPJ,
		Email:     record.Email,
		Phone:     record.Phone,
		Status:    record.Status,
		CreatedAt: record.CreatedAt,
		UpdatedAt: record.UpdatedAt,
	}, nil
}
