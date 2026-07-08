package app

import (
	"akadia/domain"
	"akadia/plarform/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	m *middleware.MiddlewareManager,
	auth domain.AuthHandler,
) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "Berjalan Perfecto"})
	})

	authApi := r.Group("/auth")
	authApi.POST("/login", auth.Login)
	authApi.GET("/profile", auth.Profile)

	return r
}
