package app

import (
	"akadia/domain"
	"akadia/internal/platform/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	m *middleware.MiddlewareManager,
	auth domain.AuthHandler,
	paymentProduct domain.PaymentProductHandler,
	paymentPolicy domain.PaymentPolicyHandler,
) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "Berjalan Perfecto"})
	})

	authApi := r.Group("/auth")
	authApi.POST("/login", auth.Login)
	authApi.GET("/profile", m.JWT, auth.Profile)

	paymentPolicyApi := r.Group("/payment-policy", m.JWT)
	paymentPolicyApi.GET("/:id", paymentPolicy.FindByID)
	paymentPolicyApi.GET("", paymentPolicy.FindAll)
	paymentPolicyApi.PUT("/:id", paymentPolicy.Update)

	return r
}
