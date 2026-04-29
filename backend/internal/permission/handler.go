package permission

import (
	"net/http"

	sharederrors "api-on/internal/shared/errors"
	"api-on/internal/shared/middleware"
	"api-on/internal/shared/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	usecase *Usecase
}

func NewHandler(usecase *Usecase) *Handler {
	return &Handler{usecase: usecase}
}

// Me retorna o mapa efetivo de permissões da conta autenticada.
//
// GET /v1/permissions/me
// Frontend: usa esta rota para decidir o que o usuário pode enxergar em pacientes,
// serviços, agenda, financeiro, IA e planos.
func (h *Handler) Me(c *gin.Context) {
	actor, ok := middleware.GetIdentity(c)
	if !ok {
		response.Fail(c, nil)
		return
	}

	result, err := h.usecase.GetMe(c.Request.Context(), actor)
	if err != nil {
		response.Fail(c, err)
		return
	}

	response.Success(c, http.StatusOK, result, nil)
}

// GetByUser retorna as permissões de uma conta específica do tenant.
//
// GET /v1/users/:id/permissions
// Frontend: owner/admin usam esta rota para auditoria e edição de acesso.
func (h *Handler) GetByUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, sharederrors.Invalid("INVALID_USER_ID", "invalid user id"))
		return
	}

	actor, ok := middleware.GetIdentity(c)
	if !ok {
		response.Fail(c, nil)
		return
	}

	result, err := h.usecase.GetByUserID(c.Request.Context(), actor, userID)
	if err != nil {
		response.Fail(c, err)
		return
	}

	response.Success(c, http.StatusOK, result, nil)
}

// UpdateByUser altera o escopo de dados que a conta pode enxergar.
//
// PATCH /v1/users/:id/permissions
// Regras:
// - somente owner/admin alteram permissões
// - escopo `all` libera visão completa do tenant
// - escopo `own` exige filtros por responsável/autor nas futuras consultas
func (h *Handler) UpdateByUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, sharederrors.Invalid("INVALID_USER_ID", "invalid user id"))
		return
	}

	var input UpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Fail(c, sharederrors.Invalid("INVALID_BODY", "invalid request body"))
		return
	}

	actor, ok := middleware.GetIdentity(c)
	if !ok {
		response.Fail(c, nil)
		return
	}

	result, err := h.usecase.UpdateByUserID(c.Request.Context(), actor, userID, input)
	if err != nil {
		response.Fail(c, err)
		return
	}

	response.Success(c, http.StatusOK, result, nil)
}
