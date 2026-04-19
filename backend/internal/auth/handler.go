package auth

import (
	"net/http"

	sharederrors "api-on/internal/shared/errors"
	"api-on/internal/shared/middleware"
	"api-on/internal/shared/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	usecase *Usecase
}

func NewHandler(usecase *Usecase) *Handler {
	return &Handler{usecase: usecase}
}

// Register cria a primeira conta da clínica, seu tenant e a assinatura inicial.
//
// POST /v1/auth/register
// Frontend: usa na etapa inicial de onboarding do SaaS.
func (h *Handler) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Fail(c, sharederrors.Invalid("INVALID_BODY", "invalid request body"))
		return
	}

	result, err := h.usecase.Register(c.Request.Context(), input)
	if err != nil {
		response.Fail(c, err)
		return
	}

	response.Success(c, http.StatusCreated, result, nil)
}

// Login autentica usuários internos do tenant.
//
// POST /v1/auth/login
// Frontend: usa na tela padrão de login da clínica.
func (h *Handler) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Fail(c, sharederrors.Invalid("INVALID_BODY", "invalid request body"))
		return
	}

	result, err := h.usecase.Login(c.Request.Context(), input)
	if err != nil {
		response.Fail(c, err)
		return
	}

	response.Success(c, http.StatusOK, result, nil)
}

// Refresh gera um novo access token para a sessão atual.
//
// POST /v1/auth/refresh
// Frontend: usa quando o access token expira.
func (h *Handler) Refresh(c *gin.Context) {
	actor, ok := middleware.GetIdentity(c)
	if !ok {
		response.Fail(c, nil)
		return
	}

	result, err := h.usecase.Refresh(c.Request.Context(), actor)
	if err != nil {
		response.Fail(c, err)
		return
	}

	response.Success(c, http.StatusOK, result, nil)
}
