package model

type Role struct {
	Code        RoleCode `gorm:"type:varchar(50);not null;uniqueIndex:uk_role_code"`
	Name        string   `gorm:"type:varchar(100);not null"`
	Description *string  `gorm:"type:text"`
	IsSystem    bool     `gorm:"default:true;not null"`

	BaseModel

	UserTenantRoles []UserTenantRole `gorm:"foreignKey:RoleID"`
}

func (Role) TableName() string {
	return SchemaMaster + ".roles"
}

type RoleCode string

const (
	RoleCodeSuperAdmin  RoleCode = "SUPER_ADMIN"
	RoleCodeSchoolAdmin RoleCode = "SCHOOL_ADMIN"
	RoleCodeTreasurer   RoleCode = "TREASURER"
	RoleCodeTeacher     RoleCode = "TEACHER"
	RoleCodeParent      RoleCode = "PARENT"
	RoleCodeStudent     RoleCode = "STUDENT"
)

/*
Only Seeder not CRUD
Code used by system
Name used by UI
*/
