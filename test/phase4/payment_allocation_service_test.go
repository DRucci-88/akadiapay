package phase4_test

import (
	"akadia/domain"
	service "akadia/internal/payment/service"
	"akadia/internal/platform/security"
	"akadia/internal/shared"
	"akadia/model"
	"context"
	"errors"
	"sort"
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
	createBatchFn          func(ctx context.Context, paymentAllocations []model.PaymentAllocation) error
	findByPaymentOrderIDFn func(ctx context.Context, paymentOrderID uuid.UUID) ([]model.PaymentAllocation, error)
}

func (f *fakePaymentAllocationRepository) CreateBatch(
	ctx context.Context,
	paymentAllocations []model.PaymentAllocation,
) error {
	if f.createBatchFn == nil {
		panic("unexpected call: CreateBatch")
	}
	return f.createBatchFn(ctx, paymentAllocations)
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
	lockByIDsFn        func(ctx context.Context, ids []uuid.UUID) ([]model.StudentObligation, error)
	updateSettlementFn func(ctx context.Context, id uuid.UUID, outstandingAmount float64, status model.StudentObligationStatus) (int, error)
	findByIDsFn        func(ctx context.Context, ids []uuid.UUID) ([]model.StudentObligation, error)
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
	if f.lockByIDsFn == nil {
		panic("unexpected call: LockByIDs")
	}
	return f.lockByIDsFn(ctx, ids)
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
	if f.updateSettlementFn == nil {
		panic("unexpected call: UpdateSettlement")
	}
	return f.updateSettlementFn(ctx, id, outstandingAmount, status)
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
	existsByPaymentOrderIDFn func(ctx context.Context, paymentOrderID uuid.UUID, tenantID uuid.UUID) (bool, error)
	createInBatchesFn        func(ctx context.Context, entries []model.LedgerEntry, batchSize int) error
}

func (f *fakeLedgerEntryRepository) FindByPaymentOrderID(
	ctx context.Context,
	paymentOrderID uuid.UUID,
	tenantID uuid.UUID,
) ([]model.LedgerEntry, error) {
	panic("unexpected call: FindByPaymentOrderID")
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
	panic("unexpected call: FindPaginate")
}

func TestPaymentAllocationServiceAllocateRejectsRequestValidation(t *testing.T) {
	tests := []struct {
		name         string
		paymentOrder *model.PaymentOrder
		req          *domain.PaymentAllocationAllocate
		expectedErr  error
	}{
		{
			name: "payment order must be pending",
			paymentOrder: func() *model.PaymentOrder {
				order := makePaymentOrder(500000, model.PaymentMethodCash)
				order.Status = model.PaymentOrderStatusCompleted
				return order
			}(),
			req: &domain.PaymentAllocationAllocate{
				Allocations: []domain.PaymentAllocationCreate{{StudentObligationID: obligationAID, AllocatedAmount: 100000}},
			},
			expectedErr: shared.ErrPaymentOrderStatusInvalid,
		},
		{
			name:         "allocation list is required",
			paymentOrder: makePaymentOrder(500000, model.PaymentMethodCash),
			req:          &domain.PaymentAllocationAllocate{},
			expectedErr:  shared.ErrPaymentAllocationRequired,
		},
		{
			name:         "allocated amount must be positive",
			paymentOrder: makePaymentOrder(500000, model.PaymentMethodCash),
			req: &domain.PaymentAllocationAllocate{
				Allocations: []domain.PaymentAllocationCreate{{StudentObligationID: obligationAID, AllocatedAmount: 0}},
			},
			expectedErr: shared.ErrPaymentAllocationAmountInvalid,
		},
		{
			name:         "duplicate obligations in request are rejected",
			paymentOrder: makePaymentOrder(500000, model.PaymentMethodCash),
			req: &domain.PaymentAllocationAllocate{
				Allocations: []domain.PaymentAllocationCreate{
					{StudentObligationID: obligationAID, AllocatedAmount: 100000},
					{StudentObligationID: obligationAID, AllocatedAmount: 200000},
				},
			},
			expectedErr: shared.ErrPaymentAllocationDuplicateObligation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			svc := service.NewPaymentAllocationService(
				manager,
				&fakePaymentOrderService{
					findByIDFn: func(ctx context.Context, authContext *security.AuthContext, id uuid.UUID) (*model.PaymentOrder, error) {
						return tt.paymentOrder, nil
					},
				},
			)

			res, err := svc.Allocate(context.Background(), adminAuthContext(), orderID, tt.req)

			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
			}
			if res != nil {
				t.Fatalf("expected nil result on validation failure")
			}
			if transactionCalls != 0 {
				t.Fatalf("expected transaction not to start")
			}
		})
	}
}

