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
	ErrStudentObligationStudentIDsEmpty = errors.New("Student Obligation student IDs required")
	ErrStudentObligationAllocated       = errors.New("Student Obligation has payment allocations")

	// Payment Order
	ErrPaymentOrderAmountInvalid            = errors.New("Payment Order amount invalid")
	ErrPaymentOrderMethodInvalid            = errors.New("Payment Order payment method invalid")
	ErrPaymentOrderDateRequired             = errors.New("Payment Order date required")
	ErrPaymentOrderOutstandingRequired      = errors.New("Outstanding bill required")
	ErrPaymentOrderAmountExceedsOutstanding = errors.New("Payment Order amount exceeds outstanding amount")
	ErrPaymentOrderNotFound                 = errors.New("Payment Order not found")
	ErrPaymentOrderStatusInvalid            = errors.New("Payment Order status invalid")
	ErrPaymentOrderAllocated                = errors.New("Payment Order has payment allocations")
	ErrPaymentAllocationNotFound            = errors.New("Payment Allocation not found")

	// Payment Allocation
	ErrPaymentAllocationRequired                 = errors.New("Payment Allocation required")
	ErrPaymentAllocationAmountInvalid            = errors.New("Payment Allocation amount invalid")
	ErrPaymentAllocationAmountExceedsOutstanding = errors.New("Payment Allocation exceeds outstanding amount")
	ErrPaymentAllocationFullPaymentRequired      = errors.New("Payment Allocation must equal outstanding amount")
	ErrPaymentAllocationBelowMinimumAmount       = errors.New("Payment Allocation below minimum amount")
	ErrPaymentAllocationBelowMinimumPercentage   = errors.New("Payment Allocation below minimum percentage")
	ErrPaymentAllocationTotalExceedsOrder        = errors.New("Payment Allocation total exceeds payment order amount")
	ErrPaymentAllocationDuplicateObligation      = errors.New("Payment Allocation duplicate obligation")

	// Ledger
	ErrLedgerAlreadyPosted              = errors.New("Ledger already posted")
	ErrLedgerEntryNotFound              = errors.New("Ledger Entry not found")
	ErrLedgerUnbalanced                 = errors.New("Ledger entries are not balanced")
	ErrLedgerUnbalancedSource           = errors.New("Payment Allocation does not match payment total")
	ErrLedgerDebitAccountNotConfigured  = errors.New("Debit account is not configured")
	ErrLedgerCreditAccountNotConfigured = errors.New("Credit account is not configured")
	ErrPostedPaymentCannotBeCancelled   = errors.New("Posted payment cannot be cancelled")
)
