package subscription

import (
	"net/http"

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

// Current retorna o plano ativo da clínica atual.
//
// GET /v1/tenant/subscription
// Frontend: usa esta rota na tela de assinatura e billing.
func (h *Handler) Current(c *gin.Context) {
	actor, ok := middleware.GetIdentity(c)
	if !ok {
		response.Fail(c, nil)
		return
	}

	result, err := h.usecase.GetCurrent(c.Request.Context(), actor)
	if err != nil {
		response.Fail(c, err)
		return
	}

	response.Success(c, http.StatusOK, result, nil)
}
