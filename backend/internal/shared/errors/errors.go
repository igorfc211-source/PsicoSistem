package errors

import "net/http"

// AppError padroniza erros de negócio e transporte para a API.
type AppError struct {
	Code    string
	Message string
	Status  int
}

func (e *AppError) Error() string {
	return e.Message
}

func New(code string, status int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

func Internal(message string) *AppError {
	return New("INTERNAL_ERROR", http.StatusInternalServerError, message)
}

func Invalid(code string, message string) *AppError {
	return New(code, http.StatusBadRequest, message)
}

func Unauthorized(message string) *AppError {
	return New("UNAUTHORIZED", http.StatusUnauthorized, message)
}

func Forbidden(message string) *AppError {
	return New("FORBIDDEN", http.StatusForbidden, message)
}

func NotFound(code string, message string) *AppError {
	return New(code, http.StatusNotFound, message)
}

func Conflict(code string, message string) *AppError {
	return New(code, http.StatusConflict, message)
}

func AsAppError(err error) *AppError {
	if err == nil {
		return nil
	}

	if appErr, ok := err.(*AppError); ok {
		return appErr
	}

	return Internal("internal server error")
}
