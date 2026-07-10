package handler

import (
	"akadia/domain"
	"akadia/internal/platform/security"
	"akadia/internal/shared"
	"net/http"

	"github.com/gin-gonic/gin"
)

type studentObligationHandler struct {
	studentObligationService domain.StudentObligationService
}

func NewStudentObligationHandler(
	studentObligationService domain.StudentObligationService,
) domain.StudentObligationHandler {
	return &studentObligationHandler{
		studentObligationService: studentObligationService,
	}
}

func (h *studentObligationHandler) Create(c *gin.Context) {
	authContextValue, _ := c.Get(domain.ContextKeyAuth)
	authContext := authContextValue.(*security.AuthContext)

	var req domain.StudentObligationCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	studentObligation, err := h.studentObligationService.Create(c.Request.Context(), authContext, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := domain.NewStudentObligationResponse(studentObligation)
	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h *studentObligationHandler) FindAll(c *gin.Context) {
	authContextValue, _ := c.Get(domain.ContextKeyAuth)
	authContext := authContextValue.(*security.AuthContext)

	pageable, err := shared.GetPageable(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var req domain.StudentObligationFilter
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	page, err := h.studentObligationService.FindPaginate(c.Request.Context(), pageable, &req, authContext)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, page)
}
