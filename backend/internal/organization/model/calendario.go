package model

import (
	"time"

	"github.com/google/uuid"
)

type CalendarEvent struct {
	ID        string    `json:"id"`
	UserID    uuid.UUID    `json:"user_id"`
	PatientID *string   `json:"patient_id,omitempty"`
	StartAt   time.Time `json:"start_at"`
	EndAt     time.Time `json:"end_at"`
	Status    string    `json:"status"`
	Location  string    `json:"location"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}