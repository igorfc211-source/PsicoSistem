package jwt

import (
	"errors"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type CustomClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwtlib.RegisteredClaims
}

type JWTService struct {
	secretKey []byte
	issuer    string
}

func NewJWTService(secret string, issuer string) *JWTService {
	return &JWTService{
		secretKey: []byte(secret),
		issuer:    issuer,
	}
}

func (j *JWTService) GenerateToken(userID uuid.UUID, email string) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwtlib.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   userID.String(),
			ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwtlib.NewNumericDate(time.Now()),
		},
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

func (j *JWTService) ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwtlib.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwtlib.Token) (interface{}, error) {
		_, ok := token.Method.(*jwtlib.SigningMethodHMAC)
		if !ok {
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