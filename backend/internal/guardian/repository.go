package guardian

import (
	"context"
	"sort"
	"strings"
	"time"

	"api-on/internal/shared/database"
	sharederrors "api-on/internal/shared/errors"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, item *Guardian) error
	GetByIDAndTenant(ctx context.Context, tenantID uuid.UUID, guardianID uuid.UUID) (*Guardian, error)
	ListByTenant(ctx context.Context, tenantID uuid.UUID, input ListInput) ([]Guardian, int, error)
	Update(ctx context.Context, item *Guardian) error
	Delete(ctx context.Context, tenantID uuid.UUID, guardianID uuid.UUID) error
	EnsureByIDs(ctx context.Context, tenantID uuid.UUID, guardianIDs []uuid.UUID) error
	ReplaceLearnerGuardians(ctx context.Context, tenantID uuid.UUID, learnerID uuid.UUID, guardianIDs []uuid.UUID) error
	CountLearnersByGuardian(ctx context.Context, tenantID uuid.UUID, guardianID uuid.UUID) (int, error)
}

type JSONRepository struct {
	store *database.Store
}

func NewRepository(store *database.Store) *JSONRepository {
	return &JSONRepository{store: store}
}

func (r *JSONRepository) Create(_ context.Context, item *Guardian) error {
	return r.store.Update(func(state *database.State) error {
		state.Guardians[item.ID.String()] = toRecord(item)
		return nil
	})
}

