package validator

import (
	"fmt"
	"regexp"
	"strings"

	sharederrors "api-on/internal/shared/errors"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
var specialCharRegex = regexp.MustCompile(`[!@#\$%\^&\*\(\)_\+\-=\[\]\{\};':"\\|,.<>\/?]`)
var slugInvalidCharsRegex = regexp.MustCompile(`[^a-z0-9\-]+`)
var duplicateDashRegex = regexp.MustCompile(`\-+`)

func NormalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}

func ValidateRequired(fieldName string, value string) error {
	if strings.TrimSpace(value) == "" {
		return sharederrors.Invalid("VALIDATION_ERROR", fmt.Sprintf("%s is required", fieldName))
	}
	return nil
}

func ValidateName(name string) error {
	if err := ValidateRequired("name", name); err != nil {
		return err
	}
	return nil
}

func ValidateClinicName(name string) error {
	if err := ValidateRequired("clinic_name", name); err != nil {
		return err
	}
	return nil
}

func ValidateEmail(email string) error {
	email = NormalizeEmail(email)
	if email == "" {
		return sharederrors.Invalid("INVALID_EMAIL", "email is required")
	}
	if !emailRegex.MatchString(email) {
		return sharederrors.Invalid("INVALID_EMAIL", "invalid email format")
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(strings.TrimSpace(password)) < 5 {
		return sharederrors.Invalid("INVALID_PASSWORD", "password must have at least 5 characters")
	}
	if !specialCharRegex.MatchString(password) {
		return sharederrors.Invalid("INVALID_PASSWORD", "password must contain at least 1 special character")
	}
	return nil
}

func ValidatePhone(phone string) error {
	if err := ValidateRequired("phone", phone); err != nil {
		return err
	}
	return nil
}

func Slugify(value string) string {
	slug := strings.TrimSpace(strings.ToLower(value))
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = slugInvalidCharsRegex.ReplaceAllString(slug, "-")
	slug = duplicateDashRegex.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		return "tenant"
	}
	return slug
}
