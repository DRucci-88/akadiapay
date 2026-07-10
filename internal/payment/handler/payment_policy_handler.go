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
