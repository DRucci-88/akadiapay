package phase1_test

import (
	"akadia/domain"
	"akadia/internal/platform/security"
	"akadia/model"
	"context"
	"math"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	tenantID  = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	userID    = uuid.MustParse("22222222-2222-2222-2222-222222222222")
	studentID = uuid.MustParse("33333333-3333-3333-3333-333333333333")
	policyID  = uuid.MustParse("44444444-4444-4444-4444-444444444444")
	productID = uuid.MustParse("55555555-5555-5555-5555-555555555555")
)

func adminAuthContext() *security.AuthContext {
	return &security.AuthContext{
		UserID:         userID,
		TenantID:       tenantID,
		RoleCode:       model.RoleCodeSchoolAdmin,
		TokenExpiredAt: time.Now().UTC().Add(time.Hour),
	}
}

func makePaymentPolicy() *model.PaymentPolicy {
	return &model.PaymentPolicy{
		TenantID:            tenantID,
		Code:                "PARTIAL_PAYMENT",
		Name:                "Partial Payment",
		Description:         "Can be paid gradually",
		AllowPartial:        true,
		MinimumAmount:       100000,
		MinimumPercentage:   20,
		AllowOverPayment:    false,
		AutoCloseObligation: true,
		BaseModel: model.BaseModel{
			ID: policyID,
		},
	}
}

func makePaymentProduct() *model.PaymentProduct {
	return &model.PaymentProduct{
		TenantID:           tenantID,
		PaymentPolicyID:    policyID,
		Code:               "SMAN1_SPP_JUL_2026",
		Name:               "SPP July",
		Description:        "Monthly tuition",
		RevenueAccountCode: "4101",
		RevenueAccountName: "Tuition Revenue",
		Price:              500000,
		Status:             model.PaymentProductStatusActive,
		BaseModel: model.BaseModel{
			ID: productID,
		},
	}
}

func assertAmountEqual(t *testing.T, expected float64, actual float64) {
	t.Helper()
	if math.Abs(expected-actual) > 0.0001 {
		t.Fatalf("expected amount %.4f, got %.4f", expected, actual)
	}
}

func boolPtr(value bool) *bool {
	return &value
}

func stringPtr(value string) *string {
	return &value
}

type fakePaymentRepositoryManager struct {
	paymentPolicyRepo  domain.PaymentPolicyRepository
	paymentProductRepo domain.PaymentProductRepository
}

func (f *fakePaymentRepositoryManager) WithDB(db *gorm.DB) domain.RepositoryManagerPayment {
	panic("unexpected call: WithDB")
}

func (f *fakePaymentRepositoryManager) Transaction(
	ctx context.Context,
	fn func(repo domain.RepositoryManagerPayment) error,
) error {
	return fn(f)
}

func (f *fakePaymentRepositoryManager) LedgerEntry() domain.LedgerEntryRepository {
	panic("unexpected call: LedgerEntry")
}

func (f *fakePaymentRepositoryManager) PaymentPolicy() domain.PaymentPolicyRepository {
	if f.paymentPolicyRepo == nil {
		panic("unexpected call: PaymentPolicy")
	}
	return f.paymentPolicyRepo
}

func (f *fakePaymentRepositoryManager) PaymentProduct() domain.PaymentProductRepository {
	if f.paymentProductRepo == nil {
		panic("unexpected call: PaymentProduct")
	}
	return f.paymentProductRepo
}

func (f *fakePaymentRepositoryManager) StudentObligation() domain.StudentObligationRepository {
	panic("unexpected call: StudentObligation")
}

func (f *fakePaymentRepositoryManager) PaymentAllocation() domain.PaymentAllocationRepository {
	panic("unexpected call: PaymentAllocation")
}

func (f *fakePaymentRepositoryManager) PaymentOrder() domain.PaymentOrderRepository {
	panic("unexpected call: PaymentOrder")
}
