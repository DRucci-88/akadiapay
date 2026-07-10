package phase4_test

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
	tenantID       = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	adminUserID    = uuid.MustParse("22222222-2222-2222-2222-222222222222")
	studentID      = uuid.MustParse("33333333-3333-3333-3333-333333333333")
	otherStudentID = uuid.MustParse("99999999-9999-9999-9999-999999999999")
	orderID        = uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	obligationAID  = uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	obligationBID  = uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")
	policyAID      = uuid.MustParse("dddddddd-dddd-dddd-dddd-dddddddddddd")
	policyBID      = uuid.MustParse("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee")
	productAID     = uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff")
	productBID     = uuid.MustParse("12121212-1212-1212-1212-121212121212")
)

func adminAuthContext() *security.AuthContext {
	return &security.AuthContext{
		UserID:         adminUserID,
		TenantID:       tenantID,
		RoleCode:       model.RoleCodeSchoolAdmin,
		TokenExpiredAt: time.Now().UTC().Add(time.Hour),
	}
}

func makePaymentPolicy(id uuid.UUID) *model.PaymentPolicy {
	return &model.PaymentPolicy{
		TenantID:            tenantID,
		Code:                "PARTIAL",
		Name:                "Partial Policy",
		AllowPartial:        true,
		MinimumAmount:       0,
		MinimumPercentage:   0,
		AllowOverPayment:    false,
		AutoCloseObligation: true,
		BaseModel: model.BaseModel{
			ID: id,
		},
	}
}

func makePaymentProduct(id uuid.UUID, policy *model.PaymentPolicy, code string, name string) *model.PaymentProduct {
	return &model.PaymentProduct{
		TenantID:           tenantID,
		PaymentPolicyID:    policy.ID,
		PaymentPolicy:      policy,
		Code:               code,
		Name:               name,
		RevenueAccountCode: "4101",
		RevenueAccountName: "Tuition Revenue",
		Price:              500000,
		Status:             model.PaymentProductStatusActive,
		BaseModel: model.BaseModel{
			ID: id,
		},
	}
}

func makePaymentOrder(totalAmount float64, method model.PaymentOrderPaymentMethod) *model.PaymentOrder {
	return &model.PaymentOrder{
		TenantID:      tenantID,
		StudentID:     studentID,
		PaidByUserID:  adminUserID,
		OrderNumber:   "PO-20260710220000-1a2b3c4d",
		OrderDate:     time.Date(2026, time.July, 10, 0, 0, 0, 0, time.UTC),
		TotalAmount:   totalAmount,
		Status:        model.PaymentOrderStatusPending,
		PaymentMethod: method,
		Notes:         "July payment",
		BaseModel: model.BaseModel{
			ID: orderID,
		},
	}
}

func makeStudentObligation(id uuid.UUID, student uuid.UUID, productID uuid.UUID, outstanding float64) model.StudentObligation {
	return model.StudentObligation{
		StudentID:         student,
		PaymentProductID:  productID,
		Period:            time.Date(2026, time.July, 1, 0, 0, 0, 0, time.UTC),
		OriginalAmount:    outstanding,
		OutstandingAmount: outstanding,
		DueDate:           time.Date(2026, time.July, 10, 0, 0, 0, 0, time.UTC),
		IssuedAt:          time.Date(2026, time.July, 1, 8, 0, 0, 0, time.UTC),
		Status:            model.StudentObligationStatusPending,
		BaseModel: model.BaseModel{
			ID: id,
		},
	}
}

func makePaymentAllocation(obligationID uuid.UUID, amount float64) model.PaymentAllocation {
	return model.PaymentAllocation{
		PaymentOrderID:      orderID,
		StudentObligationID: obligationID,
		AllocatedAmount:     amount,
	}
}

func assertAmountEqual(t *testing.T, expected float64, actual float64) {
	t.Helper()
	if math.Abs(expected-actual) > 0.0001 {
		t.Fatalf("expected amount %.4f, got %.4f", expected, actual)
	}
}

type fakePaymentRepositoryManager struct {
	ledgerEntryRepo       domain.LedgerEntryRepository
	paymentProductRepo    domain.PaymentProductRepository
	studentObligationRepo domain.StudentObligationRepository
	paymentAllocationRepo domain.PaymentAllocationRepository
	paymentOrderRepo      domain.PaymentOrderRepository
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
	if f.ledgerEntryRepo == nil {
		panic("unexpected call: LedgerEntry")
	}
	return f.ledgerEntryRepo
}

func (f *fakePaymentRepositoryManager) PaymentPolicy() domain.PaymentPolicyRepository {
	panic("unexpected call: PaymentPolicy")
}

func (f *fakePaymentRepositoryManager) PaymentProduct() domain.PaymentProductRepository {
	if f.paymentProductRepo == nil {
		panic("unexpected call: PaymentProduct")
	}
	return f.paymentProductRepo
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
