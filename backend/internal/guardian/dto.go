package guardian

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CreateInput struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	CPF     string `json:"cpf"`
}

type UpdateInput struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	CPF     string `json:"cpf"`
}

type ListInput struct {
	Search  string
	Page    int
	PerPage int
}

type Response struct {
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

type ListMeta struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

func ToResponse(item *Guardian) *Response {
	if item == nil {
		return nil
	}

	return &Response{
		ID:         item.ID,
		TenantID:   item.TenantID,
		Name:       item.Name,
		Phone:      item.Phone,
		Address:    item.Address,
		CPF:        item.CPF,
		LearnerIDs: append([]uuid.UUID(nil), item.LearnerIDs...),
		CreatedAt:  item.CreatedAt,
		UpdatedAt:  item.UpdatedAt,
	}
}

func ParseListInput(pageValue string, perPageValue string, search string) ListInput {
	page := parsePositiveInt(pageValue, 1)
	perPage := parsePositiveInt(perPageValue, 20)
	if perPage > 100 {
		perPage = 100
	}

	return ListInput{
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

func parsePositiveInt(value string, fallback int) int {
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}
