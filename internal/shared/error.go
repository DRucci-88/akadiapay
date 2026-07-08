package shared

import "errors"

var (
	/// Auth
	ErrAuthLogin        = errors.New("Email or Password is wrong")
	ErrAuthUnauthorized = errors.New("Unauthorized")
	ErrAuthTokenExpired = errors.New("Token is Expired")

	/// User
	ErrUserNotFound = errors.New("User not found")

	/// Student
	ErrStudentNotFound = errors.New("Student not found")

	/// UserTetantRole
	ErrUserTenantRoleNotFound = errors.New("User Tenant Role not found")

	/// Tenant
	ErrTenantNotFound = errors.New("Tenant Not Found")
)
