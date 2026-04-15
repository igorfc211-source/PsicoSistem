package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

//model dos pacientes
type Patient struct {
	ID        	  uuid.UUID `json:"id"`
	Document  	  string    `json:"document"`
	UserID        uuid.UUID       `json:"user_id"`		         //por qual usuário pertence
	Sex 		  *string  			`json:"sex"`
	Name          string          `json:"name"`					//nome
	BirthDate     *time.Time      `json:"birth_date,omitempty"` //data de aniversário
	Guardians     json.RawMessage `json:"guardians,omitempty"`  //resposáveis
	School        string          `json:"school"`
	MainComplaint string          `json:"main_complaint"`	 	//queixa principal
	Notes         string          `json:"notes"`				//observações gerais
	Status        string          `json:"status"`				//em qual processo ele está (em andamento, arquivado, inativo...)
	Classroom     string           `json:"classroom_grade"` 	//ano letivo escolar
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}


type Session struct {
	ID              string          `json:"id"`
	PatientID       string          `json:"patient_id"` 				
	CalendarEventID *string         `json:"calendar_event_id,omitempty"`  //pode ligar com a agenda
	SessionAt       time.Time       `json:"session_at"`
	SessionType     string          `json:"session_type"`
	Summary         string          `json:"summary"`
	Conduct         string          `json:"conduct"`
	Referrals       json.RawMessage `json:"referrals,omitempty"`
	NextAction      string          `json:"next_action"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

