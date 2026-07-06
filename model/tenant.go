package model

type Tenant struct {
	Code    string       `gorm:"type:varchar(30);uniqueIndex:uk_tenant_code;not null"`
	Name    string       `gorm:"type:varchar(150);not null"`
	Email   *string      `gorm:"type:varchar(150)"`
	Phone   string       `gorm:"type:varchar(30);not null"`
	Address string       `gorm:"type:text"`
	LogoURL string       `gorm:"type:text"`
	Status  TenantStatus `gorm:"type:varchar(20);default:ACTIVE;not null"`

	BaseModel

	Students        []Student        `gorm:"foreignKey:TenantID"`
	UserTenantRoles []UserTenantRole `gorm:"foreignKey:TenantID"`
}

func (Tenant) TableName() string {
	return SchemaMaster + ".tenants"
}

type TenantStatus string

const (
	TenantStatusActive    TenantStatus = "ACTIVE"
	TenantStatusInactive  TenantStatus = "INACTIVE"
	TenantStatusSuspended TenantStatus = "SUSPENDED"
	TenantStatusTrial     TenantStatus = "TRIAL"
	TenantStatusExpired   TenantStatus = "EXPIRED"
)

/*
Tenant Code => Used by system
Tenant Name => Used by UI
*/

/*
Email is nullable
Reason:

Some schools simply don't have an administrative email.
Phone is usually mandatory.
Email isn't.
So I'd probably use
*/
