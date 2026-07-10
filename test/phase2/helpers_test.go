package phase2_test

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
	tenantID        = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	otherTenantID   = uuid.MustParse("99999999-9999-9999-9999-999999999999")
	adminUserID     = uuid.MustParse("22222222-2222-2222-2222-222222222222")
	parentUserID    = uuid.MustParse("77777777-7777-7777-7777-777777777777")
	studentUserID   = uuid.MustParse("88888888-8888-8888-8888-888888888888")
	studentID       = uuid.MustParse("33333333-3333-3333-3333-333333333333")
	secondStudentID = uuid.MustParse("66666666-6666-6666-6666-666666666666")
	policyID        = uuid.MustParse("44444444-4444-4444-4444-444444444444")
	productID       = uuid.MustParse("55555555-5555-5555-5555-555555555555")
	obligationID    = uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
)

func adminAuthContext() *security.AuthContext {
	return &security.AuthContext{
		UserID:         adminUserID,
		TenantID:       tenantID,
		RoleCode:       model.RoleCodeSchoolAdmin,
		TokenExpiredAt: time.Now().UTC().Add(time.Hour),
	}
}

func parentAuthContext() *security.AuthContext {
	return &security.AuthContext{
		UserID:         parentUserID,
		TenantID:       tenantID,
		RoleCode:       model.RoleCodeParent,
		TokenExpiredAt: time.Now().UTC().Add(time.Hour),
	}
}

func studentAuthContext(student uuid.UUID) *security.AuthContext {
	return &security.AuthContext{
		UserID:         studentUserID,
		TenantID:       tenantID,
		StudentID:      &student,
		RoleCode:       model.RoleCodeStudent,
		TokenExpiredAt: time.Now().UTC().Add(time.Hour),
	}
}

func makeStudent(id uuid.UUID, tenant uuid.UUID) *model.Student {
	return &model.Student{
		TenantID: tenant,
		NISN:     "20260001",
		FullName: "Student Name",
		Status:   model.StudentStatusActive,
		BaseModel: model.BaseModel{
			ID: id,
		},
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
		PaymentPolicy:      makePaymentPolicy(),
		Code:               "SPP_JUL_2026",
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

func makeStudentObligation(student uuid.UUID) *model.StudentObligation {
	return &model.StudentObligation{
		StudentID:         student,
		PaymentProductID:  productID,
		Period:            time.Date(2026, time.July, 1, 0, 0, 0, 0, time.UTC),
		OriginalAmount:    500000,
		OutstandingAmount: 500000,
		DueDate:           time.Date(2026, time.July, 10, 0, 0, 0, 0, time.UTC),
		IssuedAt:          time.Date(2026, time.July, 1, 8, 0, 0, 0, time.UTC),
		Status:            model.StudentObligationStatusPending,
		Notes:             "July tuition",
		BaseModel: model.BaseModel{
			ID: obligationID,
		},
	}
}

func assertAmountEqual(t *testing.T, expected float64, actual float64) {
	t.Helper()
	if math.Abs(expected-actual) > 0.0001 {
		t.Fatalf("expected amount %.4f, got %.4f", expected, actual)
	}
}

func floatPtr(value float64) *float64 {
	return &value
}

func stringPtr(value string) *string {
	return &value
}

func timePtr(value time.Time) *time.Time {
	return &value
}

type fakePaymentRepositoryManager struct {
	studentObligationRepo domain.StudentObligationRepository
	transactionFn         func(ctx context.Context, fn func(repo domain.RepositoryManagerPayment) error) error
}

func (f *fakePaymentRepositoryManager) WithDB(db *gorm.DB) domain.RepositoryManagerPayment {
	panic("unexpected call: WithDB")
}

func (f *fakePaymentRepositoryManager) Transaction(
	ctx context.Context,
	fn func(repo domain.RepositoryManagerPayment) error,
) error {
	if f.transactionFn != nil {
		return f.transactionFn(ctx, fn)
	}
	return fn(f)
}

func (f *fakePaymentRepositoryManager) LedgerEntry() domain.LedgerEntryRepository {
	panic("unexpected call: LedgerEntry")
}

func (f *fakePaymentRepositoryManager) PaymentPolicy() domain.PaymentPolicyRepository {
	panic("unexpected call: PaymentPolicy")
}

func (f *fakePaymentRepositoryManager) PaymentProduct() domain.PaymentProductRepository {
	panic("unexpected call: PaymentProduct")
}

func (f *fakePaymentRepositoryManager) StudentObligation() domain.StudentObligationRepository {
	if f.studentObligationRepo == nil {
		panic("unexpected call: StudentObligation")
	}
	return f.studentObligationRepo
}

func (f *fakePaymentRepositoryManager) PaymentAllocation() domain.PaymentAllocationRepository {
	panic("unexpected call: PaymentAllocation")
}

func (f *fakePaymentRepositoryManager) PaymentOrder() domain.PaymentOrderRepository {
	panic("unexpected call: PaymentOrder")
}
