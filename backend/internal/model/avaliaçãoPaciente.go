package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// > Representa uma avaliação/teste feito com algum instrumento
type Assessment struct {
	ID             uuid.UUID
	PatientID      uuid.UUID
	InstrumentType string		//tipo de instrumento
	AssessmentDate time.Time
	Hypotheses     string   	//hipóteses clínicas
	ResultSummary  string		//resumo do resultado
	Attachments    *json.RawMessage //anexos em JSON
	Status         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Results        []AssessmentResult //lista dos resultados detalhados
}

// item específico da avaliação
type AssessmentResult struct {
	ID           int64    `json:"id"`
	AssessmentID uuid.UUID `json:"assessment_id"`
	Title        string    `json:"title"`
	Value        string    `json:"value"`			//ex: abaixo do esperado
	Notes        string    `json:"notes"`			//ex: dificuldade ao falar coisas especificas
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
