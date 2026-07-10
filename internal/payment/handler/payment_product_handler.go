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

// FindAll godoc
// @Summary List payment products
// @Description Lists payment products that define what the school bills, including the linked payment policy and revenue account mapping. Allowed roles: SUPER_ADMIN, SCHOOL_ADMIN, TREASURER.
// @Tags Payment Product
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param size query int false "Page size (default 15, max 100)"
// @Param keyword query string false "Search keyword for code or name"
// @Param payment_policy_id query string false "Filter by payment policy ID"
// @Param status query string false "Filter by product status: ACTIVE or INACTIVE"
// @Param sort_by query string false "Sort field: created_at, updated_at, code, name, price, status, payment_policy_id, revenue_account_code, revenue_account_name"
// @Param order query string false "Sort direction: asc or desc"
// @Success 200 {object} domain.SwaggerPaymentProductPageResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 403 {object} domain.ErrorResponse
// @Failure 422 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Security BearerAuth
// @Router /payment-product [get]
// @Router /payment-products [get]
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

// FindByID godoc
// @Summary Get payment product by ID
// @Description Returns one payment product that defines what the school bills, including the linked payment policy and revenue account mapping. Allowed roles: SUPER_ADMIN, SCHOOL_ADMIN, TREASURER.
// @Tags Payment Product
// @Produce json
// @Param id path string true "Payment product ID"
// @Success 200 {object} domain.SwaggerPaymentProductResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 403 {object} domain.ErrorResponse
// @Failure 404 {object} domain.ErrorResponse
// @Failure 422 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Security BearerAuth
// @Router /payment-product/{id} [get]
// @Router /payment-products/{id} [get]
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

// Create godoc
// @Summary Create payment product
// @Description Creates a payment product that defines what the school bills, including the linked payment policy and revenue account mapping. Allowed roles: SUPER_ADMIN, SCHOOL_ADMIN, TREASURER.
// @Tags Payment Product
// @Accept json
// @Produce json
// @Param request body domain.PaymentProductCreate true "Payment product payload"
// @Success 200 {object} domain.SwaggerPaymentProductResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 403 {object} domain.ErrorResponse
// @Failure 422 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Security BearerAuth
// @Router /payment-product [post]
// @Router /payment-products [post]
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

// Update godoc
// @Summary Update payment product
// @Description Updates a payment product that defines what the school bills while preserving the linked payment policy and revenue account mapping rules. Allowed roles: SUPER_ADMIN, SCHOOL_ADMIN, TREASURER.
// @Tags Payment Product
// @Accept json
// @Produce json
// @Param id path string true "Payment product ID"
// @Param request body domain.PaymentProductUpdate true "Payment product update payload"
// @Success 200 {object} domain.SwaggerPaymentProductResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 403 {object} domain.ErrorResponse
// @Failure 404 {object} domain.ErrorResponse
// @Failure 422 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Security BearerAuth
// @Router /payment-product/{id} [put]
// @Router /payment-products/{id} [put]
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
