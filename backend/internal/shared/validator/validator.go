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

func NormalizeDigits(value string) string {
	var normalized strings.Builder
	normalized.Grow(len(value))

	for _, char := range value {
		if char >= '0' && char <= '9' {
			normalized.WriteRune(char)
		}
	}

	return normalized.String()
}

func NormalizePhone(phone string) string {
	return NormalizeDigits(phone)
}

func NormalizeCPFOrCNPJ(value string) string {
	return NormalizeDigits(value)
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
	normalizedPhone := NormalizePhone(phone)
	if err := ValidateRequired("phone", normalizedPhone); err != nil {
		return err
	}
	if len(normalizedPhone) < 10 || len(normalizedPhone) > 11 {
		return sharederrors.Invalid("INVALID_PHONE", "phone must have 10 or 11 digits")
	}
	return nil
}

func ValidateCPFOrCNPJ(value string) error {
	normalizedValue := NormalizeCPFOrCNPJ(value)
	if err := ValidateRequired("cpf_cnpj", normalizedValue); err != nil {
		return err
	}

	switch len(normalizedValue) {
	case 11:
		if !isValidCPF(normalizedValue) {
			return sharederrors.Invalid("INVALID_TAX_DOCUMENT", "invalid CPF")
		}
	case 14:
		if !isValidCNPJ(normalizedValue) {
			return sharederrors.Invalid("INVALID_TAX_DOCUMENT", "invalid CNPJ")
		}
	default:
		return sharederrors.Invalid("INVALID_TAX_DOCUMENT", "cpf_cnpj must have 11 or 14 digits")
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

func isValidCPF(value string) bool {
	if hasRepeatedDigits(value) {
		return false
	}

	firstDigit := calculateCPFVerifier(value[:9], 10)
	secondDigit := calculateCPFVerifier(value[:10], 11)

	return value[9] == firstDigit && value[10] == secondDigit
}

func calculateCPFVerifier(base string, weight int) byte {
	sum := 0
	for _, char := range base {
		sum += int(char-'0') * weight
		weight--
	}

	remainder := (sum * 10) % 11
	if remainder == 10 {
		remainder = 0
	}

	return byte('0' + remainder)
}

func isValidCNPJ(value string) bool {
	if hasRepeatedDigits(value) {
		return false
	}

	firstDigit := calculateCNPJVerifier(value[:12], []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2})
	secondDigit := calculateCNPJVerifier(value[:12]+string(firstDigit), []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2})

	return value[12] == firstDigit && value[13] == secondDigit
}

func calculateCNPJVerifier(base string, weights []int) byte {
	sum := 0
	for index, char := range base {
		sum += int(char-'0') * weights[index]
	}

	remainder := sum % 11
	if remainder < 2 {
		return '0'
	}

	return byte('0' + (11 - remainder))
}

func hasRepeatedDigits(value string) bool {
	if len(value) == 0 {
		return true
	}

	first := value[0]
	for index := 1; index < len(value); index++ {
		if value[index] != first {
			return false
		}
	}

	return true
}
