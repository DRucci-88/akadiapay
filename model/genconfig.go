package model

import (
	"gorm.io/cli/gorm/field"
	"gorm.io/cli/gorm/genconfig"
)

var _ = genconfig.Config{
	FieldTypeMap: map[any]any{
		// AttendanceStatus(""): field.String{},
		// EmployeeStatus(""):   field.String{},
		// LeaveStatus(""):      field.String{},
		// UserRole(""):         field.String{},

		ParentStudentRelationship(""): field.String{},
		PaymentOrderStatus(""):        field.String{},
		PaymentOrderPaymentMethod(""): field.String{},
		PaymentProductStatus(""):      field.String{},
		RoleCode(""):                  field.String{},
		StudentObligationStatus(""):   field.String{},
		StudentGender(""):             field.String{},
		StudentStatus(""):             field.String{},
		TenantStatus(""):              field.String{},
		UserStatus(""):                field.String{},
	},
	// FieldNameMap: map[string]any{
	// 	// "status": field.String{},
	// },
}
