package model

import "time"

type Invoice struct {
	ID            string
	PatientID     string
	UserID        string
	Type          string
	Amount        float64
	DueDate       time.Time
	Status        string
	PaymentMethod string
	PaidAt        *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}