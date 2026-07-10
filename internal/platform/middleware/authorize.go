package middleware

import (
	"akadia/domain"
	"akadia/internal/platform/security"
	"akadia/internal/shared"
	"akadia/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRoleMiddleware(
	allowedRoles ...model.RoleCode,
) gin.HandlerFunc {
	allowed := make(map[model.RoleCode]struct{}, len(allowedRoles))
	for _, roleCode := range allowedRoles {
		allowed[roleCode] = struct{}{}
	}

	return func(c *gin.Context) {
		authContextValue, exists := c.Get(domain.ContextKeyAuth)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": shared.ErrAuthUnauthorized.Error()})
			return
		}

		authContext := authContextValue.(*security.AuthContext)
		if _, ok := allowed[authContext.RoleCode]; !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": shared.ErrAuthUnauthorized.Error()})
			return
		}

		c.Next()
	}
}
