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

type PaymentOrderHandler interface {
	FindAll(c *gin.Context)
	FindByID(c *gin.Context)
	Create(c *gin.Context)
	Cancel(c *gin.Context)
}

type PaymentOrderService interface {
	FindPaginate(
		ctx context.Context,
		pageable *shared.Pageable,
		filter *PaymentOrderFilter,
		authContext *security.AuthContext,
	) (*shared.Page[PaymentOrderResponse], error)
	FindByID(
		ctx context.Context,
		authContext *security.AuthContext,
		id uuid.UUID,
	) (*model.PaymentOrder, error)
	Create(
		ctx context.Context,
		authContext *security.AuthContext,
		req *PaymentOrderCreate,
	) (*model.PaymentOrder, error)
	Cancel(
		ctx context.Context,
		authContext *security.AuthContext,
		id uuid.UUID,
	) (*model.PaymentOrder, error)
}

type PaymentOrderRepository interface {
	Paginate(
		ctx context.Context,
		pageable *shared.Pageable,
		chain gorm.ChainInterface[model.PaymentOrder],
	) (*shared.Page[model.PaymentOrder], error)
	FindPaginate(
		ctx context.Context,
		pageable *shared.Pageable,
		filter *PaymentOrderFilter,
		authContext *security.AuthContext,
	) (*shared.Page[model.PaymentOrder], error)
	Create(
		ctx context.Context,
		paymentOrder *model.PaymentOrder,
	) error
	FirstByID(
		ctx context.Context,
		id uuid.UUID,
	) (*model.PaymentOrder, error)
	UpdateStatus(
		ctx context.Context,
		id uuid.UUID,
		status model.PaymentOrderStatus,
	) (int, error)
}

type PaymentOrderFilter struct {
	Keyword *string `form:"keyword"`

	StudentID     *uuid.UUID                       `form:"student_id"`
	Status        *model.PaymentOrderStatus        `form:"status"`
	PaymentMethod *model.PaymentOrderPaymentMethod `form:"payment_method"`
	OrderDateFrom *time.Time                       `form:"order_date_from" time_format:"2006-01-02"`
	OrderDateTo   *time.Time                       `form:"order_date_to" time_format:"2006-01-02"`

	SortBy *string `form:"sort_by,default=created_at"`
	Order  *string `form:"order,default=desc"`
}

type PaymentOrderCreate struct {
	StudentID     uuid.UUID                       `json:"student_id" binding:"required"`
	Amount        float64                         `json:"amount"`
	PaymentMethod model.PaymentOrderPaymentMethod `json:"payment_method" binding:"required"`
	PaymentDate   time.Time                       `json:"payment_date" binding:"required"`
	Notes         string                          `json:"notes"`
}

type PaymentOrderResponse struct {
	ID              uuid.UUID                       `json:"id"`
	TenantID        uuid.UUID                       `json:"tenant_id"`
	StudentID       uuid.UUID                       `json:"student_id"`
	PaidByUserID    uuid.UUID                       `json:"paid_by_user_id"`
	OrderNumber     string                          `json:"order_number"`
	OrderDate       time.Time                       `json:"order_date"`
	TotalAmount     float64                         `json:"total_amount"`
	Status          model.PaymentOrderStatus        `json:"status"`
	PaymentMethod   model.PaymentOrderPaymentMethod `json:"payment_method"`
	ReferenceNumber *string                         `json:"reference_number,omitempty"`
	Notes           string                          `json:"notes"`
}

func NewPaymentOrderResponse(model *model.PaymentOrder) *PaymentOrderResponse {
	return &PaymentOrderResponse{
		ID:              model.ID,
		TenantID:        model.TenantID,
		StudentID:       model.StudentID,
		PaidByUserID:    model.PaidByUserID,
		OrderNumber:     model.OrderNumber,
		OrderDate:       model.OrderDate,
		TotalAmount:     model.TotalAmount,
		Status:          model.Status,
		PaymentMethod:   model.PaymentMethod,
		ReferenceNumber: model.ReferenceNumber,
		Notes:           model.Notes,
	}
}

func NewPaymentOrderResponses(models []model.PaymentOrder) []PaymentOrderResponse {
	resList := make([]PaymentOrderResponse, 0, len(models))
	for i := range models {
		resList = append(resList, *NewPaymentOrderResponse(&models[i]))
	}
	return resList
}
