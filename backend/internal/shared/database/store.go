package database

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"api-on/internal/shared/permissions"
)

type PlanRecord struct {
	ID                string    `json:"id"`
	Slug              string    `json:"slug"`
	Name              string    `json:"name"`
	PriceMonthlyCents int64     `json:"price_monthly_cents"`
	HasTestsLibrary   bool      `json:"has_tests_library"`
	HasAI             bool      `json:"has_ai"`
	HasGuardianPortal bool      `json:"has_guardian_portal"`
	MaxProfessionals  int       `json:"max_professionals"`
	MaxPatients       int       `json:"max_patients"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
}

type TenantRecord struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CNPJ      string    `json:"cnpj"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRecord struct {
	ID           string                         `json:"id"`
	TenantID     string                         `json:"tenant_id"`
	Name         string                         `json:"name"`
	Email        string                         `json:"email"`
	PasswordHash string                         `json:"password_hash"`
	Role         string                         `json:"role"`
	Status       string                         `json:"status"`
	Permissions  permissions.AccountPermissions `json:"permissions"`
	LastLoginAt  *time.Time                     `json:"last_login_at,omitempty"`
	CreatedAt    time.Time                      `json:"created_at"`
	UpdatedAt    time.Time                      `json:"updated_at"`
}

type SubscriptionRecord struct {
	ID            string     `json:"id"`
	TenantID      string     `json:"tenant_id"`
	PlanID        string     `json:"plan_id"`
	Status        string     `json:"status"`
	BillingCycle  string     `json:"billing_cycle"`
	AmountMonthly int64      `json:"amount_monthly"`
	StartedAt     time.Time  `json:"started_at"`
	RenewalAt     *time.Time `json:"renewal_at,omitempty"`
	CanceledAt    *time.Time `json:"canceled_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type LearnerRecord struct {
	ID                string    `json:"id"`
	TenantID          string    `json:"tenant_id"`
	Name              string    `json:"name"`
	PhotoURL          string    `json:"photo_url"`
	Gender            string    `json:"gender"`
	Guardian          string    `json:"guardian"`
	Age               string    `json:"age"`
	Status            string    `json:"status"`
	StartDate         string    `json:"start_date"`
	EndDate           string    `json:"end_date"`
	VisitCount        int       `json:"visit_count"`
	SessionPriceCents int64     `json:"session_price_cents"`
	GeneralValueCents int64     `json:"general_value_cents"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type GuardianRecord struct {
	ID           string    `json:"id"`
	TenantID     string    `json:"tenant_id"`
	Name         string    `json:"name"`
	Relationship string    `json:"relationship"`
	Phone        string    `json:"phone"`
	Address      string    `json:"address"`
	CPF          string    `json:"cpf,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type LearnerGuardianRecord struct {
	TenantID   string    `json:"tenant_id"`
	LearnerID  string    `json:"learner_id"`
	GuardianID string    `json:"guardian_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type State struct {
	Plans                map[string]PlanRecord            `json:"plans"`
	Tenants              map[string]TenantRecord          `json:"tenants"`
	Users                map[string]UserRecord            `json:"users"`
	Subscriptions        map[string]SubscriptionRecord    `json:"subscriptions"`
	Learners             map[string]LearnerRecord         `json:"learners"`
	Guardians            map[string]GuardianRecord        `json:"guardians"`
	LearnerGuardianLinks map[string]LearnerGuardianRecord `json:"learner_guardian_links"`
}

type Store struct {
	path string
	mu   sync.Mutex
}

func NewStore(path string) *Store {
	return &Store{path: path}
}

func (s *Store) Initialize() error {
	return s.Update(func(state *State) error {
		return nil
	})
}

func (s *Store) View(fn func(state State) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	state, err := s.load()
	if err != nil {
		return err
	}

	return fn(state)
}

func (s *Store) Update(fn func(state *State) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	state, err := s.load()
	if err != nil {
		return err
	}

	if err := fn(&state); err != nil {
		return err
	}

	return s.save(state)
}

func (s *Store) load() (State, error) {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return State{}, fmt.Errorf("create data directory: %w", err)
	}

	file, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return seedState(), nil
		}
		return State{}, fmt.Errorf("read state file: %w", err)
	}

	if len(file) == 0 {
		return seedState(), nil
	}

	state := seedState()
	if err := json.Unmarshal(file, &state); err != nil {
		return State{}, fmt.Errorf("decode state file: %w", err)
	}

	ensureMaps(&state)
	seedPlans(&state)
	return state, nil
}

