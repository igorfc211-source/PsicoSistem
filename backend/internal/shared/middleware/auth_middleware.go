package middleware

import (
	"strings"

	sharederrors "api-on/internal/shared/errors"
	"api-on/internal/shared/response"
	"api-on/internal/shared/security"
	jwtpkg "api-on/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const identityContextKey = "identity"

// AuthRequired valida o JWT e injeta a identidade autenticada no contexto.
func AuthRequired(jwtSvc *jwtpkg.JWTService, allowedTypes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if strings.TrimSpace(authHeader) == "" {
			response.Fail(c, sharederrors.Unauthorized("missing authorization header"))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" || strings.TrimSpace(parts[1]) == "" {
			response.Fail(c, sharederrors.Unauthorized("invalid authorization header"))
			c.Abort()
			return
		}

		claims, err := jwtSvc.ValidateToken(parts[1])
		if err != nil {
			response.Fail(c, sharederrors.Unauthorized("invalid or expired token"))
			c.Abort()
			return
		}

		if len(allowedTypes) > 0 && !contains(allowedTypes, claims.Type) {
			response.Fail(c, sharederrors.Forbidden("token type is not allowed for this route"))
			c.Abort()
			return
		}

		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			response.Fail(c, sharederrors.Unauthorized("invalid user claim"))
			c.Abort()
			return
		}

		tenantID, err := uuid.Parse(claims.TenantID)
		if err != nil {
			response.Fail(c, sharederrors.Unauthorized("invalid tenant claim"))
			c.Abort()
			return
		}

		c.Set(identityContextKey, security.Identity{
			UserID:   userID,
			TenantID: tenantID,
			Role:     claims.Role,
			Email:    claims.Email,
			Type:     claims.Type,
		})
		c.Next()
	}
}

// AuthRequiredWithResolver valida o token e recarrega a identidade atual do banco.
func AuthRequiredWithResolver(
	jwtSvc *jwtpkg.JWTService,
	resolver security.IdentityResolver,
	allowedTypes ...string,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if strings.TrimSpace(authHeader) == "" {
			response.Fail(c, sharederrors.Unauthorized("missing authorization header"))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" || strings.TrimSpace(parts[1]) == "" {
			response.Fail(c, sharederrors.Unauthorized("invalid authorization header"))
			c.Abort()
			return
		}

		claims, err := jwtSvc.ValidateToken(parts[1])
		if err != nil {
			response.Fail(c, sharederrors.Unauthorized("invalid or expired token"))
			c.Abort()
			return
		}

		if len(allowedTypes) > 0 && !contains(allowedTypes, claims.Type) {
			response.Fail(c, sharederrors.Forbidden("token type is not allowed for this route"))
			c.Abort()
			return
		}

		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			response.Fail(c, sharederrors.Unauthorized("invalid user claim"))
			c.Abort()
			return
		}

		tenantID, err := uuid.Parse(claims.TenantID)
		if err != nil {
			response.Fail(c, sharederrors.Unauthorized("invalid tenant claim"))
			c.Abort()
			return
		}

		identity, err := resolver.ResolveInternalIdentity(c.Request.Context(), tenantID, userID)
		if err != nil {
			response.Fail(c, err)
			c.Abort()
			return
		}

		c.Set(identityContextKey, identity)
		c.Next()
	}
}

// RequireRoles restringe o acesso a papéis explícitos do painel interno.
func RequireRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		identity, ok := GetIdentity(c)
		if !ok {
			response.Fail(c, sharederrors.Unauthorized("authenticated identity not found"))
			c.Abort()
			return
		}

		if !identity.HasRole(roles...) {
			response.Fail(c, sharederrors.Forbidden("you do not have permission to access this resource"))
			c.Abort()
			return
		}

		c.Next()
	}
}

func GetIdentity(c *gin.Context) (security.Identity, bool) {
	value, exists := c.Get(identityContextKey)
	if !exists {
		return security.Identity{}, false
	}

	identity, ok := value.(security.Identity)
	return identity, ok
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
