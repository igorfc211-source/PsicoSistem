package learner

import (
	"math"
	"strconv"
	"strings"
	"time"

	sharederrors "api-on/internal/shared/errors"

	"github.com/google/uuid"
)

type CreateInput struct {
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
	GeneralValueCents int64       `json:"general_value_cents"`
	GuardianIDs       []uuid.UUID `json:"guardian_ids"`
}

type UpdateInput struct {
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
	GeneralValueCents int64       `json:"general_value_cents"`
	GuardianIDs       []uuid.UUID `json:"guardian_ids"`
}

type ListInput struct {
	Status  string
	Search  string
	Page    int
	PerPage int
}

type Response struct {
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
	GeneralValueCents int64       `json:"general_value_cents"`
	GuardianIDs       []uuid.UUID `json:"guardian_ids"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
}

type ListMeta struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

func ToResponse(item *Learner) *Response {
	if item == nil {
		return nil
	}

	return &Response{
		ID:                item.ID,
		TenantID:          item.TenantID,
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
		GeneralValueCents: item.GeneralValueCents,
		GuardianIDs:       append([]uuid.UUID(nil), item.GuardianIDs...),
		CreatedAt:         item.CreatedAt,
		UpdatedAt:         item.UpdatedAt,
	}
}

func ParseListInput(pageValue string, perPageValue string, status string, search string) ListInput {
	page := parsePositiveInt(pageValue, 1)
	perPage := parsePositiveInt(perPageValue, 20)
	if perPage > 100 {
		perPage = 100
	}

	return ListInput{
		Status:  strings.TrimSpace(status),
		Search:  strings.TrimSpace(search),
		Page:    page,
		PerPage: perPage,
	}
}

func (i ListInput) Offset() int {
	return (i.Page - 1) * i.PerPage
}

func BuildListMeta(input ListInput, total int) ListMeta {
	totalPages := int(math.Ceil(float64(total) / float64(input.PerPage)))
	if totalPages == 0 {
		totalPages = 1
	}

	return ListMeta{
		Page:       input.Page,
		PerPage:    input.PerPage,
		Total:      total,
		TotalPages: totalPages,
	}
}

func ValidateStatus(status string) error {
	switch strings.TrimSpace(status) {
	case StatusActive, StatusInactive:
		return nil
	default:
		return sharederrors.Invalid("INVALID_LEARNER_STATUS", "invalid learner status")
	}
}

func parsePositiveInt(value string, fallback int) int {
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}
