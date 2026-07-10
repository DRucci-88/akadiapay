package phase3_test

import (
	"akadia/domain"
	service "akadia/internal/payment/service"
	"akadia/internal/platform/security"
	"akadia/internal/shared"
	"akadia/model"
	"context"
	"errors"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type fakePaymentOrderRepository struct {
	findPaginateFn func(
		ctx context.Context,
		pageable *shared.Pageable,
		filter *domain.PaymentOrderFilter,
		authContext *security.AuthContext,
	) (*shared.Page[model.PaymentOrder], error)
	createFn       func(ctx context.Context, paymentOrder *model.PaymentOrder) error
	firstByIDFn    func(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (*model.PaymentOrder, error)
	updateStatusFn func(ctx context.Context, id uuid.UUID, status model.PaymentOrderStatus) (int, error)
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
	if f.findPaginateFn == nil {
		panic("unexpected call: FindPaginate")
	}
	return f.findPaginateFn(ctx, pageable, filter, authContext)
}

func (f *fakePaymentOrderRepository) Create(
	ctx context.Context,
	paymentOrder *model.PaymentOrder,
) error {
	if f.createFn == nil {
		panic("unexpected call: Create")
	}
	return f.createFn(ctx, paymentOrder)
}

func (f *fakePaymentOrderRepository) FirstByID(
	ctx context.Context,
	id uuid.UUID,
	tenantID uuid.UUID,
) (*model.PaymentOrder, error) {
	if f.firstByIDFn == nil {
		panic("unexpected call: FirstByID")
	}
	return f.firstByIDFn(ctx, id, tenantID)
}

func (f *fakePaymentOrderRepository) LockByID(
	ctx context.Context,
	id uuid.UUID,
	tenantID uuid.UUID,
) (*model.PaymentOrder, error) {
	panic("unexpected call: LockByID")
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
	panic("unexpected call: MarkLedgerPosted")
}

type fakeStudentObligationRepository struct {
	sumOutstandingByStudentIDFn func(ctx context.Context, studentID uuid.UUID) (float64, error)
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
	if f.sumOutstandingByStudentIDFn == nil {
		panic("unexpected call: SumOutstandingByStudentID")
	}
	return f.sumOutstandingByStudentIDFn(ctx, studentID)
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
	panic("unexpected call: FindByIDs")
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

type fakeStudentService struct {
	firstByIDFn func(ctx context.Context, id uuid.UUID) (*model.Student, error)
}

func (f *fakeStudentService) FirstByID(
	ctx context.Context,
	id uuid.UUID,
) (*model.Student, error) {
	if f.firstByIDFn == nil {
		panic("unexpected call: FirstByID")
	}
	return f.firstByIDFn(ctx, id)
}

func (f *fakeStudentService) FindByUserID(
	ctx context.Context,
	userID uuid.UUID,
) (*model.Student, error) {
	panic("unexpected call: FindByUserID")
}

type fakeParentStudentService struct {
	existsFn func(ctx context.Context, parentUserID uuid.UUID, studentID uuid.UUID) (bool, error)
}

func (f *fakeParentStudentService) ExistsByParentUserIDAndStudentID(
	ctx context.Context,
	parentUserID uuid.UUID,
	studentID uuid.UUID,
) (bool, error) {
	if f.existsFn == nil {
		panic("unexpected call: ExistsByParentUserIDAndStudentID")
	}
	return f.existsFn(ctx, parentUserID, studentID)
}

func TestPaymentOrderServiceCreateInitializesPendingOrder(t *testing.T) {
	orderNumberPattern := regexp.MustCompile(`^PO-\d{14}-[0-9a-f]{8}$`)
	tests := []struct {
		name          string
		authContext   *security.AuthContext
		paymentMethod model.PaymentOrderPaymentMethod
		parentLinked  bool
		expectedPayer uuid.UUID
	}{
		{
			name:          "admin can create cash order",
			authContext:   adminAuthContext(),
			paymentMethod: model.PaymentMethodCash,
			expectedPayer: adminUserID,
		},
		{
			name:          "linked parent can create virtual account order",
			authContext:   parentAuthContext(),
			paymentMethod: model.PaymentMethodVirtualAccount,
			parentLinked:  true,
			expectedPayer: parentUserID,
		},
		{
			name:          "student can create own qris order",
			authContext:   studentAuthContext(studentID),
			paymentMethod: model.PaymentMethodQRIS,
			expectedPayer: studentUserID,
		},
		{
			name:          "admin can create bank transfer order",
			authContext:   adminAuthContext(),
			paymentMethod: model.PaymentMethodBankTransfer,
			expectedPayer: adminUserID,
		},
		{
			name:          "admin can create credit card order",
			authContext:   adminAuthContext(),
			paymentMethod: model.PaymentMethodCreditCard,
			expectedPayer: adminUserID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var created *model.PaymentOrder
			sumCalls := 0
			manager := &fakePaymentRepositoryManager{
				studentObligationRepo: &fakeStudentObligationRepository{
					sumOutstandingByStudentIDFn: func(ctx context.Context, requestedStudentID uuid.UUID) (float64, error) {
						sumCalls++
						if requestedStudentID != studentID {
							t.Fatalf("expected student ID %s, got %s", studentID, requestedStudentID)
						}
						return 500000, nil
					},
				},
				paymentAllocationRepo: &fakePaymentAllocationRepository{},
				paymentOrderRepo: &fakePaymentOrderRepository{
					createFn: func(ctx context.Context, paymentOrder *model.PaymentOrder) error {
						created = paymentOrder
						return nil
					},
				},
			}
			svc := service.NewPaymentOrderService(
				manager,
				&fakeStudentService{
					firstByIDFn: func(ctx context.Context, id uuid.UUID) (*model.Student, error) {
						return makeStudent(id, tenantID), nil
					},
				},
				&fakeParentStudentService{
					existsFn: func(ctx context.Context, parentUserID uuid.UUID, requestedStudentID uuid.UUID) (bool, error) {
						return tt.parentLinked, nil
					},
				},
			)
			req := &domain.PaymentOrderCreate{
				StudentID:     studentID,
				Amount:        325000,
				PaymentMethod: tt.paymentMethod,
				PaymentDate:   time.Date(2026, time.July, 10, 0, 0, 0, 0, time.UTC),
				Notes:         "Installment payment",
			}

			res, err := svc.Create(context.Background(), tt.authContext, req)

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if created == nil || res != created {
				t.Fatalf("expected created payment order to be returned")
			}
			if sumCalls != 1 {
				t.Fatalf("expected one outstanding lookup, got %d", sumCalls)
			}
			if created.TenantID != tenantID {
				t.Fatalf("expected tenant ID %s, got %s", tenantID, created.TenantID)
			}
			if created.StudentID != studentID {
				t.Fatalf("expected student ID %s, got %s", studentID, created.StudentID)
			}
			if created.PaidByUserID != tt.expectedPayer {
				t.Fatalf("expected paid by user %s, got %s", tt.expectedPayer, created.PaidByUserID)
			}
			assertAmountEqual(t, 325000, created.TotalAmount)
			if created.Status != model.PaymentOrderStatusPending {
				t.Fatalf("expected pending status, got %s", created.Status)
			}
			if created.PaymentMethod != tt.paymentMethod {
				t.Fatalf("expected payment method %s, got %s", tt.paymentMethod, created.PaymentMethod)
			}
			if created.OrderDate != req.PaymentDate {
				t.Fatalf("expected order date %s, got %s", req.PaymentDate, created.OrderDate)
			}
			if created.Notes != req.Notes {
				t.Fatalf("expected notes to be preserved")
			}
			if created.OrderNumber == "" || !strings.HasPrefix(created.OrderNumber, "PO-") || !orderNumberPattern.MatchString(created.OrderNumber) {
				t.Fatalf("expected generated order number format, got %q", created.OrderNumber)
			}
		})
	}
}

func TestPaymentOrderServiceCreateRejectsInvalidRequests(t *testing.T) {
	validDate := time.Date(2026, time.July, 10, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name               string
		authContext        *security.AuthContext
		parentLinked       bool
		studentTenant      uuid.UUID
		req                *domain.PaymentOrderCreate
		totalOutstanding   float64
		expectedErr        error
		expectSumLookup    bool
		expectCreateCalled bool
	}{
		{
			name:        "payment date is required",
			authContext: adminAuthContext(),
			req: &domain.PaymentOrderCreate{
				StudentID:     studentID,
				Amount:        100000,
				PaymentMethod: model.PaymentMethodCash,
			},
			expectedErr: shared.ErrPaymentOrderDateRequired,
		},
		{
			name:        "amount must be positive",
			authContext: adminAuthContext(),
			req: &domain.PaymentOrderCreate{
				StudentID:     studentID,
				Amount:        0,
				PaymentMethod: model.PaymentMethodCash,
				PaymentDate:   validDate,
			},
			expectedErr: shared.ErrPaymentOrderAmountInvalid,
		},
		{
			name:        "payment method must be valid",
			authContext: adminAuthContext(),
			req: &domain.PaymentOrderCreate{
				StudentID:     studentID,
				Amount:        100000,
				PaymentMethod: model.PaymentOrderPaymentMethod("CHEQUE"),
				PaymentDate:   validDate,
			},
			expectedErr: shared.ErrPaymentOrderMethodInvalid,
		},
		{
			name:          "student must belong to tenant",
			authContext:   adminAuthContext(),
			studentTenant: otherTenantID,
			req: &domain.PaymentOrderCreate{
				StudentID:     studentID,
				Amount:        100000,
				PaymentMethod: model.PaymentMethodCash,
				PaymentDate:   validDate,
			},
			expectedErr: shared.ErrStudentNotFound,
		},
		{
			name:        "unlinked parent cannot create order",
			authContext: parentAuthContext(),
			req: &domain.PaymentOrderCreate{
				StudentID:     studentID,
				Amount:        100000,
				PaymentMethod: model.PaymentMethodCash,
				PaymentDate:   validDate,
			},
			expectedErr: shared.ErrStudentNotFound,
		},
		{
			name:        "student cannot create order for another student",
			authContext: studentAuthContext(secondStudentID),
			req: &domain.PaymentOrderCreate{
				StudentID:     studentID,
				Amount:        100000,
				PaymentMethod: model.PaymentMethodCash,
				PaymentDate:   validDate,
			},
			expectedErr: shared.ErrStudentNotFound,
		},
		{
			name:        "outstanding bill is required",
			authContext: adminAuthContext(),
			req: &domain.PaymentOrderCreate{
				StudentID:     studentID,
				Amount:        100000,
				PaymentMethod: model.PaymentMethodCash,
				PaymentDate:   validDate,
			},
			totalOutstanding: 0,
			expectedErr:      shared.ErrPaymentOrderOutstandingRequired,
			expectSumLookup:  true,
		},
		{
			name:        "amount cannot exceed total outstanding",
			authContext: adminAuthContext(),
			req: &domain.PaymentOrderCreate{
				StudentID:     studentID,
				Amount:        600000,
				PaymentMethod: model.PaymentMethodCash,
				PaymentDate:   validDate,
			},
			totalOutstanding: 500000,
			expectedErr:      shared.ErrPaymentOrderAmountExceedsOutstanding,
			expectSumLookup:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sumCalls := 0
			createCalls := 0
			studentTenant := tenantID
			if tt.studentTenant != uuid.Nil {
				studentTenant = tt.studentTenant
			}
			manager := &fakePaymentRepositoryManager{
				studentObligationRepo: &fakeStudentObligationRepository{
					sumOutstandingByStudentIDFn: func(ctx context.Context, studentID uuid.UUID) (float64, error) {
						sumCalls++
						return tt.totalOutstanding, nil
					},
				},
				paymentAllocationRepo: &fakePaymentAllocationRepository{},
				paymentOrderRepo: &fakePaymentOrderRepository{
					createFn: func(ctx context.Context, paymentOrder *model.PaymentOrder) error {
						createCalls++
						return nil
					},
				},
			}
			svc := service.NewPaymentOrderService(
				manager,
				&fakeStudentService{
					firstByIDFn: func(ctx context.Context, id uuid.UUID) (*model.Student, error) {
						return makeStudent(id, studentTenant), nil
					},
				},
				&fakeParentStudentService{
					existsFn: func(ctx context.Context, parentUserID uuid.UUID, studentID uuid.UUID) (bool, error) {
						return tt.parentLinked, nil
					},
				},
			)

			res, err := svc.Create(context.Background(), tt.authContext, tt.req)

			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
			}
			if res != nil {
				t.Fatalf("expected nil response on error")
			}
			if createCalls != 0 {
				t.Fatalf("expected create not to be called")
			}
			if tt.expectSumLookup && sumCalls != 1 {
				t.Fatalf("expected one outstanding lookup, got %d", sumCalls)
			}
			if !tt.expectSumLookup && sumCalls != 0 {
				t.Fatalf("expected outstanding lookup to be skipped, got %d calls", sumCalls)
			}
		})
	}
}