func (s *Store) save(state State) error {
	ensureMaps(&state)
	seedPlans(&state)

	payload, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("encode state file: %w", err)
	}

	tempPath := s.path + ".tmp"
	if err := os.WriteFile(tempPath, payload, 0o644); err != nil {
		return fmt.Errorf("write temp state file: %w", err)
	}

	if _, err := os.Stat(s.path); err == nil {
		if err := os.Remove(s.path); err != nil {
			_ = os.Remove(tempPath)
			return fmt.Errorf("remove previous state file: %w", err)
		}
	} else if !os.IsNotExist(err) {
		_ = os.Remove(tempPath)
		return fmt.Errorf("stat previous state file: %w", err)
	}

	if err := os.Rename(tempPath, s.path); err != nil {
		return fmt.Errorf("rename temp state file: %w", err)
	}

	return nil
}

func seedState() State {
	state := State{
		Plans:                make(map[string]PlanRecord),
		Tenants:              make(map[string]TenantRecord),
		Users:                make(map[string]UserRecord),
		Subscriptions:        make(map[string]SubscriptionRecord),
		Learners:             make(map[string]LearnerRecord),
		Guardians:            make(map[string]GuardianRecord),
		LearnerGuardianLinks: make(map[string]LearnerGuardianRecord),
	}
	seedPlans(&state)
	return state
}

func ensureMaps(state *State) {
	if state.Plans == nil {
		state.Plans = make(map[string]PlanRecord)
	}
	if state.Tenants == nil {
		state.Tenants = make(map[string]TenantRecord)
	}
	if state.Users == nil {
		state.Users = make(map[string]UserRecord)
	}
	if state.Subscriptions == nil {
		state.Subscriptions = make(map[string]SubscriptionRecord)
	}
	if state.Learners == nil {
		state.Learners = make(map[string]LearnerRecord)
	}
	if state.Guardians == nil {
		state.Guardians = make(map[string]GuardianRecord)
	}
	if state.LearnerGuardianLinks == nil {
		state.LearnerGuardianLinks = make(map[string]LearnerGuardianRecord)
	}
}

func LearnerGuardianLinkKey(learnerID string, guardianID string) string {
	return learnerID + ":" + guardianID
}

func seedPlans(state *State) {
	now := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	defaults := []PlanRecord{
		{
			ID:                "00000000-0000-0000-0000-000000000101",
			Slug:              "basico",
			Name:              "Básico",
			PriceMonthlyCents: 9700,
			HasTestsLibrary:   false,
			HasAI:             false,
			HasGuardianPortal: false,
			MaxProfessionals:  3,
			MaxPatients:       100,
			Status:            "active",
			CreatedAt:         now,
		},
		{
			ID:                "00000000-0000-0000-0000-000000000102",
			Slug:              "intermediario",
			Name:              "Intermediário",
			PriceMonthlyCents: 14700,
			HasTestsLibrary:   true,
			HasAI:             false,
			HasGuardianPortal: true,
			MaxProfessionals:  8,
			MaxPatients:       300,
			Status:            "active",
			CreatedAt:         now,
		},
		{
			ID:                "00000000-0000-0000-0000-000000000103",
			Slug:              "premium",
			Name:              "Premium",
			PriceMonthlyCents: 19700,
			HasTestsLibrary:   true,
			HasAI:             true,
			HasGuardianPortal: true,
			MaxProfessionals:  20,
			MaxPatients:       1000,
			Status:            "active",
			CreatedAt:         now,
		},
	}

	for _, plan := range defaults {
		if _, exists := state.Plans[plan.Slug]; !exists {
			state.Plans[plan.Slug] = plan
		}
	}
}
