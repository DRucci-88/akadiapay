package domain

import (
	"akadia/internal/platform/security"
	"akadia/model"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PaymentAllocationHandler interface {
	Allocate(c *gin.Context)
	FindByPaymentOrderID(c *gin.Context)
}

type PaymentAllocationService interface {
	Allocate(
		ctx context.Context,
		authContext *security.AuthContext,
		paymentOrderID uuid.UUID,
		req *PaymentAllocationAllocate,
	) (*PaymentAllocationResult, error)
	FindByPaymentOrderID(
		ctx context.Context,
		authContext *security.AuthContext,
		paymentOrderID uuid.UUID,
	) (*PaymentAllocationResult, error)
}

type PaymentAllocationRepository interface {
	CreateBatch(
		ctx context.Context,
		paymentAllocations []model.PaymentAllocation,
	) error
	FindByPaymentOrderID(
		ctx context.Context,
		paymentOrderID uuid.UUID,
	) ([]model.PaymentAllocation, error)
}

type PaymentAllocationAllocate struct {
	Allocations []PaymentAllocationCreate `json:"allocations" binding:"required"`
}

type PaymentAllocationCreate struct {
	StudentObligationID uuid.UUID `json:"student_obligation_id" binding:"required"`
	AllocatedAmount     float64   `json:"allocated_amount"`
}

type PaymentAllocationResponse struct {
	ID                  uuid.UUID `json:"id"`
	PaymentOrderID      uuid.UUID `json:"payment_order_id"`
	StudentObligationID uuid.UUID `json:"student_obligation_id"`
	AllocatedAmount     float64   `json:"allocated_amount"`
}

type PaymentAllocationResult struct {
	PaymentOrderID  uuid.UUID                   `json:"payment_order_id"`
	TotalAllocated  float64                     `json:"total_allocated"`
	RemainingAmount float64                     `json:"remaining_amount"`
	OrderStatus     model.PaymentOrderStatus    `json:"order_status"`
	Allocations     []PaymentAllocationResponse `json:"allocations"`
}

func NewPaymentAllocationResponse(model *model.PaymentAllocation) *PaymentAllocationResponse {
	return &PaymentAllocationResponse{
		ID:                  model.ID,
		PaymentOrderID:      model.PaymentOrderID,
		StudentObligationID: model.StudentObligationID,
		AllocatedAmount:     model.AllocatedAmount,
	}
}

func NewPaymentAllocationResponses(models []model.PaymentAllocation) []PaymentAllocationResponse {
	resList := make([]PaymentAllocationResponse, 0, len(models))
	for i := range models {
		resList = append(resList, *NewPaymentAllocationResponse(&models[i]))
	}
	return resList
}
