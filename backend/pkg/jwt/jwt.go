package jwt

import (
	"errors"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

// CustomClaims carrega o contexto necessário para segurança multi-tenant.
type CustomClaims struct {
	UserID   string `json:"user_id"`
	TenantID string `json:"tenant_id"`
	Role     string `json:"role"`
	Email    string `json:"email"`
	Type     string `json:"type"`
	jwtlib.RegisteredClaims
}

type TokenInput struct {
	UserID   string
	TenantID string
	Role     string
	Email    string
	Type     string
}

type JWTService struct {
	secretKey []byte
	issuer    string
	ttl       time.Duration
}

func NewJWTService(secret string, issuer string) *JWTService {
	return &JWTService{
		secretKey: []byte(secret),
		issuer:    issuer,
		ttl:       24 * time.Hour,
	}
}

func (j *JWTService) GenerateToken(input TokenInput) (string, error) {
	now := time.Now()
	claims := CustomClaims{
		UserID:   input.UserID,
		TenantID: input.TenantID,
		Role:     input.Role,
		Email:    input.Email,
		Type:     input.Type,
		RegisteredClaims: jwtlib.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   input.UserID,
			ExpiresAt: jwtlib.NewNumericDate(now.Add(j.ttl)),
			IssuedAt:  jwtlib.NewNumericDate(now),
			NotBefore: jwtlib.NewNumericDate(now),
		},
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

func (j *JWTService) ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwtlib.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwtlib.Token) (any, error) {
		if _, ok := token.Method.(*jwtlib.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return j.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
