package guardian

import (
	"time"

	"github.com/google/uuid"
)

type Guardian struct {
	ID         uuid.UUID   `json:"id"`
	TenantID   uuid.UUID   `json:"tenant_id"`
	Name       string      `json:"name"`
	Phone      string      `json:"phone"`
	Address    string      `json:"address"`
	CPF        string      `json:"cpf,omitempty"`
	LearnerIDs []uuid.UUID `json:"learner_ids"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}
