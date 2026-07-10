package handler

import (
	"akadia/domain"
	"akadia/internal/platform/security"
	"akadia/internal/shared"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type paymentPolicyHandler struct {
	paymentPolicyService domain.PaymentPolicyService
}

func NewPaymentPolicyHandler(paymentPolicyService domain.PaymentPolicyService) domain.PaymentPolicyHandler {
	return &paymentPolicyHandler{
		paymentPolicyService: paymentPolicyService,
	}
}

// FindAll godoc
// @Summary List payment policies
// @Description Lists payment policies that define how a bill can be paid, including partial payment rules, minimum amount, minimum percentage, overpayment support, and auto-close behavior. Allowed roles: SUPER_ADMIN, SCHOOL_ADMIN, TREASURER.
// @Tags Payment Policy
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param size query int false "Page size (default 15, max 100)"
// @Param keyword query string false "Search keyword for code or name"
// @Param allow_partial query bool false "Filter by partial payment support"
// @Param allow_over_payment query bool false "Filter by overpayment support"
// @Param auto_close_obligation query bool false "Filter by auto close obligation behavior"
// @Param sort_by query string false "Sort field: created_at, updated_at, code, name, minimum_amount, minimum_percentage, allow_partial, allow_over_payment, auto_close_obligation"
// @Param order query string false "Sort direction: asc or desc"
// @Success 200 {object} domain.SwaggerPaymentPolicyPageResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 403 {object} domain.ErrorResponse
// @Failure 422 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Security BearerAuth
// @Router /payment-policy [get]
// @Router /payment-policies [get]
func (h *paymentPolicyHandler) FindAll(c *gin.Context) {
	log.Println("PaymentPolicy.FindAll")

	// Auth Context
	authContextValue, _ := c.Get(domain.ContextKeyAuth)
	authContext := authContextValue.(*security.AuthContext)

	log.Println("authContext")

	// Pagination
	pageable, err := shared.GetPageable(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	log.Println("pageable")

	// Filter
	var req domain.PaymentPolicyFilter
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	log.Println("req")

	page, err := h.paymentPolicyService.FindPaginate(c.Request.Context(), pageable, &req, authContext)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("page")
	c.JSON(http.StatusOK, page)
}

// FindByID godoc
// @Summary Get payment policy by ID
// @Description Returns one payment policy that defines how a bill can be paid, including partial payment, minimum amount, minimum percentage, overpayment, and auto-close rules. Allowed roles: SUPER_ADMIN, SCHOOL_ADMIN, TREASURER.
// @Tags Payment Policy
// @Produce json
// @Param id path string true "Payment policy ID"
// @Success 200 {object} domain.SwaggerPaymentPolicyResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 403 {object} domain.ErrorResponse
// @Failure 404 {object} domain.ErrorResponse
// @Failure 422 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Security BearerAuth
// @Router /payment-policy/{id} [get]
// @Router /payment-policies/{id} [get]
func (h *paymentPolicyHandler) FindByID(c *gin.Context) {
	log.Println("PaymentPolicy.FindByID")
	authContextValue, _ := c.Get(domain.ContextKeyAuth)
	authContext := authContextValue.(*security.AuthContext)

	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": shared.ErrInvalidIDParam.Error()})
		return
	}

	paymentPolicy, err := h.paymentPolicyService.FirstByID(
		c.Request.Context(),
		authContext,
		id,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res := domain.NewPaymentPolicyResponse(paymentPolicy)

	c.JSON(http.StatusOK, gin.H{"data": res})
}

// Create godoc
// @Summary Create payment policy
// @Description Creates a payment policy that defines how a bill can be paid, including partial payment rules, minimum amount, minimum percentage, overpayment support, and auto-close behavior. Allowed roles: SUPER_ADMIN, SCHOOL_ADMIN, TREASURER.
// @Tags Payment Policy
// @Accept json
// @Produce json
// @Param request body domain.PaymentPolicyCreate true "Payment policy payload"
// @Success 200 {object} domain.SwaggerPaymentPolicyResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 403 {object} domain.ErrorResponse
// @Failure 422 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Security BearerAuth
// @Router /payment-policy [post]
// @Router /payment-policies [post]
func (h *paymentPolicyHandler) Create(c *gin.Context) {
	log.Println("PaymentPolicy.Create")

	// Auth Context
	authContextValue, _ := c.Get(domain.ContextKeyAuth)
	authContext := authContextValue.(*security.AuthContext)

	// Filter
	var req domain.PaymentPolicyCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	log.Println("req")

	paymentPolicy, err := h.paymentPolicyService.Create(c.Request.Context(), authContext, &req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res := domain.NewPaymentPolicyResponse(paymentPolicy)
	c.JSON(http.StatusOK, gin.H{"data": res})
}

// Update godoc
// @Summary Update payment policy
// @Description Updates a payment policy using pointer-based PATCH-style payload semantics while preserving current route behavior. Allowed roles: SUPER_ADMIN, SCHOOL_ADMIN, TREASURER.
// @Tags Payment Policy
// @Accept json
// @Produce json
// @Param id path string true "Payment policy ID"
// @Param request body domain.PaymentPolicyUpdate true "Payment policy update payload"
// @Success 200 {object} domain.SwaggerPaymentPolicyResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 403 {object} domain.ErrorResponse
// @Failure 404 {object} domain.ErrorResponse
// @Failure 422 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Security BearerAuth
// @Router /payment-policy/{id} [put]
// @Router /payment-policies/{id} [put]
func (h *paymentPolicyHandler) Update(c *gin.Context) {
	log.Println("PaymentPolicy.Update")

	// Auth Context
	authContextValue, _ := c.Get(domain.ContextKeyAuth)
	authContext := authContextValue.(*security.AuthContext)

	// ID Param
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": shared.ErrInvalidIDParam})
		return
	}

	// Filter
	var req domain.PaymentPolicyUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	log.Println("req")

	paymentPolicy, err := h.paymentPolicyService.Update(c.Request.Context(), authContext, id, &req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res := domain.NewPaymentPolicyResponse(paymentPolicy)
	c.JSON(http.StatusOK, gin.H{"data": res})
}
