package domain

import (
	"akadia/model"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentProductHandler interface {
	FindAll(c *gin.Context)
	FindByID(c *gin.Context)
}

type PaymentProductService interface {
	FindByID(
		ctx context.Context,
		id uuid.UUID,
		preloads ...model.PaymentProductPreload,
	) (*model.PaymentProduct, error)
}

type PaymentProductRepository interface {
	QueryWithPreloads(
		preloads ...model.PaymentProductPreload,
	) gorm.ChainInterface[model.PaymentProduct]
	FindByID(
		ctx context.Context,
		id uuid.UUID,
		preloads ...model.PaymentProductPreload,
	) (*model.PaymentProduct, error)
}

type PaymentProductResponse struct {
	ID              uuid.UUID                  `json:"id"`
	TenantID        uuid.UUID                  `json:"tenant_id"`
	PaymentPolicyID uuid.UUID                  `json:"payment_policy_id"`
	Code            string                     `json:"code"`
	Name            string                     `json:"name"`
	Description     string                     `json:"description"`
	Price           float64                    `json:"price"`
	Status          model.PaymentProductStatus `json:"status"`
}

func NewPaymentProductResponse(model *model.PaymentProduct) *PaymentProductResponse {
	return &PaymentProductResponse{
		ID:              model.ID,
		TenantID:        model.TenantID,
		PaymentPolicyID: model.PaymentPolicyID,

		Code:            model.Code,
		Name:            model.Name,
		Description:     model.Description,
		
		Price:           model.Price,
		Status:          model.Status,
	}
}
