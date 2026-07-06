package model

type User struct {
	Email    string     `gorm:"type:varchar(150);not null;uniqueIndex:uk_user_email"`
	Password string     `gorm:"type:text;not null"`
	FullName string     `gorm:"type:varchar(150);not null"`
	Phone    string     `gorm:"type:varchar(30)"`
	Status   UserStatus `gorm:"type:varchar(20);default:ACTIVE;not null"`

	BaseModel

	UserTenantRoles []UserTenantRole `gorm:"foreignKey:UserID"`
}

func (User) TableName() string {
	return SchemaMaster + ".users"
}

type UserStatus string

const (
	UserStatusActive   = "ACTIVE"
	UserStatusInactive = "INACTIVE"
)
