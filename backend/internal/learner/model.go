package learner

import (
	"time"

	"github.com/google/uuid"
)

const (
	StatusActive   = "active"
	StatusInactive = "inactive"
)

type Learner struct {
	ID                uuid.UUID   `json:"id"`
	TenantID          uuid.UUID   `json:"tenant_id"`
	Name              string      `json:"name"`
	PhotoURL          string      `json:"photo_url"`
	Gender            string      `json:"gender"`
	Guardian          string      `json:"guardian"`
	Age               string      `json:"age"`
	Status            string      `json:"status"`
	StartDate         string      `json:"start_date"`
	EndDate           string      `json:"end_date"`
	VisitCount        int         `json:"visit_count"`
	SessionPriceCents int64       `json:"session_price_cents"`
	GuardianIDs       []uuid.UUID `json:"guardian_ids"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
}
