package phase3_test

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
	orderID         = uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	allocationID    = uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
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

func makePaymentOrder(student uuid.UUID) *model.PaymentOrder {
	return &model.PaymentOrder{
		TenantID:      tenantID,
		StudentID:     student,
		PaidByUserID:  adminUserID,
		OrderNumber:   "PO-20260710220000-1a2b3c4d",
		OrderDate:     time.Date(2026, time.July, 10, 0, 0, 0, 0, time.UTC),
		TotalAmount:   500000,
		Status:        model.PaymentOrderStatusPending,
		PaymentMethod: model.PaymentMethodCash,
		Notes:         "July payment",
		BaseModel: model.BaseModel{
			ID: orderID,
		},
	}
}

func makePaymentAllocation() *model.PaymentAllocation {
	return &model.PaymentAllocation{
		PaymentOrderID:      orderID,
		StudentObligationID: uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc"),
		AllocatedAmount:     250000,
		BaseModel: model.BaseModel{
			ID: allocationID,
		},
	}
}

func assertAmountEqual(t *testing.T, expected float64, actual float64) {
	t.Helper()
	if math.Abs(expected-actual) > 0.0001 {
		t.Fatalf("expected amount %.4f, got %.4f", expected, actual)
	}
}

type fakePaymentRepositoryManager struct {
	studentObligationRepo domain.StudentObligationRepository
	paymentAllocationRepo domain.PaymentAllocationRepository
	paymentOrderRepo      domain.PaymentOrderRepository
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
	if f.paymentAllocationRepo == nil {
		panic("unexpected call: PaymentAllocation")
	}
	return f.paymentAllocationRepo
}

func (f *fakePaymentRepositoryManager) PaymentOrder() domain.PaymentOrderRepository {
	if f.paymentOrderRepo == nil {
		panic("unexpected call: PaymentOrder")
	}
	return f.paymentOrderRepo
}
