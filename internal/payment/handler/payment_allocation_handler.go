package handler

import (
	"akadia/domain"
	"akadia/internal/platform/security"
	"akadia/internal/shared"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type paymentAllocationHandler struct {
	paymentAllocationService domain.PaymentAllocationService
}

func NewPaymentAllocationHandler(
	paymentAllocationService domain.PaymentAllocationService,
) domain.PaymentAllocationHandler {
	return &paymentAllocationHandler{
		paymentAllocationService: paymentAllocationService,
	}
}

// Allocate godoc
// @Summary Allocate payment order
// @Description Distributes one payment order to one or more student obligations using the existing payment policy and no over-allocation validation rules. Allowed roles: SUPER_ADMIN, SCHOOL_ADMIN, TREASURER, PARENT, STUDENT.
// @Tags Payment Allocation
// @Accept json
// @Produce json
// @Param id path string true "Payment order ID"
// @Param request body domain.PaymentAllocationAllocate true "Payment allocation payload"
// @Success 200 {object} domain.SwaggerPaymentAllocationResultResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 403 {object} domain.ErrorResponse
// @Failure 404 {object} domain.ErrorResponse
// @Failure 422 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Security BearerAuth
// @Router /payment-orders/{id}/allocate [post]
func (h *paymentAllocationHandler) Allocate(c *gin.Context) {
	authContextValue, _ := c.Get(domain.ContextKeyAuth)
	authContext := authContextValue.(*security.AuthContext)

	paymentOrderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": shared.ErrInvalidIDParam.Error()})
		return
	}

	var req domain.PaymentAllocationAllocate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	res, err := h.paymentAllocationService.Allocate(c.Request.Context(), authContext, paymentOrderID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

// FindByPaymentOrderID godoc
// @Summary Get payment allocations by payment order
// @Description Returns allocation details for one payment order, including allocated totals and remaining amount under current finance portal access rules. Allowed roles: SUPER_ADMIN, SCHOOL_ADMIN, TREASURER, PARENT, STUDENT.
// @Tags Payment Allocation
// @Produce json
// @Param id path string true "Payment order ID"
// @Success 200 {object} domain.SwaggerPaymentAllocationResultResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 403 {object} domain.ErrorResponse
// @Failure 404 {object} domain.ErrorResponse
// @Failure 422 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Security BearerAuth
// @Router /payment-orders/{id}/allocations [get]
func (h *paymentAllocationHandler) FindByPaymentOrderID(c *gin.Context) {
	authContextValue, _ := c.Get(domain.ContextKeyAuth)
	authContext := authContextValue.(*security.AuthContext)

	paymentOrderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": shared.ErrInvalidIDParam.Error()})
		return
	}

	res, err := h.paymentAllocationService.FindByPaymentOrderID(c.Request.Context(), authContext, paymentOrderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}