func TestPaymentOrderServiceFindByIDMasksUnauthorizedAccess(t *testing.T) {
	manager := &fakePaymentRepositoryManager{
		studentObligationRepo: &fakeStudentObligationRepository{},
		paymentAllocationRepo: &fakePaymentAllocationRepository{},
		paymentOrderRepo: &fakePaymentOrderRepository{
			firstByIDFn: func(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (*model.PaymentOrder, error) {
				return makePaymentOrder(studentID), nil
			},
		},
	}
	svc := service.NewPaymentOrderService(
		manager,
		&fakeStudentService{
			firstByIDFn: func(ctx context.Context, id uuid.UUID) (*model.Student, error) {
				return makeStudent(id, tenantID), nil
			},
		},
		&fakeParentStudentService{
			existsFn: func(ctx context.Context, parentUserID uuid.UUID, studentID uuid.UUID) (bool, error) {
				return false, nil
			},
		},
	)

	res, err := svc.FindByID(context.Background(), parentAuthContext(), orderID)

	if !errors.Is(err, shared.ErrPaymentOrderNotFound) {
		t.Fatalf("expected error %v, got %v", shared.ErrPaymentOrderNotFound, err)
	}
	if res != nil {
		t.Fatalf("expected nil result when access is denied")
	}
}

func TestPaymentOrderServiceCancelPendingOrder(t *testing.T) {
	current := makePaymentOrder(studentID)
	firstByIDCalls := 0
	var capturedStatus model.PaymentOrderStatus
	manager := &fakePaymentRepositoryManager{
		studentObligationRepo: &fakeStudentObligationRepository{},
		paymentAllocationRepo: &fakePaymentAllocationRepository{
			findByPaymentOrderIDFn: func(ctx context.Context, paymentOrderID uuid.UUID) ([]model.PaymentAllocation, error) {
				return nil, nil
			},
		},
		paymentOrderRepo: &fakePaymentOrderRepository{
			firstByIDFn: func(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (*model.PaymentOrder, error) {
				firstByIDCalls++
				if firstByIDCalls == 1 {
					return current, nil
				}
				updated := *current
				updated.Status = model.PaymentOrderStatusCancelled
				return &updated, nil
			},
			updateStatusFn: func(ctx context.Context, id uuid.UUID, status model.PaymentOrderStatus) (int, error) {
				capturedStatus = status
				return 1, nil
			},
		},
	}
	svc := service.NewPaymentOrderService(
		manager,
		&fakeStudentService{
			firstByIDFn: func(ctx context.Context, id uuid.UUID) (*model.Student, error) {
				return makeStudent(id, tenantID), nil
			},
		},
		&fakeParentStudentService{},
	)

	res, err := svc.Cancel(context.Background(), adminAuthContext(), orderID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if capturedStatus != model.PaymentOrderStatusCancelled {
		t.Fatalf("expected cancel to update status to cancelled, got %s", capturedStatus)
	}
	if res == nil || res.Status != model.PaymentOrderStatusCancelled {
		t.Fatalf("expected reloaded cancelled order")
	}
	if firstByIDCalls != 2 {
		t.Fatalf("expected order reload after cancel, got %d lookups", firstByIDCalls)
	}
}

func TestPaymentOrderServiceCancelRejectsProtectedStates(t *testing.T) {
	now := time.Now().UTC()
	tests := []struct {
		name        string
		order       *model.PaymentOrder
		allocations []model.PaymentAllocation
		expectedErr error
	}{
		{
			name: "ledger posted order cannot be cancelled",
			order: func() *model.PaymentOrder {
				order := makePaymentOrder(studentID)
				order.LedgerPostedAt = &now
				return order
			}(),
			expectedErr: shared.ErrPostedPaymentCannotBeCancelled,
		},
		{
			name: "completed order cannot be cancelled",
			order: func() *model.PaymentOrder {
				order := makePaymentOrder(studentID)
				order.Status = model.PaymentOrderStatusCompleted
				return order
			}(),
			expectedErr: shared.ErrPaymentOrderStatusInvalid,
		},
		{
			name:  "allocated order cannot be cancelled",
			order: makePaymentOrder(studentID),
			allocations: []model.PaymentAllocation{
				*makePaymentAllocation(),
			},
			expectedErr: shared.ErrPaymentOrderAllocated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateCalls := 0
			manager := &fakePaymentRepositoryManager{
				studentObligationRepo: &fakeStudentObligationRepository{},
				paymentAllocationRepo: &fakePaymentAllocationRepository{
					findByPaymentOrderIDFn: func(ctx context.Context, paymentOrderID uuid.UUID) ([]model.PaymentAllocation, error) {
						return tt.allocations, nil
					},
				},
				paymentOrderRepo: &fakePaymentOrderRepository{
					firstByIDFn: func(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (*model.PaymentOrder, error) {
						return tt.order, nil
					},
					updateStatusFn: func(ctx context.Context, id uuid.UUID, status model.PaymentOrderStatus) (int, error) {
						updateCalls++
						return 1, nil
					},
				},
			}
			svc := service.NewPaymentOrderService(
				manager,
				&fakeStudentService{
					firstByIDFn: func(ctx context.Context, id uuid.UUID) (*model.Student, error) {
						return makeStudent(id, tenantID), nil
					},
				},
				&fakeParentStudentService{},
			)

			res, err := svc.Cancel(context.Background(), adminAuthContext(), orderID)

			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
			}
			if res != nil {
				t.Fatalf("expected nil result on protected-state rejection")
			}
			if updateCalls != 0 {
				t.Fatalf("expected status update not to be called")
			}
		})
	}
}
