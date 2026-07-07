package model

import "github.com/google/uuid"

type PaymentAllocation struct {
	PaymentOrderID uuid.UUID     `gorm:"type:uuid;not null;index:idx_payment_allocation_order;uniqueIndex:uk_payment_allocation"`
	PaymentOrder   *PaymentOrder `gorm:"foreignKey:PaymentOrderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	StudentObligationID uuid.UUID          `gorm:"type:uuid;not null;index:idx_payment_allocation_obligation;uniqueIndex:uk_payment_allocation"`
	StudentObligation   *StudentObligation `gorm:"foreignKey:StudentObligationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	AllocatedAmount float64 `gorm:"type:numeric(18,2);not null"`

	BaseModel
}

func (PaymentAllocation) TableName() string {
	return SchemaPayment + ".payment_allocations"
}
