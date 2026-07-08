package model

import "github.com/google/uuid"

type ParentStudent struct {
	BaseModel

	ParentUserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:uk_parent_student"`
	ParentUser   *User     `gorm:"foreignKey:ParentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	StudentID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:uk_parent_student"`
	Student   *Student  `gorm:"foreignKey:StudentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	Relationship string `gorm:"type:varchar(30);not null"`

	IsPrimary bool `gorm:"default:false;not null"`
}

func (ParentStudent) TableName() string {
	return SchemaMaster + ".parent_students"
}

type ParentStudentRelationship string

const (
	ParentStudentRelationshipFather   = "FATHER"
	ParentStudentRelationshipMother   = "MOTHER"
	ParentStudentRelationshipGuardian = "GUARDIAN"
)

/*
Why IsPrimary?

Example: Father & Mother

Who receives
- Email
- WhatsApp
- Push Notification

by default? IsPrimary = true

Exactly.
*/
