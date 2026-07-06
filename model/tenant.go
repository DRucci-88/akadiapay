package model

type Tenant struct {
	Code    string       `gorm:"type:varchar(30);uniqueIndex:uk_tenant_code;not null"`
	Name    string       `gorm:"type:varchar(150);not null"`
	Email   string       `gorm:"type:varchar(150);not null"`
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
	TenantStatusActive   TenantStatus = "ACTIVE"
	TenantStatusDeactive TenantStatus = "DEACTIVE"
)

/*
Tenant Code
*/
