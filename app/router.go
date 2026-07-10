package app

import (
	"akadia/domain"
	"akadia/internal/platform/middleware"
	"akadia/model/generated"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	m *middleware.MiddlewareManager,
	auth domain.AuthHandler,
	paymentProduct domain.PaymentProductHandler,
	paymentPolicy domain.PaymentPolicyHandler,
	studentObligation domain.StudentObligationHandler,
	paymentOrder domain.PaymentOrderHandler,
	paymentAllocation domain.PaymentAllocationHandler,
) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "Berjalan Perfecto", "aaa": generated.PaymentPolicy.Code.Column().Name})
	})

	authApi := r.Group("/auth")
	authApi.POST("/login", auth.Login)
	authApi.GET("/profile", m.JWT, auth.Profile)

	paymentPolicyApi := r.Group("/payment-policy", m.JWT)
	paymentPolicyApi.GET("/:id", paymentPolicy.FindByID)
	paymentPolicyApi.GET("", paymentPolicy.FindAll)
	paymentPolicyApi.PUT("/:id", paymentPolicy.Update)
	paymentPolicyApi.POST("", paymentPolicy.Create)

	paymentProductApi := r.Group("/payment-product", m.JWT)
	paymentProductApi.GET("/:id", paymentProduct.FindByID)
	paymentProductApi.GET("", paymentProduct.FindAll)
	paymentProductApi.PUT("/:id", paymentProduct.Update)
	paymentProductApi.POST("", paymentProduct.Create)

	studentObligationApi := r.Group("/student-obligations", m.JWT)
	studentObligationApi.GET("", studentObligation.FindAll)
	studentObligationApi.GET("/:id", studentObligation.FindByID)
	studentObligationApi.POST("/bulk", studentObligation.CreateBulk)
	studentObligationApi.POST("", studentObligation.Create)
	studentObligationApi.PUT("/:id", studentObligation.Update)
	studentObligationApi.DELETE("/:id", studentObligation.Delete)

	studentApi := r.Group("/students", m.JWT)
	studentApi.GET("/:studentId/outstanding", studentObligation.OutstandingByStudentID)

	paymentOrderApi := r.Group("/payment-orders", m.JWT)
	paymentOrderApi.GET("", paymentOrder.FindAll)
	paymentOrderApi.GET("/:id", paymentOrder.FindByID)
	paymentOrderApi.POST("", paymentOrder.Create)
	paymentOrderApi.POST("/:id/cancel", paymentOrder.Cancel)
	paymentOrderApi.POST("/:id/allocate", paymentAllocation.Allocate)
	paymentOrderApi.GET("/:id/allocations", paymentAllocation.FindByPaymentOrderID)

	return r
}
