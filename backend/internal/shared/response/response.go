package response

import (
	"api-on/internal/shared/errors"

	"github.com/gin-gonic/gin"
)

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type APIResponse struct {
	Data  any        `json:"data"`
	Meta  any        `json:"meta"`
	Error *ErrorBody `json:"error"`
}

// Success envia uma resposta previsível para o frontend.
func Success(c *gin.Context, status int, data any, meta any) {
	c.JSON(status, APIResponse{
		Data:  data,
		Meta:  meta,
		Error: nil,
	})
}

// Fail converte AppError em payload padrão da API.
func Fail(c *gin.Context, err error) {
	appErr := errors.AsAppError(err)
	if appErr == nil {
		appErr = errors.Unauthorized("unauthorized")
	}
	c.JSON(appErr.Status, APIResponse{
		Data: nil,
		Meta: nil,
		Error: &ErrorBody{
			Code:    appErr.Code,
			Message: appErr.Message,
		},
	})
}
