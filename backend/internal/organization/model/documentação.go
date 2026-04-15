package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type DocumentTemplate struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`		//utiliza o template do usuário X
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

//Ele segue estritamente o document Template dentro  <<<
type Document struct {
	ID          uuid.UUID     `json:"id"`
	PatientID   string     `json:"patient_id"`				//pertence a esse paciente
	TemplateID  *string    `json:"template_id,omitempty"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	Status      string     `json:"status"`
	GeneratedAt *time.Time `json:"generated_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}


//histórico(máquina do tempo)

type TimelineItem struct {
	Type      string
	Date      time.Time
	Reference string
	Payload   map[string]any
}


type InterventionPlan struct {
	ID                 string
	PatientID          string
	GeneralObjective   string
	SpecificObjectives json.RawMessage
	Actions            json.RawMessage
	Frequency          string
	Status             string
	StartDate          *time.Time
	EndDate            *time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
}