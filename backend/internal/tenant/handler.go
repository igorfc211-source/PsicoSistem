package tenant

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

// Me retorna os dados da clínica autenticada.
//
// GET /v1/tenant/me
// Frontend: usa esta rota na tela de configurações da clínica.
func (h *Handler) Me(c *gin.Context) {
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
