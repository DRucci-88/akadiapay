package domain

import (
	"akadia/internal/platform/security"
	"akadia/internal/shared"
	"akadia/model"
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StudentObligationHandler interface {
	Create(c *gin.Context)
	FindAll(c *gin.Context)
}

type StudentObligationService interface {
	Create(
		ctx context.Context,
		authContext *security.AuthContext,
		req *StudentObligationCreate,
	) (*model.StudentObligation, error)
	FindPaginate(
		ctx context.Context,
		pageable *shared.Pageable,
		filter *StudentObligationFilter,
		authContext *security.AuthContext,
	) (*shared.Page[StudentObligationResponse], error)
}

type StudentObligationRepository interface {
	Create(
		ctx context.Context,
		studentObligation *model.StudentObligation,
	) error
	Paginate(
		ctx context.Context,
		pageable *shared.Pageable,
		chain gorm.ChainInterface[model.StudentObligation],
	) (*shared.Page[model.StudentObligation], error)
	FindPaginate(
		ctx context.Context,
		pageable *shared.Pageable,
		filter *StudentObligationFilter,
		authContext *security.AuthContext,
	) (*shared.Page[model.StudentObligation], error)
	ExistsActiveByStudentIDAndPaymentProductIDAndPeriod(
		ctx context.Context,
		studentID uuid.UUID,
		paymentProductID uuid.UUID,
		period time.Time,
	) (bool, error)
}

type StudentObligationFilter struct {
	Keyword *string `form:"keyword"`

	StudentID        *uuid.UUID                     `form:"student_id"`
	PaymentProductID *uuid.UUID                     `form:"payment_product_id"`
	Status           *model.StudentObligationStatus `form:"status"`
	DueDateFrom      *time.Time                     `form:"due_date_from" time_format:"2006-01-02"`
	DueDateTo        *time.Time                     `form:"due_date_to" time_format:"2006-01-02"`

	SortBy *string `form:"sort_by,default=created_at"`
	Order  *string `form:"order,default=desc"`
}

type StudentObligationCreate struct {
	StudentID        uuid.UUID `json:"student_id" binding:"required"`
	PaymentProductID uuid.UUID `json:"payment_product_id" binding:"required"`
	DueDate          time.Time `json:"due_date" binding:"required"`
	Amount           *float64  `json:"amount"`
	Notes            string    `json:"notes"`
}

type StudentObligationResponse struct {
	ID                uuid.UUID                     `json:"id"`
	StudentID         uuid.UUID                     `json:"student_id"`
	PaymentProductID  uuid.UUID                     `json:"payment_product_id"`
	Period            time.Time                     `json:"period"`
	OriginalAmount    float64                       `json:"original_amount"`
	PaidAmount        float64                       `json:"paid_amount"`
	OutstandingAmount float64                       `json:"outstanding_amount"`
	DueDate           time.Time                     `json:"due_date"`
	IssuedAt          time.Time                     `json:"issued_at"`
	Status            model.StudentObligationStatus `json:"status"`
	Notes             string                        `json:"notes"`
}

func NewStudentObligationResponse(model *model.StudentObligation) *StudentObligationResponse {
	return &StudentObligationResponse{
		ID:                model.ID,
		StudentID:         model.StudentID,
		PaymentProductID:  model.PaymentProductID,
		Period:            model.Period,
		OriginalAmount:    model.OriginalAmount,
		PaidAmount:        model.OriginalAmount - model.OutstandingAmount,
		OutstandingAmount: model.OutstandingAmount,
		DueDate:           model.DueDate,
		IssuedAt:          model.IssuedAt,
		Status:            model.Status,
		Notes:             model.Notes,
	}
}

func NewStudentObligationResponses(models []model.StudentObligation) []StudentObligationResponse {
	resList := make([]StudentObligationResponse, 0, len(models))
	for i := range models {
		resList = append(resList, *NewStudentObligationResponse(&models[i]))
	}
	return resList
}
