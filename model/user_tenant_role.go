package model

import "github.com/google/uuid"

type UserTenantRole struct {
	UserID uuid.UUID `gorm:"type:uuid;not null;index:idx_user_tenant_role_user;uniqueIndex:uk_user_tenant_role"`
	User   *User     `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	TenantID uuid.UUID `gorm:"type:uuid;not null;index:idx_user_tenant_role_tenant;uniqueIndex:uk_user_tenant_role"`
	Tenant   *Tenant   `gorm:"foreignKey:TenantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	RoleID uuid.UUID `gorm:"type:uuid;not null;index:idx_user_tenant_role_role;uniqueIndex:uk_user_tenant_role"`
	Role   *Role     `gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	IsDefault bool `gorm:"default:false;not null"`
	IsActive  bool `gorm:"default:true;not null"`

	BaseModel
}

func (UserTenantRole) TableName() string {
	return SchemaMaster + ".user_tenant_roles"
}

type UserTenantRolePreload string

const (
	UserTenantRolePreloadUser   UserTenantRolePreload = "User"
	UserTenantRolePreloadTenant UserTenantRolePreload = "Tenant"
	UserTenantRolePreloadRole   UserTenantRolePreload = "Role"
)

/*
Why IsDefault?

Multi Tenant User
Imagine this user.
John
↓
School A
↓
Treasurer
----------------
School B
↓
Parent
----------------
School C
↓
School Admin
*/
