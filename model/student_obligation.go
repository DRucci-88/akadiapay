package model

import (
	"time"

	"github.com/google/uuid"
)

type StudentObligation struct {
	StudentID uuid.UUID `gorm:"type:uuid;not null;index:idx_student_obligation_student;uniqueIndex:uk_student_obligation"`
	// Student   *Student  `gorm:"foreignKey:StudentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	PaymentProductID uuid.UUID       `gorm:"type:uuid;not null;index:idx_student_obligation_product;uniqueIndex:uk_student_obligation"`
	PaymentProduct   *PaymentProduct `gorm:"foreignKey:PaymentProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	Period            time.Time `gorm:"type:date;uniqueIndex:uk_student_obligation"`
	OriginalAmount    float64   `gorm:"type:numeric(18,2);not null"`
	OutstandingAmount float64   `gorm:"type:numeric(18,2);not null"`
	DueDate           time.Time
	IssuedAt          time.Time
	Status            StudentObligationStatus `gorm:"type:varchar(20);default:PENDING;not null"`

	BaseModel

	PaymentAllocations []PaymentAllocation `gorm:"foreignKey:StudentObligationID"`
}

func (StudentObligation) TableName() string {
	return SchemaPayment + ".student_obligations"
}

type StudentObligationStatus string

const (
	StudentObligationStatusPending   StudentObligationStatus = "PENDING"
	StudentObligationStatusPartial   StudentObligationStatus = "PARTIAL"
	StudentObligationStatusPaid      StudentObligationStatus = "PAID"
	StudentObligationStatusCancelled StudentObligationStatus = "CANCELLED"
)

/*
Why OutstandingAmount?
Changes every payment.
*/
