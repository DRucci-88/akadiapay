package model

import (
	"time"

	"github.com/google/uuid"
)

type Student struct {
	TenantID uuid.UUID `gorm:"type:uuid;not null;index:idx_student_tenant"`
	Tenant   *Tenant   `gorm:"foreignKey:TenantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	UserID *uuid.UUID `gorm:"type:uuid;index:idx_student_user"`
	User   *User      `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	NISN       string        `gorm:"type:varchar(50);not null;uniqueIndex:uk_student_nisn"`
	FullName   string        `gorm:"type:varchar(150);not null"`
	Gender     StudentGender `gorm:"type:varchar(20)"`
	BirthPlace string        `gorm:"type:varchar(100)"`
	BirthDate  *time.Time
	Email      string        `gorm:"type:varchar(150)"`
	Phone      string        `gorm:"type:varchar(30)"`
	Address    string        `gorm:"type:text"`
	Status     StudentStatus `gorm:"type:varchar(20);not null"`

	BaseModel

	ParentStudents []ParentStudent `gorm:"foreignKey:StudentID"`
}

func (Student) TableName() string {
	return SchemaMaster + ".students"
}

type StudentGender string

const (
	StudentGenderMale   StudentGender = "MALE"
	StudentGenderFemale StudentGender = "FEMALE"
)

type StudentStatus string

const (
	StudentStatusActive     StudentStatus = "ACTIVE"
	StudentStatusGraduated  StudentStatus = "GRADUATED"
	StudentStatusDroppedOut StudentStatus = "DROPPED_OUT"
	StudentStatusSuspended  StudentStatus = "SUSPENDED"
)

/*
Why UserID is nullable?
This is intentional.
Imagine this flow.
Today
Admin imports
500 students
None of them has an account.
Perfectly valid.

Tomorrow
Parent activates account.
Student activates account.
Now
Student
↓
User
is connected.
No migration.
*/

/*
No Class ?
Belongs to Akadia Academic / Course Management
*/

/*
Why BirthPlace?

Your friend's roadmap includes: Course & Certificate
Certificates usually print

Birth Place
Birth Date
*/
