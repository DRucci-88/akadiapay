package handler

import (
	"akadia/domain"
	"akadia/internal/platform/security"
	"akadia/internal/shared"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

// Create godoc
// @Summary Create student obligation
// @Description Creates one student financial obligation record with outstanding, paid, closed, or cancelled lifecycle semantics driven by the existing billing flow. Allowed roles: SUPER_ADMIN, SCHOOL_ADMIN, TREASURER.
// @Tags Student Obligation
// @Accept json
// @Produce json
// @Param request body domain.StudentObligationCreate true "Student obligation payload"
// @Success 200 {object} domain.SwaggerStudentObligationResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 403 {object} domain.ErrorResponse
// @Failure 422 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Security BearerAuth
// @Router /student-obligations [post]
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

// CreateBulk godoc
// @Summary Create student obligations in bulk
// @Description Assigns one payment product to multiple students and creates obligation rows in bulk using the current billing rules. Allowed roles: SUPER_ADMIN, SCHOOL_ADMIN, TREASURER.
// @Tags Student Obligation
// @Accept json
// @Produce json
// @Param request body domain.StudentObligationBulkCreate true "Bulk student obligation payload"
// @Success 200 {object} domain.SwaggerStudentObligationListResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 403 {object} domain.ErrorResponse
// @Failure 422 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Security BearerAuth
// @Router /student-obligations/bulk [post]
func (h *studentObligationHandler) CreateBulk(c *gin.Context) {
	authContextValue, _ := c.Get(domain.ContextKeyAuth)
	authContext := authContextValue.(*security.AuthContext)

	var req domain.StudentObligationBulkCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	studentObligations, err := h.studentObligationService.CreateBulk(c.Request.Context(), authContext, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := domain.NewStudentObligationResponses(studentObligations)
	c.JSON(http.StatusOK, gin.H{"data": res})
}

// FindAll godoc
// @Summary List student obligations
// @Description Lists student financial obligations with outstanding, paid, closed, and cancelled states using the current billing and tenant isolation rules. Allowed roles: SUPER_ADMIN, SCHOOL_ADMIN, TREASURER.
// @Tags Student Obligation
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param size query int false "Page size (default 15, max 100)"
// @Param keyword query string false "Search keyword for student or payment product"
// @Param student_id query string false "Filter by student ID"
// @Param payment_product_id query string false "Filter by payment product ID"
// @Param status query string false "Filter by obligation status: PENDING, PARTIAL, PAID, CLOSED, CANCELLED"
// @Param due_date_from query string false "Filter due date from (YYYY-MM-DD)"
// @Param due_date_to query string false "Filter due date to (YYYY-MM-DD)"
// @Param sort_by query string false "Sort field: created_at, updated_at, student_id, payment_product_id, period, original_amount, outstanding_amount, due_date, issued_at, status"
// @Param order query string false "Sort direction: asc or desc"
// @Success 200 {object} domain.SwaggerStudentObligationPageResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 403 {object} domain.ErrorResponse
// @Failure 422 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Security BearerAuth
// @Router /student-obligations [get]
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

// FindByID godoc
// @Summary Get student obligation by ID
// @Description Returns one student financial obligation record including outstanding balance and current status. Allowed roles: SUPER_ADMIN, SCHOOL_ADMIN, TREASURER.
// @Tags Student Obligation
// @Produce json
// @Param id path string true "Student obligation ID"
// @Success 200 {object} domain.SwaggerStudentObligationResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 403 {object} domain.ErrorResponse
// @Failure 404 {object} domain.ErrorResponse
// @Failure 422 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Security BearerAuth
// @Router /student-obligations/{id} [get]
func (h *studentObligationHandler) FindByID(c *gin.Context) {
	authContextValue, _ := c.Get(domain.ContextKeyAuth)
	authContext := authContextValue.(*security.AuthContext)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": shared.ErrInvalidIDParam.Error()})
		return
	}

	studentObligation, err := h.studentObligationService.FindByID(c.Request.Context(), authContext, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := domain.NewStudentObligationResponse(studentObligation)
	c.JSON(http.StatusOK, gin.H{"data": res})
}

// Update godoc
// @Summary Update student obligation
// @Description Updates mutable student obligation fields without changing the billing flow semantics for outstanding, paid, closed, or cancelled states. Allowed roles: SUPER_ADMIN, SCHOOL_ADMIN, TREASURER.
// @Tags Student Obligation
// @Accept json
// @Produce json
// @Param id path string true "Student obligation ID"
// @Param request body domain.StudentObligationUpdate true "Student obligation update payload"
// @Success 200 {object} domain.SwaggerStudentObligationResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 403 {object} domain.ErrorResponse
// @Failure 404 {object} domain.ErrorResponse
// @Failure 422 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Security BearerAuth
// @Router /student-obligations/{id} [put]
func (h *studentObligationHandler) Update(c *gin.Context) {
	authContextValue, _ := c.Get(domain.ContextKeyAuth)
	authContext := authContextValue.(*security.AuthContext)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": shared.ErrInvalidIDParam.Error()})
		return
	}

	var req domain.StudentObligationUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	studentObligation, err := h.studentObligationService.Update(c.Request.Context(), authContext, id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := domain.NewStudentObligationResponse(studentObligation)
	c.JSON(http.StatusOK, gin.H{"data": res})
}

// Delete godoc
// @Summary Delete student obligation
// @Description Deletes a student obligation only when current business rules allow it, without altering existing billing behavior. Allowed roles: SUPER_ADMIN, SCHOOL_ADMIN, TREASURER.
// @Tags Student Obligation
// @Produce json
// @Param id path string true "Student obligation ID"
// @Success 200 {object} domain.SwaggerBoolResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 403 {object} domain.ErrorResponse
// @Failure 404 {object} domain.ErrorResponse
// @Failure 422 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Security BearerAuth
// @Router /student-obligations/{id} [delete]
func (h *studentObligationHandler) Delete(c *gin.Context) {
	authContextValue, _ := c.Get(domain.ContextKeyAuth)
	authContext := authContextValue.(*security.AuthContext)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": shared.ErrInvalidIDParam.Error()})
		return
	}

	if err := h.studentObligationService.Delete(c.Request.Context(), authContext, id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}

// OutstandingByStudentID godoc
// @Summary Get outstanding obligations by student
// @Description Returns outstanding student financial obligations visible to the current finance portal user. Allowed roles: SUPER_ADMIN, SCHOOL_ADMIN, TREASURER, PARENT, STUDENT.
// @Tags Student Outstanding
// @Produce json
// @Param studentId path string true "Student ID"
// @Success 200 {object} domain.SwaggerStudentOutstandingResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 403 {object} domain.ErrorResponse
// @Failure 404 {object} domain.ErrorResponse
// @Failure 422 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Security BearerAuth
// @Router /students/{studentId}/outstanding [get]
func (h *studentObligationHandler) OutstandingByStudentID(c *gin.Context) {
	authContextValue, _ := c.Get(domain.ContextKeyAuth)
	authContext := authContextValue.(*security.AuthContext)

	studentID, err := uuid.Parse(c.Param("studentId"))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": shared.ErrInvalidIDParam.Error()})
		return
	}

	res, err := h.studentObligationService.FindOutstandingByStudentID(c.Request.Context(), authContext, studentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}
