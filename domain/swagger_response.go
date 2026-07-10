package domain

import "akadia/internal/shared"

type ErrorResponse struct {
	Error string `json:"error" example:"validation error message"`
}

type SwaggerHealthResponse struct {
	Status string `json:"status" example:"Berjalan Perfecto"`
	AAA    string `json:"aaa" example:"code"`
}

type SwaggerBoolResponse struct {
	Data bool `json:"data"`
}

type SwaggerAuthLoginResponse struct {
	Data []AuthLoginResponse `json:"data"`
}

type SwaggerAuthProfileResponse struct {
	Data AuthProfileResponse `json:"data"`
}

type SwaggerPaymentPolicyResponse struct {
	Data PaymentPolicyResponse `json:"data"`
}

type SwaggerPaymentProductResponse struct {
	Data PaymentProductResponse `json:"data"`
}

type SwaggerStudentObligationResponse struct {
	Data StudentObligationResponse `json:"data"`
}

type SwaggerStudentObligationListResponse struct {
	Data []StudentObligationResponse `json:"data"`
}

type SwaggerStudentOutstandingResponse struct {
	Data StudentOutstandingResponse `json:"data"`
}

type SwaggerPaymentOrderResponse struct {
	Data PaymentOrderResponse `json:"data"`
}

type SwaggerPaymentAllocationResultResponse struct {
	Data PaymentAllocationResult `json:"data"`
}

type SwaggerLedgerEntryListResponse struct {
	Data []LedgerEntryResponse `json:"data"`
}

type SwaggerPaymentPolicyPageResponse struct {
	Data       []PaymentPolicyResponse `json:"data"`
	Pagination shared.Pagination       `json:"pagination"`
}

type SwaggerPaymentProductPageResponse struct {
	Data       []PaymentProductResponse `json:"data"`
	Pagination shared.Pagination        `json:"pagination"`
}

type SwaggerStudentObligationPageResponse struct {
	Data       []StudentObligationResponse `json:"data"`
	Pagination shared.Pagination           `json:"pagination"`
}

type SwaggerPaymentOrderPageResponse struct {
	Data       []PaymentOrderResponse `json:"data"`
	Pagination shared.Pagination      `json:"pagination"`
}

type SwaggerLedgerEntryPageResponse struct {
	Data       []LedgerEntryResponse `json:"data"`
	Pagination shared.Pagination     `json:"pagination"`
}
