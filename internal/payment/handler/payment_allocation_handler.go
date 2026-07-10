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
