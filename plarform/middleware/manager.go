package middleware

import (
	"akadia/domain"

	"github.com/gin-gonic/gin"
)

type MiddlewareManager struct {
	JWT gin.HandlerFunc
}

func NewMiddlewareManager(
	appConfig domain.AppConfigProvider,
) *MiddlewareManager {
	return &MiddlewareManager{
		JWT: NewJWTMiddleware(appConfig),
	}
}
