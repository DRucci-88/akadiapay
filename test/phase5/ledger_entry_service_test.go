package phase5_test

import (
	"akadia/domain"
	service "akadia/internal/payment/service"
	"akadia/internal/platform/security"
	"akadia/internal/shared"
	"akadia/model"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type fakePaymentOrderService struct {
	findByIDFn func(ctx context.Context, authContext *security.AuthContext, id uuid.UUID) (*model.PaymentOrder, error)
}

func (f *fakePaymentOrderService) FindPaginate(
	ctx context.Context,
	pageable *shared.Pageable,
	filter *domain.PaymentOrderFilter,
	authContext *security.AuthContext,
) (*shared.Page[domain.PaymentOrderResponse], error) {
	panic("unexpected call: FindPaginate")
}

func (f *fakePaymentOrderService) FindByID(
	ctx context.Context,
	authContext *security.AuthContext,
	id uuid.UUID,
) (*model.PaymentOrder, error) {
	if f.findByIDFn == nil {
		panic("unexpected call: FindByID")
	}
	return f.findByIDFn(ctx, authContext, id)
}

func (f *fakePaymentOrderService) Create(
	ctx context.Context,
	authContext *security.AuthContext,
	req *domain.PaymentOrderCreate,
) (*model.PaymentOrder, error) {
	panic("unexpected call: Create")
}

func (f *fakePaymentOrderService) Cancel(
	ctx context.Context,
	authContext *security.AuthContext,
	id uuid.UUID,
) (*model.PaymentOrder, error) {
	panic("unexpected call: Cancel")
}

type fakePaymentOrderRepository struct {
	lockByIDFn         func(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (*model.PaymentOrder, error)
	updateStatusFn     func(ctx context.Context, id uuid.UUID, status model.PaymentOrderStatus) (int, error)
	markLedgerPostedFn func(ctx context.Context, id uuid.UUID, ledgerPostedAt *time.Time) (int, error)
}

func (f *fakePaymentOrderRepository) Paginate(
	ctx context.Context,
	pageable *shared.Pageable,
	chain gorm.ChainInterface[model.PaymentOrder],
) (*shared.Page[model.PaymentOrder], error) {
	panic("unexpected call: Paginate")
}

func (f *fakePaymentOrderRepository) FindPaginate(
	ctx context.Context,
	pageable *shared.Pageable,
	filter *domain.PaymentOrderFilter,
	authContext *security.AuthContext,
) (*shared.Page[model.PaymentOrder], error) {
	panic("unexpected call: FindPaginate")
}

func (f *fakePaymentOrderRepository) Create(
	ctx context.Context,
	paymentOrder *model.PaymentOrder,
) error {
	panic("unexpected call: Create")
}

func (f *fakePaymentOrderRepository) FirstByID(
	ctx context.Context,
	id uuid.UUID,
	tenantID uuid.UUID,
) (*model.PaymentOrder, error) {
	panic("unexpected call: FirstByID")
}

func (f *fakePaymentOrderRepository) LockByID(
	ctx context.Context,
	id uuid.UUID,
	tenantID uuid.UUID,
) (*model.PaymentOrder, error) {
	if f.lockByIDFn == nil {
		panic("unexpected call: LockByID")
	}
	return f.lockByIDFn(ctx, id, tenantID)
}

func (f *fakePaymentOrderRepository) UpdateStatus(
	ctx context.Context,
	id uuid.UUID,
	status model.PaymentOrderStatus,
) (int, error) {
	if f.updateStatusFn == nil {
		panic("unexpected call: UpdateStatus")
	}
	return f.updateStatusFn(ctx, id, status)
}

func (f *fakePaymentOrderRepository) MarkLedgerPosted(
	ctx context.Context,
	id uuid.UUID,
	ledgerPostedAt *time.Time,
) (int, error) {
	if f.markLedgerPostedFn == nil {
		panic("unexpected call: MarkLedgerPosted")
	}
	return f.markLedgerPostedFn(ctx, id, ledgerPostedAt)
}

type fakePaymentAllocationRepository struct {
	findByPaymentOrderIDFn func(ctx context.Context, paymentOrderID uuid.UUID) ([]model.PaymentAllocation, error)
}

func (f *fakePaymentAllocationRepository) CreateBatch(
	ctx context.Context,
	paymentAllocations []model.PaymentAllocation,
) error {
	panic("unexpected call: CreateBatch")
}

func (f *fakePaymentAllocationRepository) FindByPaymentOrderID(
	ctx context.Context,
	paymentOrderID uuid.UUID,
) ([]model.PaymentAllocation, error) {
	if f.findByPaymentOrderIDFn == nil {
		panic("unexpected call: FindByPaymentOrderID")
	}
	return f.findByPaymentOrderIDFn(ctx, paymentOrderID)
}

type fakeStudentObligationRepository struct {
	findByIDsFn func(ctx context.Context, ids []uuid.UUID) ([]model.StudentObligation, error)
}

func (f *fakeStudentObligationRepository) Create(
	ctx context.Context,
	studentObligation *model.StudentObligation,
) error {
	panic("unexpected call: Create")
}

func (f *fakeStudentObligationRepository) CreateBatch(
	ctx context.Context,
	studentObligations []model.StudentObligation,
) error {
	panic("unexpected call: CreateBatch")
}

func (f *fakeStudentObligationRepository) Paginate(
	ctx context.Context,
	pageable *shared.Pageable,
	chain gorm.ChainInterface[model.StudentObligation],
) (*shared.Page[model.StudentObligation], error) {
	panic("unexpected call: Paginate")
}

func (f *fakeStudentObligationRepository) FindPaginate(
	ctx context.Context,
	pageable *shared.Pageable,
	filter *domain.StudentObligationFilter,
	authContext *security.AuthContext,
) (*shared.Page[model.StudentObligation], error) {
	panic("unexpected call: FindPaginate")
}

func (f *fakeStudentObligationRepository) SumOutstandingByStudentID(
	ctx context.Context,
	studentID uuid.UUID,
) (float64, error) {
	panic("unexpected call: SumOutstandingByStudentID")
}

func (f *fakeStudentObligationRepository) FindOutstandingByStudentID(
	ctx context.Context,
	studentID uuid.UUID,
) ([]model.StudentObligation, error) {
	panic("unexpected call: FindOutstandingByStudentID")
}

func (f *fakeStudentObligationRepository) FirstByID(
	ctx context.Context,
	id uuid.UUID,
	tenantID uuid.UUID,
) (*model.StudentObligation, error) {
	panic("unexpected call: FirstByID")
}

func (f *fakeStudentObligationRepository) LockByIDs(
	ctx context.Context,
	ids []uuid.UUID,
) ([]model.StudentObligation, error) {
	panic("unexpected call: LockByIDs")
}

func (f *fakeStudentObligationRepository) Update(
	ctx context.Context,
	id uuid.UUID,
	req *domain.StudentObligationUpdate,
) (int, error) {
	panic("unexpected call: Update")
}

func (f *fakeStudentObligationRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
) (int, error) {
	panic("unexpected call: Delete")
}

func (f *fakeStudentObligationRepository) HasPaymentAllocations(
	ctx context.Context,
	id uuid.UUID,
) (bool, error) {
	panic("unexpected call: HasPaymentAllocations")
}

func (f *fakeStudentObligationRepository) FindByIDs(
	ctx context.Context,
	ids []uuid.UUID,
) ([]model.StudentObligation, error) {
	if f.findByIDsFn == nil {
		panic("unexpected call: FindByIDs")
	}
	return f.findByIDsFn(ctx, ids)
}

func (f *fakeStudentObligationRepository) UpdateSettlement(
	ctx context.Context,
	id uuid.UUID,
	outstandingAmount float64,
	status model.StudentObligationStatus,
) (int, error) {
	panic("unexpected call: UpdateSettlement")
}

func (f *fakeStudentObligationRepository) ExistsActiveByStudentIDAndPaymentProductIDAndPeriod(
	ctx context.Context,
	studentID uuid.UUID,
	paymentProductID uuid.UUID,
	period time.Time,
) (bool, error) {
	panic("unexpected call: ExistsActiveByStudentIDAndPaymentProductIDAndPeriod")
}

type fakePaymentProductRepository struct {
	findByIDsIncludingDeletedFn func(ctx context.Context, ids []uuid.UUID) ([]model.PaymentProduct, error)
}

func (f *fakePaymentProductRepository) QueryWithPreloads(
	preloads ...model.PaymentProductPreload,
) gorm.ChainInterface[model.PaymentProduct] {
	panic("unexpected call: QueryWithPreloads")
}

func (f *fakePaymentProductRepository) Paginate(
	ctx context.Context,
	pageable *shared.Pageable,
	chain gorm.ChainInterface[model.PaymentProduct],
) (*shared.Page[model.PaymentProduct], error) {
	panic("unexpected call: Paginate")
}

func (f *fakePaymentProductRepository) FindPaginate(
	ctx context.Context,
	pageable *shared.Pageable,
	filter *domain.PaymentProductFilter,
	authContext *security.AuthContext,
) (*shared.Page[model.PaymentProduct], error) {
	panic("unexpected call: FindPaginate")
}

func (f *fakePaymentProductRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
	tenantID uuid.UUID,
	preloads ...model.PaymentProductPreload,
) (*model.PaymentProduct, error) {
	panic("unexpected call: FindByID")
}

func (f *fakePaymentProductRepository) FindByIDsIncludingDeleted(
	ctx context.Context,
	ids []uuid.UUID,
) ([]model.PaymentProduct, error) {
	if f.findByIDsIncludingDeletedFn == nil {
		panic("unexpected call: FindByIDsIncludingDeleted")
	}
	return f.findByIDsIncludingDeletedFn(ctx, ids)
}

func (f *fakePaymentProductRepository) Create(
	ctx context.Context,
	paymentProduct *model.PaymentProduct,
) error {
	panic("unexpected call: Create")
}

func (f *fakePaymentProductRepository) Update(
	ctx context.Context,
	id uuid.UUID,
	tenantID uuid.UUID,
	req *domain.PaymentProductUpdate,
) (int, error) {
	panic("unexpected call: Update")
}

type fakeLedgerEntryRepository struct {
	findByPaymentOrderIDFn   func(ctx context.Context, paymentOrderID uuid.UUID, tenantID uuid.UUID) ([]model.LedgerEntry, error)
	existsByPaymentOrderIDFn func(ctx context.Context, paymentOrderID uuid.UUID, tenantID uuid.UUID) (bool, error)
	createInBatchesFn        func(ctx context.Context, entries []model.LedgerEntry, batchSize int) error
	findPaginateFn           func(ctx context.Context, tenantID uuid.UUID, pageable *shared.Pageable, filter *domain.LedgerEntryFilter) (*shared.Page[model.LedgerEntry], error)
}

func (f *fakeLedgerEntryRepository) FindByPaymentOrderID(
	ctx context.Context,
	paymentOrderID uuid.UUID,
	tenantID uuid.UUID,
) ([]model.LedgerEntry, error) {
	if f.findByPaymentOrderIDFn == nil {
		panic("unexpected call: FindByPaymentOrderID")
	}
	return f.findByPaymentOrderIDFn(ctx, paymentOrderID, tenantID)
}

func (f *fakeLedgerEntryRepository) ExistsByPaymentOrderID(
	ctx context.Context,
	paymentOrderID uuid.UUID,
	tenantID uuid.UUID,
) (bool, error) {
	if f.existsByPaymentOrderIDFn == nil {
		panic("unexpected call: ExistsByPaymentOrderID")
	}
	return f.existsByPaymentOrderIDFn(ctx, paymentOrderID, tenantID)
}

func (f *fakeLedgerEntryRepository) CreateInBatches(
	ctx context.Context,
	entries []model.LedgerEntry,
	batchSize int,
) error {
	if f.createInBatchesFn == nil {
		panic("unexpected call: CreateInBatches")
	}
	return f.createInBatchesFn(ctx, entries, batchSize)
}

func (f *fakeLedgerEntryRepository) FindPaginate(
	ctx context.Context,
	tenantID uuid.UUID,
	pageable *shared.Pageable,
	filter *domain.LedgerEntryFilter,
) (*shared.Page[model.LedgerEntry], error) {
	if f.findPaginateFn == nil {
		panic("unexpected call: FindPaginate")
	}
	return f.findPaginateFn(ctx, tenantID, pageable, filter)
}

func TestLedgerEntryServicePostPaymentRejectsUnauthorizedRole(t *testing.T) {
	transactionCalls := 0
	var manager *fakePaymentRepositoryManager
	manager = &fakePaymentRepositoryManager{
		ledgerEntryRepo:       &fakeLedgerEntryRepository{},
		paymentProductRepo:    &fakePaymentProductRepository{},
		studentObligationRepo: &fakeStudentObligationRepository{},
		paymentAllocationRepo: &fakePaymentAllocationRepository{},
		paymentOrderRepo:      &fakePaymentOrderRepository{},
		transactionFn: func(ctx context.Context, fn func(repo domain.RepositoryManagerPayment) error) error {
			transactionCalls++
			return fn(manager)
		},
	}
	svc := service.NewLedgerEntryService(
		manager,
		&fakePaymentOrderService{
			findByIDFn: func(ctx context.Context, authContext *security.AuthContext, id uuid.UUID) (*model.PaymentOrder, error) {
				t.Fatalf("payment order lookup should not run for unauthorized role")
				return nil, nil
			},
		},
	)

	err := svc.PostPayment(context.Background(), parentAuthContext(), orderID)

	if !errors.Is(err, shared.ErrAuthUnauthorized) {
		t.Fatalf("expected error %v, got %v", shared.ErrAuthUnauthorized, err)
	}
	if transactionCalls != 0 {
		t.Fatalf("expected transaction not to start for unauthorized role")
	}
}

func TestLedgerEntryServicePostPaymentCreatesBalancedEntries(t *testing.T) {
	currentOrder := makePaymentOrder(300000, model.PaymentMethodVirtualAccount)
	policyA := makePaymentPolicy(policyAID)
	policyB := makePaymentPolicy(policyBID)
	productA := makePaymentProduct(productAID, policyA, "SPP-2026", "Tuition July")
	productA.RevenueAccountCode = "4101"
	productA.RevenueAccountName = "Tuition Revenue"
	productB := makePaymentProduct(productBID, policyB, "REGISTRATION-2026", "Registration Fee")
	productB.RevenueAccountCode = ""
	productB.RevenueAccountName = ""
	allocations := []model.PaymentAllocation{
		makePaymentAllocation(obligationAID, 200000),
		makePaymentAllocation(obligationBID, 100000),
	}
	obligations := []model.StudentObligation{
		makeStudentObligation(obligationAID, productAID),
		makeStudentObligation(obligationBID, productBID),
	}

	transactionCalls := 0
	markPostedCalls := 0
	var createdEntries []model.LedgerEntry
	var manager *fakePaymentRepositoryManager
	manager = &fakePaymentRepositoryManager{
		ledgerEntryRepo: &fakeLedgerEntryRepository{
			existsByPaymentOrderIDFn: func(ctx context.Context, paymentOrderID uuid.UUID, tenantID uuid.UUID) (bool, error) {
				return false, nil
			},
			createInBatchesFn: func(ctx context.Context, entries []model.LedgerEntry, batchSize int) error {
				createdEntries = append([]model.LedgerEntry(nil), entries...)
				if batchSize != 100 {
					t.Fatalf("expected batch size 100, got %d", batchSize)
				}
				return nil
			},
		},
		paymentProductRepo: &fakePaymentProductRepository{
			findByIDsIncludingDeletedFn: func(ctx context.Context, ids []uuid.UUID) ([]model.PaymentProduct, error) {
				return []model.PaymentProduct{*productA, *productB}, nil
			},
		},
		studentObligationRepo: &fakeStudentObligationRepository{
			findByIDsFn: func(ctx context.Context, ids []uuid.UUID) ([]model.StudentObligation, error) {
				return obligations, nil
			},
		},
		paymentAllocationRepo: &fakePaymentAllocationRepository{
			findByPaymentOrderIDFn: func(ctx context.Context, paymentOrderID uuid.UUID) ([]model.PaymentAllocation, error) {
				return allocations, nil
			},
		},
		paymentOrderRepo: &fakePaymentOrderRepository{
			lockByIDFn: func(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (*model.PaymentOrder, error) {
				return currentOrder, nil
			},
			updateStatusFn: func(ctx context.Context, id uuid.UUID, status model.PaymentOrderStatus) (int, error) {
				t.Fatalf("ledger posting should not update order status")
				return 0, nil
			},
			markLedgerPostedFn: func(ctx context.Context, id uuid.UUID, ledgerPostedAt *time.Time) (int, error) {
				markPostedCalls++
				currentOrder.LedgerPostedAt = ledgerPostedAt
				return 1, nil
			},
		},
		transactionFn: func(ctx context.Context, fn func(repo domain.RepositoryManagerPayment) error) error {
			transactionCalls++
			return fn(manager)
		},
	}
	svc := service.NewLedgerEntryService(
		manager,
		&fakePaymentOrderService{
			findByIDFn: func(ctx context.Context, authContext *security.AuthContext, id uuid.UUID) (*model.PaymentOrder, error) {
				return currentOrder, nil
			},
		},
	)

	err := svc.PostPayment(context.Background(), treasurerAuthContext(), orderID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if transactionCalls != 1 {
		t.Fatalf("expected one transaction, got %d", transactionCalls)
	}
	if len(createdEntries) != 3 {
		t.Fatalf("expected one debit and two credit entries, got %d", len(createdEntries))
	}
	assertAmountEqual(t, 300000, createdEntries[0].Debit)
	if createdEntries[0].AccountCode != "1104" || createdEntries[0].AccountName != "Virtual Account Clearing" {
		t.Fatalf("expected virtual account debit mapping, got %+v", createdEntries[0])
	}
	assertAmountEqual(t, 200000, createdEntries[1].Credit)
	if createdEntries[1].AccountCode != "4101" || createdEntries[1].AccountName != "Tuition Revenue" {
		t.Fatalf("expected first credit entry for tuition revenue, got %+v", createdEntries[1])
	}
	assertAmountEqual(t, 100000, createdEntries[2].Credit)
	if createdEntries[2].AccountCode != "4102" || createdEntries[2].AccountName != "Registration Revenue" {
		t.Fatalf("expected fallback registration revenue mapping, got %+v", createdEntries[2])
	}
	if markPostedCalls != 1 || currentOrder.LedgerPostedAt == nil {
		t.Fatalf("expected ledger posted timestamp to be marked once")
	}
}

func TestLedgerEntryServicePostPaymentIsIdempotentWhenEntriesAlreadyExist(t *testing.T) {
	currentOrder := makePaymentOrder(300000, model.PaymentMethodCash)
	createCalls := 0
	markPostedCalls := 0
	var manager *fakePaymentRepositoryManager
	manager = &fakePaymentRepositoryManager{
		ledgerEntryRepo: &fakeLedgerEntryRepository{
			existsByPaymentOrderIDFn: func(ctx context.Context, paymentOrderID uuid.UUID, tenantID uuid.UUID) (bool, error) {
				return true, nil
			},
			createInBatchesFn: func(ctx context.Context, entries []model.LedgerEntry, batchSize int) error {
				createCalls++
				return nil
			},
		},
		paymentProductRepo:    &fakePaymentProductRepository{},
		studentObligationRepo: &fakeStudentObligationRepository{},
		paymentAllocationRepo: &fakePaymentAllocationRepository{
			findByPaymentOrderIDFn: func(ctx context.Context, paymentOrderID uuid.UUID) ([]model.PaymentAllocation, error) {
				t.Fatalf("allocations should not be loaded when ledger entries already exist")
				return nil, nil
			},
		},
		paymentOrderRepo: &fakePaymentOrderRepository{
			lockByIDFn: func(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (*model.PaymentOrder, error) {
				return currentOrder, nil
			},
			updateStatusFn: func(ctx context.Context, id uuid.UUID, status model.PaymentOrderStatus) (int, error) {
				t.Fatalf("order status should not be updated during repost")
				return 0, nil
			},
			markLedgerPostedFn: func(ctx context.Context, id uuid.UUID, ledgerPostedAt *time.Time) (int, error) {
				markPostedCalls++
				currentOrder.LedgerPostedAt = ledgerPostedAt
				return 1, nil
			},
		},
		transactionFn: func(ctx context.Context, fn func(repo domain.RepositoryManagerPayment) error) error {
			return fn(manager)
		},
	}
	svc := service.NewLedgerEntryService(
		manager,
		&fakePaymentOrderService{
			findByIDFn: func(ctx context.Context, authContext *security.AuthContext, id uuid.UUID) (*model.PaymentOrder, error) {
				return currentOrder, nil
			},
		},
	)

	err := svc.PostPayment(context.Background(), adminAuthContext(), orderID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if createCalls != 0 {
		t.Fatalf("expected no new ledger batch creation on idempotent repost")
	}
	if markPostedCalls != 1 || currentOrder.LedgerPostedAt == nil {
		t.Fatalf("expected ledger posted timestamp to be refreshed once")
	}
}

func TestLedgerEntryServicePostPaymentRejectsInvalidSourceState(t *testing.T) {
	tests := []struct {
		name        string
		order       *model.PaymentOrder
		allocations []model.PaymentAllocation
		obligations []model.StudentObligation
		products    []model.PaymentProduct
		expectedErr error
	}{
		{
			name: "payment order must be completed",
			order: func() *model.PaymentOrder {
				order := makePaymentOrder(300000, model.PaymentMethodCash)
				order.Status = model.PaymentOrderStatusPending
				return order
			}(),
			expectedErr: shared.ErrPaymentOrderStatusInvalid,
		},
		{
			name:        "payment allocations must exist",
			order:       makePaymentOrder(300000, model.PaymentMethodCash),
			allocations: nil,
			expectedErr: shared.ErrPaymentAllocationNotFound,
		},
		{
			name:  "allocation amount must be positive",
			order: makePaymentOrder(300000, model.PaymentMethodCash),
			allocations: []model.PaymentAllocation{
				makePaymentAllocation(obligationAID, 0),
			},
			expectedErr: shared.ErrPaymentAllocationAmountInvalid,
		},
		{
			name:  "allocation total must equal payment total",
			order: makePaymentOrder(300000, model.PaymentMethodCash),
			allocations: []model.PaymentAllocation{
				makePaymentAllocation(obligationAID, 200000),
			},
			expectedErr: shared.ErrLedgerUnbalancedSource,
		},
		{
			name:  "all obligations must resolve",
			order: makePaymentOrder(300000, model.PaymentMethodCash),
			allocations: []model.PaymentAllocation{
				makePaymentAllocation(obligationAID, 300000),
			},
			obligations: nil,
			expectedErr: shared.ErrStudentObligationNotFound,
		},
		{
			name:  "all payment products must resolve",
			order: makePaymentOrder(300000, model.PaymentMethodCash),
			allocations: []model.PaymentAllocation{
				makePaymentAllocation(obligationAID, 300000),
			},
			obligations: []model.StudentObligation{
				makeStudentObligation(obligationAID, productAID),
			},
			products:    nil,
			expectedErr: shared.ErrPaymentProductNotFound,
		},
		{
			name: "payment method must map to a debit account",
			order: func() *model.PaymentOrder {
				order := makePaymentOrder(300000, model.PaymentOrderPaymentMethod("CHEQUE"))
				return order
			}(),
			allocations: []model.PaymentAllocation{
				makePaymentAllocation(obligationAID, 300000),
			},
			obligations: []model.StudentObligation{
				makeStudentObligation(obligationAID, productAID),
			},
			products: []model.PaymentProduct{
				*makePaymentProduct(productAID, makePaymentPolicy(policyAID), "SPP", "Tuition"),
			},
			expectedErr: shared.ErrLedgerDebitAccountNotConfigured,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var manager *fakePaymentRepositoryManager
			manager = &fakePaymentRepositoryManager{
				ledgerEntryRepo: &fakeLedgerEntryRepository{
					existsByPaymentOrderIDFn: func(ctx context.Context, paymentOrderID uuid.UUID, tenantID uuid.UUID) (bool, error) {
						return false, nil
					},
					createInBatchesFn: func(ctx context.Context, entries []model.LedgerEntry, batchSize int) error {
						t.Fatalf("ledger entries should not be written on source validation failure")
						return nil
					},
				},
				paymentProductRepo: &fakePaymentProductRepository{
					findByIDsIncludingDeletedFn: func(ctx context.Context, ids []uuid.UUID) ([]model.PaymentProduct, error) {
						return tt.products, nil
					},
				},
				studentObligationRepo: &fakeStudentObligationRepository{
					findByIDsFn: func(ctx context.Context, ids []uuid.UUID) ([]model.StudentObligation, error) {
						return tt.obligations, nil
					},
				},
				paymentAllocationRepo: &fakePaymentAllocationRepository{
					findByPaymentOrderIDFn: func(ctx context.Context, paymentOrderID uuid.UUID) ([]model.PaymentAllocation, error) {
						return tt.allocations, nil
					},
				},
				paymentOrderRepo: &fakePaymentOrderRepository{
					lockByIDFn: func(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (*model.PaymentOrder, error) {
						return tt.order, nil
					},
					updateStatusFn: func(ctx context.Context, id uuid.UUID, status model.PaymentOrderStatus) (int, error) {
						t.Fatalf("ledger posting should not update order status")
						return 0, nil
					},
					markLedgerPostedFn: func(ctx context.Context, id uuid.UUID, ledgerPostedAt *time.Time) (int, error) {
						t.Fatalf("ledger posted timestamp should not be written on validation failure")
						return 0, nil
					},
				},
				transactionFn: func(ctx context.Context, fn func(repo domain.RepositoryManagerPayment) error) error {
					return fn(manager)
				},
			}
			svc := service.NewLedgerEntryService(
				manager,
				&fakePaymentOrderService{
					findByIDFn: func(ctx context.Context, authContext *security.AuthContext, id uuid.UUID) (*model.PaymentOrder, error) {
						return tt.order, nil
					},
				},
			)

			err := svc.PostPayment(context.Background(), adminAuthContext(), orderID)

			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestLedgerEntryServiceFindByPaymentOrderIDRequiresRoleAndMapsResponses(t *testing.T) {
	svc := service.NewLedgerEntryService(
		&fakePaymentRepositoryManager{
			ledgerEntryRepo: &fakeLedgerEntryRepository{
				findByPaymentOrderIDFn: func(ctx context.Context, paymentOrderID uuid.UUID, tenantID uuid.UUID) ([]model.LedgerEntry, error) {
					return []model.LedgerEntry{
						makeLedgerEntry(entryAID, 300000, 0, "1101", "Cash"),
						makeLedgerEntry(entryBID, 0, 300000, "4101", "Tuition Revenue"),
					}, nil
				},
			},
			paymentProductRepo:    &fakePaymentProductRepository{},
			studentObligationRepo: &fakeStudentObligationRepository{},
			paymentAllocationRepo: &fakePaymentAllocationRepository{},
			paymentOrderRepo:      &fakePaymentOrderRepository{},
		},
		&fakePaymentOrderService{
			findByIDFn: func(ctx context.Context, authContext *security.AuthContext, id uuid.UUID) (*model.PaymentOrder, error) {
				return makePaymentOrder(300000, model.PaymentMethodCash), nil
			},
		},
	)

	_, err := svc.FindByPaymentOrderID(context.Background(), parentAuthContext(), orderID)
	if !errors.Is(err, shared.ErrAuthUnauthorized) {
		t.Fatalf("expected unauthorized error, got %v", err)
	}

	res, err := svc.FindByPaymentOrderID(context.Background(), adminAuthContext(), orderID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(res) != 2 {
		t.Fatalf("expected two ledger responses, got %d", len(res))
	}
	assertAmountEqual(t, 300000, res[0].Debit)
	assertAmountEqual(t, 300000, res[1].Credit)
}

func TestLedgerEntryServiceFindPaginateEnforcesRoleAndReturnsEmptyPageOnError(t *testing.T) {
	pageable := &shared.Pageable{Page: 2, Size: 10}
	pageable.Normalize()
	svc := service.NewLedgerEntryService(
		&fakePaymentRepositoryManager{
			ledgerEntryRepo: &fakeLedgerEntryRepository{
				findPaginateFn: func(ctx context.Context, tenantID uuid.UUID, pageable *shared.Pageable, filter *domain.LedgerEntryFilter) (*shared.Page[model.LedgerEntry], error) {
					return &shared.Page[model.LedgerEntry]{
						Data: []model.LedgerEntry{
							makeLedgerEntry(entryAID, 100000, 0, "1101", "Cash"),
						},
						Pagination: shared.Pagination{
							Size:         10,
							TotalCount:   11,
							CurrentPage:  2,
							PreviousPage: func() *int { v := 1; return &v }(),
							NextPage:     nil,
							TotalPage:    2,
						},
					}, nil
				},
			},
			paymentProductRepo:    &fakePaymentProductRepository{},
			studentObligationRepo: &fakeStudentObligationRepository{},
			paymentAllocationRepo: &fakePaymentAllocationRepository{},
			paymentOrderRepo:      &fakePaymentOrderRepository{},
		},
		&fakePaymentOrderService{},
	)

	page, err := svc.FindPaginate(context.Background(), parentAuthContext(), pageable, &domain.LedgerEntryFilter{})
	if !errors.Is(err, shared.ErrAuthUnauthorized) {
		t.Fatalf("expected unauthorized error, got %v", err)
	}
	if page == nil || len(page.Data) != 0 || page.Pagination.CurrentPage != pageable.Page {
		t.Fatalf("expected empty page on unauthorized access")
	}

	page, err = svc.FindPaginate(context.Background(), adminAuthContext(), pageable, &domain.LedgerEntryFilter{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(page.Data) != 1 || page.Data[0].AccountCode != "1101" {
		t.Fatalf("expected mapped ledger page data")
	}
}
