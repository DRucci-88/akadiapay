package model

import (
	"time"

	"github.com/google/uuid"
)

type PaymentOrder struct {
	TenantID uuid.UUID `gorm:"type:uuid;not null;index:idx_payment_order_tenant"`
	// Tenant   *Tenant   `gorm:"foreignKey:TenantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	StudentID uuid.UUID `gorm:"type:uuid;not null;index:idx_payment_order_student"`
	// Student   *Student  `gorm:"foreignKey:StudentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	PaidByUserID uuid.UUID `gorm:"type:uuid;not null;index:idx_payment_order_paid_by_user"`
	// PaidByUser   *User  `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	OrderNumber string    `gorm:"type:varchar(50);not null;uniqueIndex:uk_payment_order_number"`
	OrderDate   time.Time `gorm:"not null"`

	TotalAmount     float64            `gorm:"type:numeric(18,2);not null"`
	Status          PaymentOrderStatus `gorm:"type:varchar(20);default:PENDING;not null"`
	PaymentMethod   PaymentMethod      `gorm:"type:varchar(30);not null"`
	ReferenceNumber *string            `gorm:"type:varchar(100)"`
	Notes           string             `gorm:"type:text"`

	BaseModel

	PaymentAllocations []PaymentAllocation `gorm:"foreignKey:PaymentOrderID"`
}

func (PaymentOrder) TableName() string {
	return SchemaPayment + ".payment_orders"
}

type PaymentOrderStatus string

const (
	PaymentOrderStatusPending   PaymentOrderStatus = "PENDING"
	PaymentOrderStatusCompleted PaymentOrderStatus = "COMPLETED"
	PaymentOrderStatusCancelled PaymentOrderStatus = "CANCELLED"
	PaymentOrderStatusExpired   PaymentOrderStatus = "EXPIRED"
)

type PaymentMethod string

const (
	PaymentMethodCash           PaymentMethod = "CASH"
	PaymentMethodBankTransfer   PaymentMethod = "BANK_TRANSFER"
	PaymentMethodVirtualAccount PaymentMethod = "VIRTUAL_ACCOUNT"
	PaymentMethodQRIS           PaymentMethod = "QRIS"
	PaymentMethodCreditCard     PaymentMethod = "CREDIT_CARD"
)
