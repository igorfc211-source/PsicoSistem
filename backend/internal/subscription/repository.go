package subscription

import (
	"context"

	"api-on/internal/shared/database"
	sharederrors "api-on/internal/shared/errors"

	"github.com/google/uuid"
)

type Repository interface {
	FindPlanBySlug(ctx context.Context, slug string) (*Plan, error)
	FindPlanByID(ctx context.Context, id uuid.UUID) (*Plan, error)
	Create(ctx context.Context, item *Subscription) error
	GetByTenantID(ctx context.Context, tenantID uuid.UUID) (*Subscription, *Plan, error)
	DeleteByTenantID(ctx context.Context, tenantID uuid.UUID) error
}

type JSONRepository struct {
	store *database.Store
}

func NewRepository(store *database.Store) *JSONRepository {
	return &JSONRepository{store: store}
}

func (r *JSONRepository) FindPlanBySlug(_ context.Context, slug string) (*Plan, error) {
	var result *Plan

	err := r.store.View(func(state database.State) error {
		record, exists := state.Plans[slug]
		if !exists {
			return sharederrors.NotFound("PLAN_NOT_FOUND", "plan not found")
		}

		plan, err := planFromRecord(record)
		if err != nil {
			return err
		}
		result = plan
		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *JSONRepository) FindPlanByID(_ context.Context, id uuid.UUID) (*Plan, error) {
	var result *Plan

	err := r.store.View(func(state database.State) error {
		for _, record := range state.Plans {
			if record.ID != id.String() {
				continue
			}

			plan, err := planFromRecord(record)
			if err != nil {
				return err
			}
			result = plan
			return nil
		}

		return sharederrors.NotFound("PLAN_NOT_FOUND", "plan not found")
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *JSONRepository) Create(_ context.Context, item *Subscription) error {
	return r.store.Update(func(state *database.State) error {
		state.Subscriptions[item.ID.String()] = database.SubscriptionRecord{
			ID:            item.ID.String(),
			TenantID:      item.TenantID.String(),
			PlanID:        item.PlanID.String(),
			Status:        item.Status,
			BillingCycle:  item.BillingCycle,
			AmountMonthly: item.AmountMonthly,
			StartedAt:     item.StartedAt,
			RenewalAt:     item.RenewalAt,
			CanceledAt:    item.CanceledAt,
			CreatedAt:     item.CreatedAt,
			UpdatedAt:     item.UpdatedAt,
		}
		return nil
	})
}

func (r *JSONRepository) GetByTenantID(_ context.Context, tenantID uuid.UUID) (*Subscription, *Plan, error) {
	var result *Subscription
	var plan *Plan

	err := r.store.View(func(state database.State) error {
		for _, record := range state.Subscriptions {
			if record.TenantID != tenantID.String() {
				continue
			}

			subscription, err := subscriptionFromRecord(record)
			if err != nil {
				return err
			}

			planRecord, exists := findPlanRecordByID(state.Plans, record.PlanID)
			if !exists {
				return sharederrors.Internal("subscription plan not found in state")
			}

			parsedPlan, err := planFromRecord(planRecord)
			if err != nil {
				return err
			}

			result = subscription
			plan = parsedPlan
			return nil
		}

		return sharederrors.NotFound("SUBSCRIPTION_NOT_FOUND", "subscription not found")
	})
	if err != nil {
		return nil, nil, err
	}

	return result, plan, nil
}

func (r *JSONRepository) DeleteByTenantID(_ context.Context, tenantID uuid.UUID) error {
	return r.store.Update(func(state *database.State) error {
		for id, record := range state.Subscriptions {
			if record.TenantID == tenantID.String() {
				delete(state.Subscriptions, id)
			}
		}
		return nil
	})
}

func findPlanRecordByID(plans map[string]database.PlanRecord, id string) (database.PlanRecord, bool) {
	for _, record := range plans {
		if record.ID == id {
			return record, true
		}
	}
	return database.PlanRecord{}, false
}

func planFromRecord(record database.PlanRecord) (*Plan, error) {
	id, err := uuid.Parse(record.ID)
	if err != nil {
		return nil, sharederrors.Internal("invalid plan record id")
	}

	return &Plan{
		ID:                id,
		Slug:              record.Slug,
		Name:              record.Name,
		PriceMonthlyCents: record.PriceMonthlyCents,
		HasTestsLibrary:   record.HasTestsLibrary,
		HasAI:             record.HasAI,
		HasGuardianPortal: record.HasGuardianPortal,
		MaxProfessionals:  record.MaxProfessionals,
		MaxPatients:       record.MaxPatients,
		Status:            record.Status,
		CreatedAt:         record.CreatedAt,
	}, nil
}

func subscriptionFromRecord(record database.SubscriptionRecord) (*Subscription, error) {
	id, err := uuid.Parse(record.ID)
	if err != nil {
		return nil, sharederrors.Internal("invalid subscription record id")
	}
	tenantID, err := uuid.Parse(record.TenantID)
	if err != nil {
		return nil, sharederrors.Internal("invalid subscription tenant id")
	}
	planID, err := uuid.Parse(record.PlanID)
	if err != nil {
		return nil, sharederrors.Internal("invalid subscription plan id")
	}

	return &Subscription{
		ID:            id,
		TenantID:      tenantID,
		PlanID:        planID,
		Status:        record.Status,
		BillingCycle:  record.BillingCycle,
		AmountMonthly: record.AmountMonthly,
		StartedAt:     record.StartedAt,
		RenewalAt:     record.RenewalAt,
		CanceledAt:    record.CanceledAt,
		CreatedAt:     record.CreatedAt,
		UpdatedAt:     record.UpdatedAt,
	}, nil
}
