package phase5_test

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
	tenantID      = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	adminUserID   = uuid.MustParse("22222222-2222-2222-2222-222222222222")
	treasurerID   = uuid.MustParse("55555555-5555-5555-5555-555555555555")
	parentUserID  = uuid.MustParse("77777777-7777-7777-7777-777777777777")
	studentID     = uuid.MustParse("33333333-3333-3333-3333-333333333333")
	orderID       = uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	obligationAID = uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	obligationBID = uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")
	policyAID     = uuid.MustParse("dddddddd-dddd-dddd-dddd-dddddddddddd")
	policyBID     = uuid.MustParse("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee")
	productAID    = uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff")
	productBID    = uuid.MustParse("12121212-1212-1212-1212-121212121212")
	entryAID      = uuid.MustParse("23232323-2323-2323-2323-232323232323")
	entryBID      = uuid.MustParse("34343434-3434-3434-3434-343434343434")
)

func adminAuthContext() *security.AuthContext {
	return &security.AuthContext{
		UserID:         adminUserID,
		TenantID:       tenantID,
		RoleCode:       model.RoleCodeSchoolAdmin,
		TokenExpiredAt: time.Now().UTC().Add(time.Hour),
	}
}

func treasurerAuthContext() *security.AuthContext {
	return &security.AuthContext{
		UserID:         treasurerID,
		TenantID:       tenantID,
		RoleCode:       model.RoleCodeTreasurer,
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

func makePaymentPolicy(id uuid.UUID) *model.PaymentPolicy {
	return &model.PaymentPolicy{
		TenantID:            tenantID,
		Code:                "PARTIAL",
		Name:                "Partial Policy",
		AllowPartial:        true,
		AllowOverPayment:    false,
		AutoCloseObligation: true,
		BaseModel: model.BaseModel{
			ID: id,
		},
	}
}

func makePaymentProduct(id uuid.UUID, policy *model.PaymentPolicy, code string, name string) *model.PaymentProduct {
	return &model.PaymentProduct{
		TenantID:        tenantID,
		PaymentPolicyID: policy.ID,
		PaymentPolicy:   policy,
		Code:            code,
		Name:            name,
		Status:          model.PaymentProductStatusActive,
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
		OrderDate:     time.Date(2026, time.July, 10, 15, 30, 0, 0, time.UTC),
		TotalAmount:   totalAmount,
		Status:        model.PaymentOrderStatusCompleted,
		PaymentMethod: method,
		Notes:         "July payment",
		BaseModel: model.BaseModel{
			ID: orderID,
		},
	}
}

func makeStudentObligation(id uuid.UUID, productID uuid.UUID) model.StudentObligation {
	return model.StudentObligation{
		StudentID:         studentID,
		PaymentProductID:  productID,
		Period:            time.Date(2026, time.July, 1, 0, 0, 0, 0, time.UTC),
		OriginalAmount:    500000,
		OutstandingAmount: 0,
		DueDate:           time.Date(2026, time.July, 10, 0, 0, 0, 0, time.UTC),
		IssuedAt:          time.Date(2026, time.July, 1, 8, 0, 0, 0, time.UTC),
		Status:            model.StudentObligationStatusClosed,
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

func makeLedgerEntry(id uuid.UUID, debit float64, credit float64, accountCode string, accountName string) model.LedgerEntry {
	return model.LedgerEntry{
		PaymentOrderID: orderID,
		EntryDate:      time.Date(2026, time.July, 10, 0, 0, 0, 0, time.UTC),
		AccountCode:    accountCode,
		AccountName:    accountName,
		Debit:          debit,
		Credit:         credit,
		Description:    "Posting",
		BaseModel: model.BaseModel{
			ID: id,
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
