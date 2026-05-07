package guardian

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

// List lista responsáveis da clínica autenticada.
//
// GET /v1/guardians
func (h *Handler) List(c *gin.Context) {
	actor, ok := middleware.GetIdentity(c)
	if !ok {
		response.Fail(c, nil)
		return
	}

	input := ParseListInput(
		c.Query("page"),
		c.Query("per_page"),
		c.Query("search"),
	)

	result, meta, err := h.usecase.List(c.Request.Context(), actor, input)
	if err != nil {
		response.Fail(c, err)
		return
	}

	response.Success(c, http.StatusOK, result, meta)
}

// Get retorna um responsável do tenant atual.
//
// GET /v1/guardians/:id
func (h *Handler) Get(c *gin.Context) {
	guardianID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, sharederrors.Invalid("INVALID_GUARDIAN_ID", "invalid guardian id"))
		return
	}

	actor, ok := middleware.GetIdentity(c)
	if !ok {
		response.Fail(c, nil)
		return
	}

	result, err := h.usecase.Get(c.Request.Context(), actor, guardianID)
	if err != nil {
		response.Fail(c, err)
		return
	}

	response.Success(c, http.StatusOK, result, nil)
}

// Create cria o cadastro de um responsável.
//
// POST /v1/guardians
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

// Update altera dados cadastrais de um responsável.
//
// PATCH /v1/guardians/:id
func (h *Handler) Update(c *gin.Context) {
	guardianID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, sharederrors.Invalid("INVALID_GUARDIAN_ID", "invalid guardian id"))
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

	result, err := h.usecase.Update(c.Request.Context(), actor, guardianID, input)
	if err != nil {
		response.Fail(c, err)
		return
	}

	response.Success(c, http.StatusOK, result, nil)
}

// Delete remove um responsável sem aprendentes vinculados.
//
// DELETE /v1/guardians/:id
func (h *Handler) Delete(c *gin.Context) {
	guardianID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, sharederrors.Invalid("INVALID_GUARDIAN_ID", "invalid guardian id"))
		return
	}

	actor, ok := middleware.GetIdentity(c)
	if !ok {
		response.Fail(c, nil)
		return
	}

	if err := h.usecase.Delete(c.Request.Context(), actor, guardianID); err != nil {
		response.Fail(c, err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "guardian deleted successfully"}, nil)
}
