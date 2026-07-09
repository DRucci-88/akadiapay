package middleware

import (
	"akadia/domain"
	"akadia/internal/platform/security"
	"akadia/internal/shared"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func NewJWTMiddleware(
	appConfig domain.AppConfigProvider,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		log.Println("authHeader" + authHeader)
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": shared.ErrAuthUnauthorized.Error()})
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		// log.Println(tokenParts)
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": shared.ErrAuthUnauthorized.Error()})
			return
		}

		tokenString := tokenParts[1]

		// isTokenBlackListed, err := blackListedTokenRepo.IsTokenBlackListed(context.Background(), tokenString)
		// if err != nil {
		// 	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": model.ErrTokenValidatedFailed.Error()})
		// 	return
		// }
		// if isTokenBlackListed {
		// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": model.ErrTokenIsBlackListed.Error()})
		// 	return
		// }

		claims := &security.JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
			return appConfig.JWT_SECRET_BYTE(), nil
		})
		log.Println(err)
		log.Println(token.Valid)

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": shared.ErrAuthTokenExpired.Error()})
			return
		}

		authContext := &security.AuthContext{
			Email:          claims.Email,
			UserID:         claims.UserID,
			TenantID:       claims.TenantID,
			StudentID:      claims.StudentID,
			RoleCode:       claims.RoleCode,
			Token:          tokenString,
			TokenExpiredAt: claims.ExpiresAt.Time,
		}
		log.Println(authContext)

		c.Set(domain.ContextKeyAuth, authContext)
		c.Next()
	}
}
