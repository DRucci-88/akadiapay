package model

import (
	"github.com/google/uuid"
)

// Memisahkan aturan umum dengan tipe data ketat dan menyediakan kolom dinamis
type PaymentPolicy struct {
	TenantID uuid.UUID `gorm:"type:uuid;not null;index:idx_payment_policy_tenant;uniqueIndex:uk_payment_policy"`

	Code        string `gorm:"type:varchar(50);not null;uniqueIndex:uk_payment_policy_code;uniqueIndex:uk_payment_policy"`
	Name        string `gorm:"type:varchar(100);not null"`
	Description string `gorm:"type:text"`

	AllowPartial      bool    `gorm:"default:false;not null"`
	MinimumAmount     float64 `gorm:"type:numeric(18,2);default:0;not null"`
	MinimumPercentage float64 `gorm:"type:numeric(5,2);default:0;not null"`

	AllowOverPayment    bool `gorm:"default:false;not null"`
	AutoCloseObligation bool `gorm:"default:true;not null"`

	BaseModel

	PaymentProducts []PaymentProduct `gorm:"foreignKey:PaymentPolicyID"`
}

func (PaymentPolicy) TableName() string {
	return SchemaPayment + ".payment_policies"
}
