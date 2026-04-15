package middleware

import (
	"net/http"
	"strings"

	jwtpkg "api-on/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AuthRequired(jwtSvc *jwtpkg.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := jwtSvc.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()
	}
}

func GetUserIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	value, exists := c.Get("userID")
	if !exists {
		return uuid.Nil, false
	}

	userID, ok := value.(uuid.UUID)
	return userID, ok
}