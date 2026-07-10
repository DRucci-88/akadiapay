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

	TotalAmount     float64                   `gorm:"type:numeric(18,2);not null"`
	Status          PaymentOrderStatus        `gorm:"type:varchar(20);default:PENDING;not null"`
	PaymentMethod   PaymentOrderPaymentMethod `gorm:"type:varchar(30);not null"`
	ReferenceNumber *string                   `gorm:"type:varchar(100)"`
	LedgerPostedAt  *time.Time                `gorm:"index"`
	Notes           string                    `gorm:"type:text"`

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

type PaymentOrderPaymentMethod string

const (
	PaymentMethodCash           PaymentOrderPaymentMethod = "CASH"
	PaymentMethodBankTransfer   PaymentOrderPaymentMethod = "BANK_TRANSFER"
	PaymentMethodVirtualAccount PaymentOrderPaymentMethod = "VIRTUAL_ACCOUNT"
	PaymentMethodQRIS           PaymentOrderPaymentMethod = "QRIS"
	PaymentMethodCreditCard     PaymentOrderPaymentMethod = "CREDIT_CARD"
)
