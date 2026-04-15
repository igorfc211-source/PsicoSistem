package validation

import (
	"errors"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
var specialCharRegex = regexp.MustCompile(`[!@#\$%\^&\*\(\)_\+\-=\[\]\{\};':"\\|,.<>\/?]`)

func ValidateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("name is required")
	}
	return nil
}

func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return errors.New("email is required")
	}
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 5 {
		return errors.New("password must have at least 5 characters")
	}
	if !specialCharRegex.MatchString(password) {
		return errors.New("password must contain at least 1 special character")
	}
	return nil
}

func ValidateTipoPlano(tipoPlano string) error {
	if strings.TrimSpace(tipoPlano) == "" {
		return errors.New("tipo_plano is required")
	}
	return nil
}