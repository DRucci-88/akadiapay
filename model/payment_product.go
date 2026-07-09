package model

import (
	"github.com/google/uuid"
)

type PaymentProduct struct {
	TenantID uuid.UUID `gorm:"type:uuid;not null;index:idx_payment_product_tenant;uniqueIndex:uk_payment_product"`
	// Tenant   *Tenant   `gorm:"foreignKey:TenantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	PaymentPolicyID uuid.UUID      `gorm:"type:uuid;not null;index:idx_payment_product_policy"`
	PaymentPolicy   *PaymentPolicy `gorm:"foreignKey:PaymentPolicyID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	Code        string `gorm:"type:varchar(50);not null;uniqueIndex:uk_payment_product_code;uniqueIndex:uk_payment_product"`
	Name        string `gorm:"type:varchar(150);not null"`
	Description string `gorm:"type:text"`

	Price  float64              `gorm:"type:numeric(18,2);not null"`
	Status PaymentProductStatus `gorm:"type:varchar(20);default:ACTIVE;not null"`

	BaseModel
}

func (PaymentProduct) TableName() string {
	return SchemaPayment + ".payment_products"
}

type PaymentProductStatus string

const (
	PaymentProductStatusActive   PaymentProductStatus = "ACTIVE"
	PaymentProductStatusInactive PaymentProductStatus = "INACTIVE"
)

type PaymentProductPreload string

const (
	PaymentProductPreloadPaymentPolicy PaymentProductPreload = "PaymentPolicy"
)
