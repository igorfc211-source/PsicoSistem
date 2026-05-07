package learner

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
	Create(ctx context.Context, item *Learner) error
	GetByIDAndTenant(ctx context.Context, tenantID uuid.UUID, learnerID uuid.UUID) (*Learner, error)
	ExistsByIDAndTenant(ctx context.Context, tenantID uuid.UUID, learnerID uuid.UUID) (bool, error)
	ListByTenant(ctx context.Context, tenantID uuid.UUID, input ListInput) ([]Learner, int, error)
	Update(ctx context.Context, item *Learner) error
	Deactivate(ctx context.Context, tenantID uuid.UUID, learnerID uuid.UUID) error
	CountActiveByTenant(ctx context.Context, tenantID uuid.UUID) (int, error)
}

type JSONRepository struct {
	store *database.Store
}

func NewRepository(store *database.Store) *JSONRepository {
	return &JSONRepository{store: store}
}

func (r *JSONRepository) Create(_ context.Context, item *Learner) error {
	return r.store.Update(func(state *database.State) error {
		state.Learners[item.ID.String()] = toRecord(item)
		return nil
	})
}

func (r *JSONRepository) GetByIDAndTenant(_ context.Context, tenantID uuid.UUID, learnerID uuid.UUID) (*Learner, error) {
	var result *Learner

	err := r.store.View(func(state database.State) error {
		record, exists := state.Learners[learnerID.String()]
		if !exists || record.TenantID != tenantID.String() {
			return sharederrors.NotFound("LEARNER_NOT_FOUND", "learner not found")
		}

		item, err := fromRecord(record, guardianIDsForLearner(state, tenantID.String(), learnerID.String()))
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

func (r *JSONRepository) ExistsByIDAndTenant(_ context.Context, tenantID uuid.UUID, learnerID uuid.UUID) (bool, error) {
	exists := false

	err := r.store.View(func(state database.State) error {
		record, found := state.Learners[learnerID.String()]
		exists = found && record.TenantID == tenantID.String()
		return nil
	})
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *JSONRepository) ListByTenant(_ context.Context, tenantID uuid.UUID, input ListInput) ([]Learner, int, error) {
	result := make([]Learner, 0)

	err := r.store.View(func(state database.State) error {
		for _, record := range state.Learners {
			if record.TenantID != tenantID.String() {
				continue
			}
			if input.Status != "" && record.Status != input.Status {
				continue
			}
			if input.Search != "" {
				search := strings.ToLower(input.Search)
				searchText := strings.ToLower(strings.Join([]string{
					record.Name,
					record.Guardian,
					record.Gender,
					record.Age,
				}, " "))
				if !strings.Contains(searchText, search) {
					continue
				}
			}

			item, err := fromRecord(record, guardianIDsForLearner(state, tenantID.String(), record.ID))
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
		return []Learner{}, total, nil
	}

	end := start + input.PerPage
	if end > total {
		end = total
	}

	return result[start:end], total, nil
}

func (r *JSONRepository) Update(_ context.Context, item *Learner) error {
	return r.store.Update(func(state *database.State) error {
		record, exists := state.Learners[item.ID.String()]
		if !exists || record.TenantID != item.TenantID.String() {
			return sharederrors.NotFound("LEARNER_NOT_FOUND", "learner not found")
		}

		state.Learners[item.ID.String()] = toRecord(item)
		return nil
	})
}

func (r *JSONRepository) Deactivate(_ context.Context, tenantID uuid.UUID, learnerID uuid.UUID) error {
	return r.store.Update(func(state *database.State) error {
		record, exists := state.Learners[learnerID.String()]
		if !exists || record.TenantID != tenantID.String() {
			return sharederrors.NotFound("LEARNER_NOT_FOUND", "learner not found")
		}

		record.Status = StatusInactive
		record.UpdatedAt = time.Now()
		state.Learners[learnerID.String()] = record
		return nil
	})
}

func (r *JSONRepository) CountActiveByTenant(_ context.Context, tenantID uuid.UUID) (int, error) {
	total := 0

	err := r.store.View(func(state database.State) error {
		for _, record := range state.Learners {
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

func toRecord(item *Learner) database.LearnerRecord {
	return database.LearnerRecord{
		ID:                item.ID.String(),
		TenantID:          item.TenantID.String(),
		Name:              item.Name,
		PhotoURL:          item.PhotoURL,
		Gender:            item.Gender,
		Guardian:          item.Guardian,
		Age:               item.Age,
		Status:            item.Status,
		StartDate:         item.StartDate,
		EndDate:           item.EndDate,
		VisitCount:        item.VisitCount,
		SessionPriceCents: item.SessionPriceCents,
		CreatedAt:         item.CreatedAt,
		UpdatedAt:         item.UpdatedAt,
	}
}

func fromRecord(record database.LearnerRecord, guardianIDs []uuid.UUID) (*Learner, error) {
	id, err := uuid.Parse(record.ID)
	if err != nil {
		return nil, sharederrors.Internal("invalid learner record id")
	}
	tenantID, err := uuid.Parse(record.TenantID)
	if err != nil {
		return nil, sharederrors.Internal("invalid learner record tenant id")
	}

	return &Learner{
		ID:                id,
		TenantID:          tenantID,
		Name:              record.Name,
		PhotoURL:          record.PhotoURL,
		Gender:            record.Gender,
		Guardian:          record.Guardian,
		Age:               record.Age,
		Status:            record.Status,
		StartDate:         record.StartDate,
		EndDate:           record.EndDate,
		VisitCount:        record.VisitCount,
		SessionPriceCents: record.SessionPriceCents,
		GuardianIDs:       guardianIDs,
		CreatedAt:         record.CreatedAt,
		UpdatedAt:         record.UpdatedAt,
	}, nil
}

func guardianIDsForLearner(state database.State, tenantID string, learnerID string) []uuid.UUID {
	guardianIDs := make([]uuid.UUID, 0)
	for _, record := range state.LearnerGuardianLinks {
		if record.TenantID != tenantID || record.LearnerID != learnerID {
			continue
		}

		guardianID, err := uuid.Parse(record.GuardianID)
		if err != nil {
			continue
		}
		guardianIDs = append(guardianIDs, guardianID)
	}

	sort.Slice(guardianIDs, func(i, j int) bool {
		return guardianIDs[i].String() < guardianIDs[j].String()
	})

	return guardianIDs
}