func TestPaymentAllocationServiceAllocateUsesLocksAndUpdatesSettlements(t *testing.T) {
	order := makePaymentOrder(500000, model.PaymentMethodCash)
	policyA := makePaymentPolicy(policyAID)
	policyA.AutoCloseObligation = true
	policyB := makePaymentPolicy(policyBID)
	policyB.MinimumAmount = 100000
	productA := makePaymentProduct(productAID, policyA, "SPP-JUL", "Tuition")
	productB := makePaymentProduct(productBID, policyB, "SPP-AUG", "Tuition")
	obligationA := makeStudentObligation(obligationAID, studentID, productAID, 100000)
	obligationB := makeStudentObligation(obligationBID, studentID, productBID, 500000)

	transactionCalls := 0
	lockByOrderCalls := 0
	lockByIDsCalls := 0
	var lockedIDs []uuid.UUID
	var createdAllocations []model.PaymentAllocation
	type settlementCall struct {
		id          uuid.UUID
		outstanding float64
		status      model.StudentObligationStatus
	}
	var settlements []settlementCall
	var updatedOrderStatus model.PaymentOrderStatus

	manager := &fakePaymentRepositoryManager{
		ledgerEntryRepo: &fakeLedgerEntryRepository{
			existsByPaymentOrderIDFn: func(ctx context.Context, paymentOrderID uuid.UUID, tenantID uuid.UUID) (bool, error) {
				t.Fatalf("ledger should not be posted for partial order")
				return false, nil
			},
			createInBatchesFn: func(ctx context.Context, entries []model.LedgerEntry, batchSize int) error {
				t.Fatalf("ledger should not be created for partial order")
				return nil
			},
		},
		paymentProductRepo: &fakePaymentProductRepository{
			findByIDsIncludingDeletedFn: func(ctx context.Context, ids []uuid.UUID) ([]model.PaymentProduct, error) {
				return []model.PaymentProduct{*productA, *productB}, nil
			},
		},
		studentObligationRepo: &fakeStudentObligationRepository{
			lockByIDsFn: func(ctx context.Context, ids []uuid.UUID) ([]model.StudentObligation, error) {
				lockByIDsCalls++
				lockedIDs = append([]uuid.UUID(nil), ids...)
				return []model.StudentObligation{obligationA, obligationB}, nil
			},
			updateSettlementFn: func(ctx context.Context, id uuid.UUID, outstandingAmount float64, status model.StudentObligationStatus) (int, error) {
				settlements = append(settlements, settlementCall{id: id, outstanding: outstandingAmount, status: status})
				return 1, nil
			},
			findByIDsFn: func(ctx context.Context, ids []uuid.UUID) ([]model.StudentObligation, error) {
				t.Fatalf("FindByIDs should not be called when order remains pending")
				return nil, nil
			},
		},
		paymentAllocationRepo: &fakePaymentAllocationRepository{
			findByPaymentOrderIDFn: func(ctx context.Context, paymentOrderID uuid.UUID) ([]model.PaymentAllocation, error) {
				return createdAllocations, nil
			},
			createBatchFn: func(ctx context.Context, paymentAllocations []model.PaymentAllocation) error {
				createdAllocations = append([]model.PaymentAllocation(nil), paymentAllocations...)
				return nil
			},
		},
		paymentOrderRepo: &fakePaymentOrderRepository{
			lockByIDFn: func(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (*model.PaymentOrder, error) {
				lockByOrderCalls++
				return order, nil
			},
			updateStatusFn: func(ctx context.Context, id uuid.UUID, status model.PaymentOrderStatus) (int, error) {
				updatedOrderStatus = status
				order.Status = status
				return 1, nil
			},
			markLedgerPostedFn: func(ctx context.Context, id uuid.UUID, ledgerPostedAt *time.Time) (int, error) {
				t.Fatalf("ledger should not be marked for partial order")
				return 0, nil
			},
		},
	}
	manager.transactionFn = func(ctx context.Context, fn func(repo domain.RepositoryManagerPayment) error) error {
		transactionCalls++
		return fn(manager)
	}
	svc := service.NewPaymentAllocationService(
		manager,
		&fakePaymentOrderService{
			findByIDFn: func(ctx context.Context, authContext *security.AuthContext, id uuid.UUID) (*model.PaymentOrder, error) {
				return order, nil
			},
		},
	)

	res, err := svc.Allocate(context.Background(), adminAuthContext(), orderID, &domain.PaymentAllocationAllocate{
		Allocations: []domain.PaymentAllocationCreate{
			{StudentObligationID: obligationBID, AllocatedAmount: 200000},
			{StudentObligationID: obligationAID, AllocatedAmount: 100000},
		},
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if transactionCalls != 1 || lockByOrderCalls != 1 || lockByIDsCalls != 1 {
		t.Fatalf("expected one transaction and one lock path, got tx=%d orderLocks=%d obligationLocks=%d", transactionCalls, lockByOrderCalls, lockByIDsCalls)
	}
	expectedLockedIDs := []uuid.UUID{obligationAID, obligationBID}
	if len(lockedIDs) != len(expectedLockedIDs) {
		t.Fatalf("expected %d locked ids, got %d", len(expectedLockedIDs), len(lockedIDs))
	}
	for i := range expectedLockedIDs {
		if lockedIDs[i] != expectedLockedIDs[i] {
			t.Fatalf("expected sorted lock ids %v, got %v", expectedLockedIDs, lockedIDs)
		}
	}
	if len(createdAllocations) != 2 {
		t.Fatalf("expected two persisted allocations, got %d", len(createdAllocations))
	}
	if updatedOrderStatus != model.PaymentOrderStatusPending {
		t.Fatalf("expected partial allocation to keep order pending, got %s", updatedOrderStatus)
	}
	if len(settlements) != 2 {
		t.Fatalf("expected two settlement updates, got %d", len(settlements))
	}
	sort.Slice(settlements, func(i, j int) bool { return settlements[i].id.String() < settlements[j].id.String() })
	if settlements[0].id != obligationAID || settlements[0].status != model.StudentObligationStatusClosed {
		t.Fatalf("expected first obligation to close after full payment, got %+v", settlements[0])
	}
	assertAmountEqual(t, 0, settlements[0].outstanding)
	if settlements[1].id != obligationBID || settlements[1].status != model.StudentObligationStatusPartial {
		t.Fatalf("expected second obligation to become partial, got %+v", settlements[1])
	}
	assertAmountEqual(t, 300000, settlements[1].outstanding)
	if res == nil {
		t.Fatalf("expected allocation summary result")
	}
	assertAmountEqual(t, 300000, res.TotalAllocated)
	assertAmountEqual(t, 200000, res.RemainingAmount)
	if res.OrderStatus != model.PaymentOrderStatusPending {
		t.Fatalf("expected pending order status in summary, got %s", res.OrderStatus)
	}
}

func TestPaymentAllocationServiceAllocateCompletesOrderAndPostsLedger(t *testing.T) {
	order := makePaymentOrder(300000, model.PaymentMethodQRIS)
	policy := makePaymentPolicy(policyAID)
	policy.AutoCloseObligation = false
	product := makePaymentProduct(productAID, policy, "UNKNOWN-FEE", "Misc fee")
	product.RevenueAccountCode = ""
	product.RevenueAccountName = ""
	obligation := makeStudentObligation(obligationAID, studentID, productAID, 300000)

	var allocationsState []model.PaymentAllocation
	var createdEntries []model.LedgerEntry
	markLedgerCalls := 0
	settlementStatus := model.StudentObligationStatusPending
	currentOrder := order

	manager := &fakePaymentRepositoryManager{
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
				return []model.PaymentProduct{*product}, nil
			},
		},
		studentObligationRepo: &fakeStudentObligationRepository{
			lockByIDsFn: func(ctx context.Context, ids []uuid.UUID) ([]model.StudentObligation, error) {
				return []model.StudentObligation{obligation}, nil
			},
			updateSettlementFn: func(ctx context.Context, id uuid.UUID, outstandingAmount float64, status model.StudentObligationStatus) (int, error) {
				settlementStatus = status
				obligation.OutstandingAmount = outstandingAmount
				obligation.Status = status
				return 1, nil
			},
			findByIDsFn: func(ctx context.Context, ids []uuid.UUID) ([]model.StudentObligation, error) {
				return []model.StudentObligation{obligation}, nil
			},
		},
		paymentAllocationRepo: &fakePaymentAllocationRepository{
			findByPaymentOrderIDFn: func(ctx context.Context, paymentOrderID uuid.UUID) ([]model.PaymentAllocation, error) {
				return allocationsState, nil
			},
			createBatchFn: func(ctx context.Context, paymentAllocations []model.PaymentAllocation) error {
				allocationsState = append([]model.PaymentAllocation(nil), paymentAllocations...)
				return nil
			},
		},
		paymentOrderRepo: &fakePaymentOrderRepository{
			lockByIDFn: func(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (*model.PaymentOrder, error) {
				return currentOrder, nil
			},
			updateStatusFn: func(ctx context.Context, id uuid.UUID, status model.PaymentOrderStatus) (int, error) {
				currentOrder.Status = status
				return 1, nil
			},
			markLedgerPostedFn: func(ctx context.Context, id uuid.UUID, ledgerPostedAt *time.Time) (int, error) {
				markLedgerCalls++
				currentOrder.LedgerPostedAt = ledgerPostedAt
				return 1, nil
			},
		},
	}
	svc := service.NewPaymentAllocationService(
		manager,
		&fakePaymentOrderService{
			findByIDFn: func(ctx context.Context, authContext *security.AuthContext, id uuid.UUID) (*model.PaymentOrder, error) {
				return currentOrder, nil
			},
		},
	)

	res, err := svc.Allocate(context.Background(), adminAuthContext(), orderID, &domain.PaymentAllocationAllocate{
		Allocations: []domain.PaymentAllocationCreate{
			{StudentObligationID: obligationAID, AllocatedAmount: 300000},
		},
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res == nil || res.OrderStatus != model.PaymentOrderStatusCompleted {
		t.Fatalf("expected completed allocation summary")
	}
	if settlementStatus != model.StudentObligationStatusPaid {
		t.Fatalf("expected fully paid obligation to become PAID when auto close disabled, got %s", settlementStatus)
	}
	if len(createdEntries) != 2 {
		t.Fatalf("expected one debit and one credit ledger entry, got %d", len(createdEntries))
	}
	assertAmountEqual(t, 300000, createdEntries[0].Debit)
	if createdEntries[0].AccountCode != "1103" || createdEntries[0].AccountName != "QRIS Clearing" {
		t.Fatalf("expected QRIS debit account mapping, got %+v", createdEntries[0])
	}
	assertAmountEqual(t, 300000, createdEntries[1].Credit)
	if createdEntries[1].AccountCode != "4199" || createdEntries[1].AccountName != "Other Education Revenue" {
		t.Fatalf("expected default revenue fallback account, got %+v", createdEntries[1])
	}
	if markLedgerCalls != 1 || currentOrder.LedgerPostedAt == nil {
		t.Fatalf("expected ledger posted marker to be set once")
	}
}

func TestPaymentAllocationServiceAllocateRejectsBusinessRuleViolations(t *testing.T) {
	tests := []struct {
		name                string
		orderTotal          float64
		existingAllocations []model.PaymentAllocation
		obligation          model.StudentObligation
		product             *model.PaymentProduct
		req                 *domain.PaymentAllocationAllocate
		expectedErr         error
	}{
		{
			name:       "existing allocation cannot be duplicated",
			orderTotal: 500000,
			existingAllocations: []model.PaymentAllocation{
				makePaymentAllocation(obligationAID, 100000),
			},
			obligation: makeStudentObligation(obligationAID, studentID, productAID, 500000),
			product:    makePaymentProduct(productAID, makePaymentPolicy(policyAID), "SPP", "Tuition"),
			req: &domain.PaymentAllocationAllocate{
				Allocations: []domain.PaymentAllocationCreate{{StudentObligationID: obligationAID, AllocatedAmount: 200000}},
			},
			expectedErr: shared.ErrPaymentAllocationDuplicateObligation,
		},
		{
			name:       "total allocation cannot exceed order amount",
			orderTotal: 250000,
			existingAllocations: []model.PaymentAllocation{
				makePaymentAllocation(obligationBID, 100000),
			},
			obligation: makeStudentObligation(obligationAID, studentID, productAID, 500000),
			product:    makePaymentProduct(productAID, makePaymentPolicy(policyAID), "SPP", "Tuition"),
			req: &domain.PaymentAllocationAllocate{
				Allocations: []domain.PaymentAllocationCreate{{StudentObligationID: obligationAID, AllocatedAmount: 200000}},
			},
			expectedErr: shared.ErrPaymentAllocationTotalExceedsOrder,
		},
		{
			name:       "allocation cannot exceed obligation outstanding",
			orderTotal: 500000,
			obligation: makeStudentObligation(obligationAID, studentID, productAID, 100000),
			product:    makePaymentProduct(productAID, makePaymentPolicy(policyAID), "SPP", "Tuition"),
			req: &domain.PaymentAllocationAllocate{
				Allocations: []domain.PaymentAllocationCreate{{StudentObligationID: obligationAID, AllocatedAmount: 150000}},
			},
			expectedErr: shared.ErrPaymentAllocationAmountExceedsOutstanding,
		},
		{
			name:       "non partial policy requires full payment",
			orderTotal: 500000,
			obligation: makeStudentObligation(obligationAID, studentID, productAID, 100000),
			product: func() *model.PaymentProduct {
				policy := makePaymentPolicy(policyAID)
				policy.AllowPartial = false
				return makePaymentProduct(productAID, policy, "SPP", "Tuition")
			}(),
			req: &domain.PaymentAllocationAllocate{
				Allocations: []domain.PaymentAllocationCreate{{StudentObligationID: obligationAID, AllocatedAmount: 50000}},
			},
			expectedErr: shared.ErrPaymentAllocationFullPaymentRequired,
		},
		{
			name:       "minimum amount is enforced",
			orderTotal: 500000,
			obligation: makeStudentObligation(obligationAID, studentID, productAID, 500000),
			product: func() *model.PaymentProduct {
				policy := makePaymentPolicy(policyAID)
				policy.MinimumAmount = 100000
				return makePaymentProduct(productAID, policy, "SPP", "Tuition")
			}(),
			req: &domain.PaymentAllocationAllocate{
				Allocations: []domain.PaymentAllocationCreate{{StudentObligationID: obligationAID, AllocatedAmount: 50000}},
			},
			expectedErr: shared.ErrPaymentAllocationBelowMinimumAmount,
		},
		{
			name:       "minimum percentage is enforced",
			orderTotal: 500000,
			obligation: makeStudentObligation(obligationAID, studentID, productAID, 500000),
			product: func() *model.PaymentProduct {
				policy := makePaymentPolicy(policyAID)
				policy.MinimumPercentage = 60
				return makePaymentProduct(productAID, policy, "SPP", "Tuition")
			}(),
			req: &domain.PaymentAllocationAllocate{
				Allocations: []domain.PaymentAllocationCreate{{StudentObligationID: obligationAID, AllocatedAmount: 200000}},
			},
			expectedErr: shared.ErrPaymentAllocationBelowMinimumPercentage,
		},
		{
			name:       "obligation from another student is rejected",
			orderTotal: 500000,
			obligation: makeStudentObligation(obligationAID, otherStudentID, productAID, 500000),
			product:    makePaymentProduct(productAID, makePaymentPolicy(policyAID), "SPP", "Tuition"),
			req: &domain.PaymentAllocationAllocate{
				Allocations: []domain.PaymentAllocationCreate{{StudentObligationID: obligationAID, AllocatedAmount: 100000}},
			},
			expectedErr: shared.ErrStudentObligationNotFound,
		},
		{
			name:       "missing payment product is rejected",
			orderTotal: 500000,
			obligation: makeStudentObligation(obligationAID, studentID, productAID, 500000),
			product:    nil,
			req: &domain.PaymentAllocationAllocate{
				Allocations: []domain.PaymentAllocationCreate{{StudentObligationID: obligationAID, AllocatedAmount: 100000}},
			},
			expectedErr: shared.ErrPaymentProductNotFound,
		},
		{
			name:       "missing payment policy is rejected",
			orderTotal: 500000,
			obligation: makeStudentObligation(obligationAID, studentID, productAID, 500000),
			product: func() *model.PaymentProduct {
				product := makePaymentProduct(productAID, makePaymentPolicy(policyAID), "SPP", "Tuition")
				product.PaymentPolicy = nil
				return product
			}(),
			req: &domain.PaymentAllocationAllocate{
				Allocations: []domain.PaymentAllocationCreate{{StudentObligationID: obligationAID, AllocatedAmount: 100000}},
			},
			expectedErr: shared.ErrPaymentPolicyNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currentOrder := makePaymentOrder(tt.orderTotal, model.PaymentMethodCash)
			manager := &fakePaymentRepositoryManager{
				ledgerEntryRepo: &fakeLedgerEntryRepository{
					existsByPaymentOrderIDFn: func(ctx context.Context, paymentOrderID uuid.UUID, tenantID uuid.UUID) (bool, error) {
						t.Fatalf("ledger should not be reached on validation failure")
						return false, nil
					},
					createInBatchesFn: func(ctx context.Context, entries []model.LedgerEntry, batchSize int) error {
						t.Fatalf("ledger should not be reached on validation failure")
						return nil
					},
				},
				paymentProductRepo: &fakePaymentProductRepository{
					findByIDsIncludingDeletedFn: func(ctx context.Context, ids []uuid.UUID) ([]model.PaymentProduct, error) {
						if tt.product == nil {
							return nil, nil
						}
						return []model.PaymentProduct{*tt.product}, nil
					},
				},
				studentObligationRepo: &fakeStudentObligationRepository{
					lockByIDsFn: func(ctx context.Context, ids []uuid.UUID) ([]model.StudentObligation, error) {
						return []model.StudentObligation{tt.obligation}, nil
					},
					updateSettlementFn: func(ctx context.Context, id uuid.UUID, outstandingAmount float64, status model.StudentObligationStatus) (int, error) {
						t.Fatalf("settlement should not be written on validation failure")
						return 0, nil
					},
					findByIDsFn: func(ctx context.Context, ids []uuid.UUID) ([]model.StudentObligation, error) {
						t.Fatalf("ledger should not be reached on validation failure")
						return nil, nil
					},
				},
				paymentAllocationRepo: &fakePaymentAllocationRepository{
					findByPaymentOrderIDFn: func(ctx context.Context, paymentOrderID uuid.UUID) ([]model.PaymentAllocation, error) {
						return tt.existingAllocations, nil
					},
					createBatchFn: func(ctx context.Context, paymentAllocations []model.PaymentAllocation) error {
						t.Fatalf("allocations should not be created on validation failure")
						return nil
					},
				},
				paymentOrderRepo: &fakePaymentOrderRepository{
					lockByIDFn: func(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (*model.PaymentOrder, error) {
						return currentOrder, nil
					},
					updateStatusFn: func(ctx context.Context, id uuid.UUID, status model.PaymentOrderStatus) (int, error) {
						t.Fatalf("order status should not be updated on validation failure")
						return 0, nil
					},
					markLedgerPostedFn: func(ctx context.Context, id uuid.UUID, ledgerPostedAt *time.Time) (int, error) {
						t.Fatalf("ledger marker should not be written on validation failure")
						return 0, nil
					},
				},
			}
			svc := service.NewPaymentAllocationService(
				manager,
				&fakePaymentOrderService{
					findByIDFn: func(ctx context.Context, authContext *security.AuthContext, id uuid.UUID) (*model.PaymentOrder, error) {
						return currentOrder, nil
					},
				},
			)

			res, err := svc.Allocate(context.Background(), adminAuthContext(), orderID, tt.req)

			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
			}
			if res != nil {
				t.Fatalf("expected nil result on business rule violation")
			}
		})
	}
}

func TestPaymentAllocationServiceFindByPaymentOrderIDBuildsSummary(t *testing.T) {
	currentOrder := makePaymentOrder(500000, model.PaymentMethodCash)
	allocationA := makePaymentAllocation(obligationAID, 125000)
	allocationB := makePaymentAllocation(obligationBID, 175000)
	svc := service.NewPaymentAllocationService(
		&fakePaymentRepositoryManager{
			ledgerEntryRepo:       &fakeLedgerEntryRepository{},
			paymentProductRepo:    &fakePaymentProductRepository{},
			studentObligationRepo: &fakeStudentObligationRepository{},
			paymentAllocationRepo: &fakePaymentAllocationRepository{
				findByPaymentOrderIDFn: func(ctx context.Context, paymentOrderID uuid.UUID) ([]model.PaymentAllocation, error) {
					return []model.PaymentAllocation{allocationA, allocationB}, nil
				},
			},
			paymentOrderRepo: &fakePaymentOrderRepository{},
		},
		&fakePaymentOrderService{
			findByIDFn: func(ctx context.Context, authContext *security.AuthContext, id uuid.UUID) (*model.PaymentOrder, error) {
				return currentOrder, nil
			},
		},
	)

	res, err := svc.FindByPaymentOrderID(context.Background(), adminAuthContext(), orderID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res == nil {
		t.Fatalf("expected payment allocation summary")
	}
	assertAmountEqual(t, 300000, res.TotalAllocated)
	assertAmountEqual(t, 200000, res.RemainingAmount)
	if res.OrderStatus != currentOrder.Status {
		t.Fatalf("expected order status %s, got %s", currentOrder.Status, res.OrderStatus)
	}
	if len(res.Allocations) != 2 {
		t.Fatalf("expected two allocation responses, got %d", len(res.Allocations))
	}
}
