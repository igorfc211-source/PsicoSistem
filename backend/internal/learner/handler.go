package learner

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

// List lista os aprendentes da clínica autenticada.
//
// GET /v1/learners
// Frontend: usa nas telas de aprendentes, agenda e financeiro.
func (h *Handler) List(c *gin.Context) {
	actor, ok := middleware.GetIdentity(c)
	if !ok {
		response.Fail(c, nil)
		return
	}

	input := ParseListInput(
		c.Query("page"),
		c.Query("per_page"),
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

// Get retorna um aprendente do tenant atual.
//
// GET /v1/learners/:id
func (h *Handler) Get(c *gin.Context) {
	learnerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, sharederrors.Invalid("INVALID_LEARNER_ID", "invalid learner id"))
		return
	}

	actor, ok := middleware.GetIdentity(c)
	if !ok {
		response.Fail(c, nil)
		return
	}

	result, err := h.usecase.Get(c.Request.Context(), actor, learnerID)
	if err != nil {
		response.Fail(c, err)
		return
	}

	response.Success(c, http.StatusOK, result, nil)
}

// Create cria um aprendente com o valor financeiro padrão da sessão.
//
// POST /v1/learners
// Campo financeiro: session_price_cents pertence ao aprendente, não à assinatura.
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

// Update altera dados cadastrais e o valor financeiro padrão do aprendente.
//
// PATCH /v1/learners/:id
func (h *Handler) Update(c *gin.Context) {
	learnerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, sharederrors.Invalid("INVALID_LEARNER_ID", "invalid learner id"))
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

	result, err := h.usecase.Update(c.Request.Context(), actor, learnerID, input)
	if err != nil {
		response.Fail(c, err)
		return
	}

	response.Success(c, http.StatusOK, result, nil)
}

// Delete desativa um aprendente da clínica.
//
// DELETE /v1/learners/:id
func (h *Handler) Delete(c *gin.Context) {
	learnerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, sharederrors.Invalid("INVALID_LEARNER_ID", "invalid learner id"))
		return
	}

	actor, ok := middleware.GetIdentity(c)
	if !ok {
		response.Fail(c, nil)
		return
	}

	if err := h.usecase.Deactivate(c.Request.Context(), actor, learnerID); err != nil {
		response.Fail(c, err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "learner deactivated successfully"}, nil)
}
