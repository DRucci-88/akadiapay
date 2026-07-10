package middleware

import (
	"akadia/domain"
	"akadia/model"

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

func (m *MiddlewareManager) Roles(
	allowedRoles ...model.RoleCode,
) gin.HandlerFunc {
	return NewRoleMiddleware(allowedRoles...)
}
