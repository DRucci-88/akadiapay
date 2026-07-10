package shared

import "errors"

var (
	/* GENERAL */
	ErrInvalidIDParam         = errors.New("Invalid ID param")
	ErrInvalidPaginationParam = errors.New("Invalid Pagination Param")

	/* Auth */
	ErrInvalidCredential = errors.New("Email or Password is wrong")
	ErrAuthUnauthorized  = errors.New("Unauthorized")
	ErrAuthTokenExpired  = errors.New("Token is Expired")

	/* MASTER */
	/// User
	ErrUserNotFound = errors.New("User not found")

	/// Student
	ErrStudentNotFound = errors.New("Student not found")

	/// UserTetantRole
	ErrUserTenantRoleNotFound = errors.New("User Tenant Role not found")

	/// Tenant
	ErrTenantNotFound = errors.New("Tenant Not Found")

	/* PAYMENT */
	// PaymentPolicy
	ErrPaymentPolicyNotFound                 = errors.New("Payment Policy not found")
	ErrPaymentPolicyMinimumAmountInvalid     = errors.New("Minimum Amount Invalid")
	ErrPaymentPolicyMinimumPercentageInvalid = errors.New("Minimum Percentage Invalid")
	ErrPaymentPolicyMinimumPaymentRequired   = errors.New("Minimum Payment Required")

	// Payment Product
	ErrPaymentProductNotFound      = errors.New("Payment Product not found")
	ErrPaymentProductPriceInvalid  = errors.New("Payment Product price invalid")
	ErrPaymentProductStatusInvalid = errors.New("Payment Product status invalid")

	// Student Obligation
	ErrStudentObligationNotFound        = errors.New("Student Obligation not found")
	ErrStudentObligationAlreadyExists   = errors.New("Student Obligation already exists")
	ErrStudentObligationAmountInvalid   = errors.New("Student Obligation amount invalid")
	ErrStudentObligationDueDateRequired = errors.New("Student Obligation due date required")
)
