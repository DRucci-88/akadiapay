package handler

import (
	"akadia/domain"
	"akadia/internal/platform/security"
	"akadia/internal/shared"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type paymentOrderHandler struct {
	paymentOrderService domain.PaymentOrderService
}

func NewPaymentOrderHandler(
	paymentOrderService domain.PaymentOrderService,
) domain.PaymentOrderHandler {
	return &paymentOrderHandler{
		paymentOrderService: paymentOrderService,
	}
}

// FindAll godoc
// @Summary List payment orders
// @Description Lists payment transaction orders with pending, completed, and cancelled statuses under the current finance portal access rules. Allowed roles: SUPER_ADMIN, SCHOOL_ADMIN, TREASURER, PARENT, STUDENT.
// @Tags Payment Order
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param size query int false "Page size (default 15, max 100)"
// @Param keyword query string false "Search keyword for student, order number, reference number, or notes"
// @Param student_id query string false "Filter by student ID"
// @Param status query string false "Filter by order status: PENDING, COMPLETED, CANCELLED, EXPIRED"
// @Param payment_method query string false "Filter by payment method: CASH, BANK_TRANSFER, VIRTUAL_ACCOUNT, QRIS, CREDIT_CARD"
// @Param order_date_from query string false "Filter order date from (YYYY-MM-DD)"
// @Param order_date_to query string false "Filter order date to (YYYY-MM-DD)"
// @Param sort_by query string false "Sort field: created_at, updated_at, student_id, order_number, order_date, total_amount, status, payment_method, ledger_posted_at"
// @Param order query string false "Sort direction: asc or desc"
// @Success 200 {object} domain.SwaggerPaymentOrderPageResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 403 {object} domain.ErrorResponse
// @Failure 422 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Security BearerAuth
// @Router /payment-orders [get]
func (h *paymentOrderHandler) FindAll(c *gin.Context) {
	authContextValue, _ := c.Get(domain.ContextKeyAuth)
	authContext := authContextValue.(*security.AuthContext)

	pageable, err := shared.GetPageable(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var req domain.PaymentOrderFilter
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	page, err := h.paymentOrderService.FindPaginate(c.Request.Context(), pageable, &req, authContext)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, page)
}

// FindByID godoc
// @Summary Get payment order by ID
// @Description Returns one payment transaction order, including its pending, completed, or cancelled status under current finance portal access rules. Allowed roles: SUPER_ADMIN, SCHOOL_ADMIN, TREASURER, PARENT, STUDENT.
// @Tags Payment Order
// @Produce json
// @Param id path string true "Payment order ID"
// @Success 200 {object} domain.SwaggerPaymentOrderResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 403 {object} domain.ErrorResponse
// @Failure 404 {object} domain.ErrorResponse
// @Failure 422 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Security BearerAuth
// @Router /payment-orders/{id} [get]
func (h *paymentOrderHandler) FindByID(c *gin.Context) {
	authContextValue, _ := c.Get(domain.ContextKeyAuth)
	authContext := authContextValue.(*security.AuthContext)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": shared.ErrInvalidIDParam.Error()})
		return
	}

	paymentOrder, err := h.paymentOrderService.FindByID(c.Request.Context(), authContext, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := domain.NewPaymentOrderResponse(paymentOrder)
	c.JSON(http.StatusOK, gin.H{"data": res})
}

// Create godoc
// @Summary Create payment order
// @Description Creates a payment transaction order for the selected student obligation balance using the existing payment method, outstanding amount, and overpayment validation rules. Allowed roles: SUPER_ADMIN, SCHOOL_ADMIN, TREASURER, PARENT, STUDENT.
// @Tags Payment Order
// @Accept json
// @Produce json
// @Param request body domain.PaymentOrderCreate true "Payment order payload"
// @Success 200 {object} domain.SwaggerPaymentOrderResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 403 {object} domain.ErrorResponse
// @Failure 422 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Security BearerAuth
// @Router /payment-orders [post]
func (h *paymentOrderHandler) Create(c *gin.Context) {
	authContextValue, _ := c.Get(domain.ContextKeyAuth)
	authContext := authContextValue.(*security.AuthContext)

	var req domain.PaymentOrderCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	paymentOrder, err := h.paymentOrderService.Create(c.Request.Context(), authContext, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := domain.NewPaymentOrderResponse(paymentOrder)
	c.JSON(http.StatusOK, gin.H{"data": res})
}

// Cancel godoc
// @Summary Cancel payment order
// @Description Cancels a pending payment transaction order when current business rules allow it. Allowed roles: SUPER_ADMIN, SCHOOL_ADMIN, TREASURER, PARENT, STUDENT.
// @Tags Payment Order
// @Produce json
// @Param id path string true "Payment order ID"
// @Success 200 {object} domain.SwaggerPaymentOrderResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 403 {object} domain.ErrorResponse
// @Failure 404 {object} domain.ErrorResponse
// @Failure 422 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Security BearerAuth
// @Router /payment-orders/{id}/cancel [post]
func (h *paymentOrderHandler) Cancel(c *gin.Context) {
	authContextValue, _ := c.Get(domain.ContextKeyAuth)
	authContext := authContextValue.(*security.AuthContext)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": shared.ErrInvalidIDParam.Error()})
		return
	}

	paymentOrder, err := h.paymentOrderService.Cancel(c.Request.Context(), authContext, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := domain.NewPaymentOrderResponse(paymentOrder)
	c.JSON(http.StatusOK, gin.H{"data": res})
}
