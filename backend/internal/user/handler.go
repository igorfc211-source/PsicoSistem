package user

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

// Me retorna o usuário interno autenticado.
//
// GET /v1/users/me
// Frontend: usa para carregar o perfil atual e permissões da sessão.
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

// List lista os profissionais da clínica autenticada.
//
// GET /v1/users
// Frontend: usa na tela de equipe clínica com filtros e paginação.
func (h *Handler) List(c *gin.Context) {
	actor, ok := middleware.GetIdentity(c)
	if !ok {
		response.Fail(c, nil)
		return
	}

	input := ParseListInput(
		c.Query("page"),
		c.Query("per_page"),
		c.Query("role"),
		c.Query("status"),
		c.Query("search"),
	)

	result, meta, err := h.usecase.List(c.Request.Context(), actor, input)
	if err != nil {
		response.Fail(c, err)
		return
	}

	response.Success(c, http.StatusOK, result, meta)
}

// Create cria um novo usuário interno dentro do tenant autenticado.
//
// POST /v1/users
// Regras:
// - tenant_id vem do JWT, nunca do body
// - apenas owner e admin podem criar
// - respeita o limite do plano atual
func (h *Handler) Create(c *gin.Context) {
	var input CreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Fail(c, sharederrors.Invalid("INVALID_BODY", "invalid request body"))
		return
	}

	actor, ok := middleware.GetIdentity(c)
	if !ok {
		response.Fail(c, nil)
		return
	}

	result, err := h.usecase.Create(c.Request.Context(), actor, input)
	if err != nil {
		response.Fail(c, err)
		return
	}

	response.Success(c, http.StatusCreated, result, nil)
}

// Update altera nome, email, papel e status de um usuário do tenant atual.
//
// PATCH /v1/users/:id
// Frontend: usa na edição de membros da equipe.
func (h *Handler) Update(c *gin.Context) {
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

	result, err := h.usecase.Update(c.Request.Context(), actor, userID, input)
	if err != nil {
		response.Fail(c, err)
		return
	}

	response.Success(c, http.StatusOK, result, nil)
}

// Delete desativa um usuário da clínica.
//
// DELETE /v1/users/:id
// Em produção a base já segue soft delete por status `inactive`.
func (h *Handler) Delete(c *gin.Context) {
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

	if err := h.usecase.Deactivate(c.Request.Context(), actor, userID); err != nil {
		response.Fail(c, err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "user deactivated successfully"}, nil)
}
