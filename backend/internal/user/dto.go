package user

import (
	"math"
	"strconv"
	"strings"
	"time"

	sharederrors "api-on/internal/shared/errors"

	"github.com/google/uuid"
)

type CreateInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}

type UpdateInput struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Status string `json:"status"`
}

type ListInput struct {
	Role    string
	Status  string
	Search  string
	Page    int
	PerPage int
}

type Response struct {
	ID          uuid.UUID  `json:"user_id"`
	TenantID    uuid.UUID  `json:"tenant_id"`
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	Role        string     `json:"role"`
	Status      string     `json:"status"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type ListMeta struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

func ToResponse(item *User) *Response {
	if item == nil {
		return nil
	}

	return &Response{
		ID:          item.ID,
		TenantID:    item.TenantID,
		Name:        item.Name,
		Email:       item.Email,
		Role:        item.Role,
		Status:      item.Status,
		LastLoginAt: item.LastLoginAt,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}

func ParseListInput(pageValue string, perPageValue string, role string, status string, search string) ListInput {
	page := parsePositiveInt(pageValue, 1)
	perPage := parsePositiveInt(perPageValue, 20)
	if perPage > 100 {
		perPage = 100
	}

	return ListInput{
		Role:    strings.TrimSpace(role),
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

func ValidateRole(role string) error {
	switch strings.TrimSpace(role) {
	case RoleOwner, RoleAdmin, RoleCoordinator, RoleProfessional, RoleFinancial:
		return nil
	default:
		return sharederrors.Invalid("INVALID_ROLE", "invalid role")
	}
}

func ValidateStatus(status string) error {
	switch strings.TrimSpace(status) {
	case StatusActive, StatusInactive:
		return nil
	default:
		return sharederrors.Invalid("INVALID_STATUS", "invalid user status")
	}
}

func parsePositiveInt(value string, fallback int) int {
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}
