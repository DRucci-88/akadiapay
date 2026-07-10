package app

import (
	_ "akadia/docs"
	"akadia/domain"
	"akadia/internal/platform/middleware"
	"akadia/model"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(
	m *middleware.MiddlewareManager,
	auth domain.AuthHandler,
	paymentProduct domain.PaymentProductHandler,
	paymentPolicy domain.PaymentPolicyHandler,
	studentObligation domain.StudentObligationHandler,
	paymentOrder domain.PaymentOrderHandler,
	paymentAllocation domain.PaymentAllocationHandler,
	ledgerEntry domain.LedgerEntryHandler,
) *gin.Engine {
	r := gin.Default()

	r.GET("/health", health)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	financeAdminRoles := m.Roles(
		model.RoleCodeSuperAdmin,
		model.RoleCodeSchoolAdmin,
		model.RoleCodeTreasurer,
	)
	financePortalRoles := m.Roles(
		model.RoleCodeSuperAdmin,
		model.RoleCodeSchoolAdmin,
		model.RoleCodeTreasurer,
		model.RoleCodeParent,
		model.RoleCodeStudent,
	)

	authApi := r.Group("/auth")
	authApi.POST("/login", auth.Login)
	authApi.GET("/profile", m.JWT, auth.Profile)

	registerPaymentPolicyRoutes := func(group *gin.RouterGroup) {
		group.GET("/:id", paymentPolicy.FindByID)
		group.GET("", paymentPolicy.FindAll)
		group.PUT("/:id", paymentPolicy.Update)
		group.POST("", paymentPolicy.Create)
	}
	registerPaymentProductRoutes := func(group *gin.RouterGroup) {
		group.GET("/:id", paymentProduct.FindByID)
		group.GET("", paymentProduct.FindAll)
		group.PUT("/:id", paymentProduct.Update)
		group.POST("", paymentProduct.Create)
	}

	registerPaymentPolicyRoutes(r.Group("/payment-policy", m.JWT, financeAdminRoles))
	registerPaymentPolicyRoutes(r.Group("/payment-policies", m.JWT, financeAdminRoles))

	registerPaymentProductRoutes(r.Group("/payment-product", m.JWT, financeAdminRoles))
	registerPaymentProductRoutes(r.Group("/payment-products", m.JWT, financeAdminRoles))

	studentObligationApi := r.Group("/student-obligations", m.JWT, financeAdminRoles)
	studentObligationApi.GET("", studentObligation.FindAll)
	studentObligationApi.GET("/:id", studentObligation.FindByID)
	studentObligationApi.POST("/bulk", studentObligation.CreateBulk)
	studentObligationApi.POST("", studentObligation.Create)
	studentObligationApi.PUT("/:id", studentObligation.Update)
	studentObligationApi.DELETE("/:id", studentObligation.Delete)

	studentApi := r.Group("/students", m.JWT, financePortalRoles)
	studentApi.GET("/:studentId/outstanding", studentObligation.OutstandingByStudentID)

	paymentOrderApi := r.Group("/payment-orders", m.JWT, financePortalRoles)
	paymentOrderApi.GET("", paymentOrder.FindAll)
	paymentOrderApi.GET("/:id", paymentOrder.FindByID)
	paymentOrderApi.POST("", paymentOrder.Create)
	paymentOrderApi.POST("/:id/cancel", paymentOrder.Cancel)
	paymentOrderApi.POST("/:id/allocate", paymentAllocation.Allocate)
	paymentOrderApi.GET("/:id/allocations", paymentAllocation.FindByPaymentOrderID)
	paymentOrderApi.GET("/:id/ledger", financeAdminRoles, ledgerEntry.FindByPaymentOrderID)
	paymentOrderApi.POST("/:id/post-ledger", financeAdminRoles, ledgerEntry.PostPayment)

	ledgerEntryApi := r.Group("/ledger-entries", m.JWT, financeAdminRoles)
	ledgerEntryApi.GET("", ledgerEntry.FindAll)

	return r
}
