package handler

import (
	"akadia/domain"
	"akadia/internal/shared"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type paymentProductHandler struct {
	paymentProductService domain.PaymentProductService
}

func NewPaymentProductHandler(paymentProductService domain.PaymentProductService) domain.PaymentProductHandler {
	return &paymentProductHandler{
		paymentProductService: paymentProductService,
	}
}

func (h *paymentProductHandler) FindAll(c *gin.Context) {

}
func (h *paymentProductHandler) FindByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": shared.ErrInvalidIDParam})
		return
	}

	paymentProduct, err := h.paymentProductService.FindByID(
		c.Request.Context(),
		id,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := domain.NewPaymentProductResponse(paymentProduct)

	c.JSON(http.StatusOK, gin.H{"data": res})
}