func (r *JSONRepository) GetByIDAndTenant(_ context.Context, tenantID uuid.UUID, guardianID uuid.UUID) (*Guardian, error) {
	var result *Guardian

	err := r.store.View(func(state database.State) error {
		record, exists := state.Guardians[guardianID.String()]
		if !exists || record.TenantID != tenantID.String() {
			return sharederrors.NotFound("GUARDIAN_NOT_FOUND", "guardian not found")
		}

		item, err := fromRecord(record, learnerIDsForGuardian(state, tenantID.String(), guardianID.String()))
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

func (r *JSONRepository) ListByTenant(_ context.Context, tenantID uuid.UUID, input ListInput) ([]Guardian, int, error) {
	result := make([]Guardian, 0)

	err := r.store.View(func(state database.State) error {
		for _, record := range state.Guardians {
			if record.TenantID != tenantID.String() {
				continue
			}
			if input.Search != "" {
				search := strings.ToLower(input.Search)
				searchText := strings.ToLower(strings.Join([]string{
					record.Name,
					record.Phone,
					record.Address,
					record.CPF,
				}, " "))
				if !strings.Contains(searchText, search) {
					continue
				}
			}

			item, err := fromRecord(record, learnerIDsForGuardian(state, tenantID.String(), record.ID))
			if err != nil {
				return err
			}
			result = append(result, *item)
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
		return []Guardian{}, total, nil
	}

	end := start + input.PerPage
	if end > total {
		end = total
	}

	return result[start:end], total, nil
}

func (r *JSONRepository) Update(_ context.Context, item *Guardian) error {
	return r.store.Update(func(state *database.State) error {
		record, exists := state.Guardians[item.ID.String()]
		if !exists || record.TenantID != item.TenantID.String() {
			return sharederrors.NotFound("GUARDIAN_NOT_FOUND", "guardian not found")
		}

		state.Guardians[item.ID.String()] = toRecord(item)
		return nil
	})
}

func (r *JSONRepository) Delete(_ context.Context, tenantID uuid.UUID, guardianID uuid.UUID) error {
	return r.store.Update(func(state *database.State) error {
		record, exists := state.Guardians[guardianID.String()]
		if !exists || record.TenantID != tenantID.String() {
			return sharederrors.NotFound("GUARDIAN_NOT_FOUND", "guardian not found")
		}

		delete(state.Guardians, guardianID.String())
		return nil
	})
}

func (r *JSONRepository) EnsureByIDs(_ context.Context, tenantID uuid.UUID, guardianIDs []uuid.UUID) error {
	return r.store.View(func(state database.State) error {
		for _, guardianID := range guardianIDs {
			record, exists := state.Guardians[guardianID.String()]
			if !exists || record.TenantID != tenantID.String() {
				return sharederrors.NotFound("GUARDIAN_NOT_FOUND", "guardian not found")
			}
		}
		return nil
	})
}

func (r *JSONRepository) ReplaceLearnerGuardians(_ context.Context, tenantID uuid.UUID, learnerID uuid.UUID, guardianIDs []uuid.UUID) error {
	return r.store.Update(func(state *database.State) error {
		learnerRecord, exists := state.Learners[learnerID.String()]
		if !exists || learnerRecord.TenantID != tenantID.String() {
			return sharederrors.NotFound("LEARNER_NOT_FOUND", "learner not found")
		}

		for _, guardianID := range guardianIDs {
			record, exists := state.Guardians[guardianID.String()]
			if !exists || record.TenantID != tenantID.String() {
				return sharederrors.NotFound("GUARDIAN_NOT_FOUND", "guardian not found")
			}
		}

		for key, record := range state.LearnerGuardianLinks {
			if record.TenantID == tenantID.String() && record.LearnerID == learnerID.String() {
				delete(state.LearnerGuardianLinks, key)
			}
		}

		now := time.Now()
		for _, guardianID := range guardianIDs {
			key := database.LearnerGuardianLinkKey(learnerID.String(), guardianID.String())
			state.LearnerGuardianLinks[key] = database.LearnerGuardianRecord{
				TenantID:   tenantID.String(),
				LearnerID:  learnerID.String(),
				GuardianID: guardianID.String(),
				CreatedAt:  now,
			}
		}

		return nil
	})
}

func (r *JSONRepository) CountLearnersByGuardian(_ context.Context, tenantID uuid.UUID, guardianID uuid.UUID) (int, error) {
	total := 0

	err := r.store.View(func(state database.State) error {
		for _, record := range state.LearnerGuardianLinks {
			if record.TenantID == tenantID.String() && record.GuardianID == guardianID.String() {
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

func toRecord(item *Guardian) database.GuardianRecord {
	return database.GuardianRecord{
		ID:        item.ID.String(),
		TenantID:  item.TenantID.String(),
		Name:      item.Name,
		Phone:     item.Phone,
		Address:   item.Address,
		CPF:       item.CPF,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}

func fromRecord(record database.GuardianRecord, learnerIDs []uuid.UUID) (*Guardian, error) {
	id, err := uuid.Parse(record.ID)
	if err != nil {
		return nil, sharederrors.Internal("invalid guardian record id")
	}
	tenantID, err := uuid.Parse(record.TenantID)
	if err != nil {
		return nil, sharederrors.Internal("invalid guardian record tenant id")
	}

	return &Guardian{
		ID:         id,
		TenantID:   tenantID,
		Name:       record.Name,
		Phone:      record.Phone,
		Address:    record.Address,
		CPF:        record.CPF,
		LearnerIDs: learnerIDs,
		CreatedAt:  record.CreatedAt,
		UpdatedAt:  record.UpdatedAt,
	}, nil
}

func learnerIDsForGuardian(state database.State, tenantID string, guardianID string) []uuid.UUID {
	learnerIDs := make([]uuid.UUID, 0)
	for _, record := range state.LearnerGuardianLinks {
		if record.TenantID != tenantID || record.GuardianID != guardianID {
			continue
		}

		learnerID, err := uuid.Parse(record.LearnerID)
		if err != nil {
			continue
		}
		learnerIDs = append(learnerIDs, learnerID)
	}

	sort.Slice(learnerIDs, func(i, j int) bool {
		return learnerIDs[i].String() < learnerIDs[j].String()
	})

	return learnerIDs
}
