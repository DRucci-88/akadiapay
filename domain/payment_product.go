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

type PaymentProductHandler interface {
	FindAll(c *gin.Context)
	FindByID(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
}

type PaymentProductService interface {
	FindPaginate(
		ctx context.Context,
		pageable *shared.Pageable,
		filter *PaymentProductFilter,
		authContext *security.AuthContext,
	) (*shared.Page[PaymentProductResponse], error)
	FindByID(
		ctx context.Context,
		authContext *security.AuthContext,
		id uuid.UUID,
		preloads ...model.PaymentProductPreload,
	) (*model.PaymentProduct, error)
	Create(
		ctx context.Context,
		authContext *security.AuthContext,
		req *PaymentProductCreate,
	) (*model.PaymentProduct, error)
	Update(
		ctx context.Context,
		authContext *security.AuthContext,
		id uuid.UUID,
		req *PaymentProductUpdate,
	) (*model.PaymentProduct, error)
}

type PaymentProductRepository interface {
	QueryWithPreloads(
		preloads ...model.PaymentProductPreload,
	) gorm.ChainInterface[model.PaymentProduct]
	Paginate(
		ctx context.Context,
		pageable *shared.Pageable,
		chain gorm.ChainInterface[model.PaymentProduct],
	) (*shared.Page[model.PaymentProduct], error)
	FindPaginate(
		ctx context.Context,
		pageable *shared.Pageable,
		filter *PaymentProductFilter,
		authContext *security.AuthContext,
	) (*shared.Page[model.PaymentProduct], error)
	FindByID(
		ctx context.Context,
		id uuid.UUID,
		tenantID uuid.UUID,
		preloads ...model.PaymentProductPreload,
	) (*model.PaymentProduct, error)
	FindByIDsIncludingDeleted(
		ctx context.Context,
		ids []uuid.UUID,
	) ([]model.PaymentProduct, error)
	Create(
		ctx context.Context,
		paymentProduct *model.PaymentProduct,
	) error
	Update(
		ctx context.Context,
		id uuid.UUID,
		tenantID uuid.UUID,
		req *PaymentProductUpdate,
	) (int, error)
}

type PaymentProductFilter struct {
	Keyword *string `form:"keyword"`

	PaymentPolicyID *uuid.UUID                  `form:"payment_policy_id"`
	Status          *model.PaymentProductStatus `form:"status"`

	SortBy *string `form:"sort_by,default=created_at"`
	Order  *string `form:"order,default=desc"`
}

type PaymentProductCreate struct {
	PaymentPolicyID    uuid.UUID                  `json:"payment_policy_id" binding:"required"`
	Code               string                     `json:"code" binding:"required,max=50"`
	Name               string                     `json:"name" binding:"required,max=150"`
	Description        string                     `json:"description"`
	RevenueAccountCode string                     `json:"revenue_account_code" binding:"max=30"`
	RevenueAccountName string                     `json:"revenue_account_name" binding:"max=150"`
	Price              float64                    `json:"price"`
	Status             model.PaymentProductStatus `json:"status"`
}

type PaymentProductUpdate struct {
	PaymentPolicyID    *uuid.UUID                  `json:"payment_policy_id"`
	Code               *string                     `json:"code" binding:"max=50"`
	Name               *string                     `json:"name" binding:"max=150"`
	Description        *string                     `json:"description"`
	RevenueAccountCode *string                     `json:"revenue_account_code" binding:"max=30"`
	RevenueAccountName *string                     `json:"revenue_account_name" binding:"max=150"`
	Price              *float64                    `json:"price"`
	Status             *model.PaymentProductStatus `json:"status"`
}

type PaymentProductResponse struct {
	ID                 uuid.UUID                  `json:"id"`
	TenantID           uuid.UUID                  `json:"tenant_id"`
	PaymentPolicyID    uuid.UUID                  `json:"payment_policy_id"`
	Code               string                     `json:"code"`
	Name               string                     `json:"name"`
	Description        string                     `json:"description"`
	RevenueAccountCode string                     `json:"revenue_account_code"`
	RevenueAccountName string                     `json:"revenue_account_name"`
	Price              float64                    `json:"price"`
	Status             model.PaymentProductStatus `json:"status"`
}

func NewPaymentProductResponse(model *model.PaymentProduct) *PaymentProductResponse {
	return &PaymentProductResponse{
		ID:                 model.ID,
		TenantID:           model.TenantID,
		PaymentPolicyID:    model.PaymentPolicyID,
		Code:               model.Code,
		Name:               model.Name,
		Description:        model.Description,
		RevenueAccountCode: model.RevenueAccountCode,
		RevenueAccountName: model.RevenueAccountName,
		Price:              model.Price,
		Status:             model.Status,
	}
}

func NewPaymentProductResponses(models []model.PaymentProduct) []PaymentProductResponse {
	resList := make([]PaymentProductResponse, 0, len(models))
	for i := range models {
		resList = append(resList, *NewPaymentProductResponse(&models[i]))
	}
	return resList
}
