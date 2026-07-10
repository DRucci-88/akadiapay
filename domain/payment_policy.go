package domain

import (
	"akadia/internal/platform/security"
	"akadia/internal/shared"
	"akadia/model"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentPolicyHandler interface {
	FindAll(c *gin.Context)
	FindByID(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
}

type PaymentPolicyService interface {
	FirstByID(
		ctx context.Context,
		authContext *security.AuthContext,
		id uuid.UUID,
	) (*model.PaymentPolicy, error)
	FindPaginate(
		ctx context.Context,
		pageable *shared.Pageable,
		filter *PaymentPolicyFilter,
		authContext *security.AuthContext,
	) (*shared.Page[PaymentPolicyResponse], error)
	Create(
		ctx context.Context,
		authContext *security.AuthContext,
		req *PaymentPolicyCreate,
	) (*model.PaymentPolicy, error)
	Update(
		ctx context.Context,
		authContext *security.AuthContext,
		id uuid.UUID,
		req *PaymentPolicyUpdate,
	) (*model.PaymentPolicy, error)
}

type PaymentPolicyRepository interface {
	Paginate(
		ctx context.Context,
		pageable *shared.Pageable,
		chain gorm.ChainInterface[model.PaymentPolicy],
	) (*shared.Page[model.PaymentPolicy], error)
	FirstByID(
		ctx context.Context,
		id uuid.UUID,
		tenantID uuid.UUID,
	) (*model.PaymentPolicy, error)
	FindPaginate(
		ctx context.Context,
		pageable *shared.Pageable,
		filter *PaymentPolicyFilter,
		authContext *security.AuthContext,
	) (*shared.Page[model.PaymentPolicy], error)
	Create(
		ctx context.Context,
		paymentPolicy *model.PaymentPolicy,
	) error
	Update(
		ctx context.Context,
		id uuid.UUID,
		tenantID uuid.UUID,
		req *PaymentPolicyUpdate,
	) (int, error)
}

type PaymentPolicyFilter struct {
	Keyword *string `form:"keyword"`

	AllowPartial        *bool `form:"allow_partial"`
	AllowOverPayment    *bool `form:"allow_over_payment"`
	AutoCloseObligation *bool `form:"auto_close_obligation"`

	SortBy *string `form:"sort_by,default=created_at"`
	Order  *string `form:"order,default=desc"`
}

type PaymentPolicyCreate struct {
	Code        string `json:"code" binding:"required,max=50"`
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description"`

	AllowPartial      bool    `json:"allow_partial"`
	MinimumAmount     float64 `json:"minimum_amount"`
	MinimumPercentage float64 `json:"minimum_percentage"`

	AllowOverPayment    bool `json:"allow_over_payment"`
	AutoCloseObligation bool `json:"auto_close_obligation"`
}

type PaymentPolicyUpdate struct {
	Code        *string `json:"code" binding:"max=50"`
	Name        *string `json:"name" binding:"max=100"`
	Description *string `json:"description"`

	AllowPartial      *bool    `json:"allow_partial"`
	MinimumAmount     *float64 `json:"minimum_amount"`
	MinimumPercentage *float64 `json:"minimum_percentage"`

	AllowOverPayment    *bool `json:"allow_over_payment"`
	AutoCloseObligation *bool `json:"auto_close_obligation"`
}

type PaymentPolicyResponse struct {
	ID          uuid.UUID `json:"id"`
	TenantID    uuid.UUID `json:"tenant_id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`

	AllowPartial      bool    `json:"allow_partial"`
	MinimumAmount     float64 `json:"minimum_amount"`
	MinimumPercentage float64 `json:"minimum_percentage"`

	AllowOverPayment    bool `json:"allow_over_payment"`
	AutoCloseObligation bool `json:"auto_close_obligation"`
}

func NewPaymentPolicyResponse(model *model.PaymentPolicy) *PaymentPolicyResponse {
	return &PaymentPolicyResponse{
		ID:       model.ID,
		TenantID: model.TenantID,

		Code:        model.Code,
		Name:        model.Name,
		Description: model.Description,

		AllowPartial:        model.AllowPartial,
		MinimumAmount:       model.MinimumAmount,
		MinimumPercentage:   model.MinimumPercentage,
		AllowOverPayment:    model.AllowOverPayment,
		AutoCloseObligation: model.AutoCloseObligation,
	}
}

func NewPaymentPolicyResponses(models []model.PaymentPolicy) []PaymentPolicyResponse {
	resList := make([]PaymentPolicyResponse, 0, len(models))
	for i := range models {
		resList = append(resList, *NewPaymentPolicyResponse(&models[i]))
	}
	return resList
}
