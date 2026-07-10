package handler

import (
	"akadia/domain"
	"akadia/internal/platform/security"
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
	authContextValue, _ := c.Get(domain.ContextKeyAuth)
	authContext := authContextValue.(*security.AuthContext)

	pageable, err := shared.GetPageable(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var req domain.PaymentProductFilter
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	page, err := h.paymentProductService.FindPaginate(c.Request.Context(), pageable, &req, authContext)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, page)
}

func (h *paymentProductHandler) FindByID(c *gin.Context) {
	authContextValue, _ := c.Get(domain.ContextKeyAuth)
	authContext := authContextValue.(*security.AuthContext)

	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": shared.ErrInvalidIDParam.Error()})
		return
	}

	paymentProduct, err := h.paymentProductService.FindByID(
		c.Request.Context(),
		authContext,
		id,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := domain.NewPaymentProductResponse(paymentProduct)

	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h *paymentProductHandler) Create(c *gin.Context) {
	authContextValue, _ := c.Get(domain.ContextKeyAuth)
	authContext := authContextValue.(*security.AuthContext)

	var req domain.PaymentProductCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	paymentProduct, err := h.paymentProductService.Create(c.Request.Context(), authContext, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := domain.NewPaymentProductResponse(paymentProduct)
	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h *paymentProductHandler) Update(c *gin.Context) {
	authContextValue, _ := c.Get(domain.ContextKeyAuth)
	authContext := authContextValue.(*security.AuthContext)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": shared.ErrInvalidIDParam.Error()})
		return
	}

	var req domain.PaymentProductUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	paymentProduct, err := h.paymentProductService.Update(c.Request.Context(), authContext, id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := domain.NewPaymentProductResponse(paymentProduct)
	c.JSON(http.StatusOK, gin.H{"data": res})
}
