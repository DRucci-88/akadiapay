package phase2_test

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

type fakeStudentObligationRepository struct {
	createFn                     func(ctx context.Context, studentObligation *model.StudentObligation) error
	createBatchFn                func(ctx context.Context, studentObligations []model.StudentObligation) error
	sumOutstandingByStudentIDFn  func(ctx context.Context, studentID uuid.UUID) (float64, error)
	findOutstandingByStudentIDFn func(ctx context.Context, studentID uuid.UUID) ([]model.StudentObligation, error)
	firstByIDFn                  func(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (*model.StudentObligation, error)
	updateFn                     func(ctx context.Context, id uuid.UUID, req *domain.StudentObligationUpdate) (int, error)
	deleteFn                     func(ctx context.Context, id uuid.UUID) (int, error)
	hasPaymentAllocationsFn      func(ctx context.Context, id uuid.UUID) (bool, error)
	existsActiveFn               func(ctx context.Context, studentID uuid.UUID, paymentProductID uuid.UUID, period time.Time) (bool, error)
}

func (f *fakeStudentObligationRepository) Create(
	ctx context.Context,
	studentObligation *model.StudentObligation,
) error {
	if f.createFn == nil {
		panic("unexpected call: Create")
	}
	return f.createFn(ctx, studentObligation)
}

func (f *fakeStudentObligationRepository) CreateBatch(
	ctx context.Context,
	studentObligations []model.StudentObligation,
) error {
	if f.createBatchFn == nil {
		panic("unexpected call: CreateBatch")
	}
	return f.createBatchFn(ctx, studentObligations)
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
	if f.findOutstandingByStudentIDFn == nil {
		panic("unexpected call: FindOutstandingByStudentID")
	}
	return f.findOutstandingByStudentIDFn(ctx, studentID)
}

func (f *fakeStudentObligationRepository) FirstByID(
	ctx context.Context,
	id uuid.UUID,
	tenantID uuid.UUID,
) (*model.StudentObligation, error) {
	if f.firstByIDFn == nil {
		panic("unexpected call: FirstByID")
	}
	return f.firstByIDFn(ctx, id, tenantID)
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
	if f.updateFn == nil {
		panic("unexpected call: Update")
	}
	return f.updateFn(ctx, id, req)
}

func (f *fakeStudentObligationRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
) (int, error) {
	if f.deleteFn == nil {
		panic("unexpected call: Delete")
	}
	return f.deleteFn(ctx, id)
}

func (f *fakeStudentObligationRepository) HasPaymentAllocations(
	ctx context.Context,
	id uuid.UUID,
) (bool, error) {
	if f.hasPaymentAllocationsFn == nil {
		panic("unexpected call: HasPaymentAllocations")
	}
	return f.hasPaymentAllocationsFn(ctx, id)
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
	if f.existsActiveFn == nil {
		panic("unexpected call: ExistsActiveByStudentIDAndPaymentProductIDAndPeriod")
	}
	return f.existsActiveFn(ctx, studentID, paymentProductID, period)
}

type fakeStudentService struct {
	firstByIDFn  func(ctx context.Context, id uuid.UUID) (*model.Student, error)
	findByUserID func(ctx context.Context, userID uuid.UUID) (*model.Student, error)
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
	if f.findByUserID == nil {
		panic("unexpected call: FindByUserID")
	}
	return f.findByUserID(ctx, userID)
}

type fakePaymentProductService struct {
	findByIDFn func(
		ctx context.Context,
		authContext *security.AuthContext,
		id uuid.UUID,
		preloads ...model.PaymentProductPreload,
	) (*model.PaymentProduct, error)
}

func (f *fakePaymentProductService) FindPaginate(
	ctx context.Context,
	pageable *shared.Pageable,
	filter *domain.PaymentProductFilter,
	authContext *security.AuthContext,
) (*shared.Page[domain.PaymentProductResponse], error) {
	panic("unexpected call: FindPaginate")
}

func (f *fakePaymentProductService) FindByID(
	ctx context.Context,
	authContext *security.AuthContext,
	id uuid.UUID,
	preloads ...model.PaymentProductPreload,
) (*model.PaymentProduct, error) {
	if f.findByIDFn == nil {
		panic("unexpected call: FindByID")
	}
	return f.findByIDFn(ctx, authContext, id, preloads...)
}

func (f *fakePaymentProductService) Create(
	ctx context.Context,
	authContext *security.AuthContext,
	req *domain.PaymentProductCreate,
) (*model.PaymentProduct, error) {
	panic("unexpected call: Create")
}

func (f *fakePaymentProductService) Update(
	ctx context.Context,
	authContext *security.AuthContext,
	id uuid.UUID,
	req *domain.PaymentProductUpdate,
) (*model.PaymentProduct, error) {
	panic("unexpected call: Update")
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

func TestStudentObligationServiceCreateInitializesBillingState(t *testing.T) {
	dueDate := time.Date(2026, time.July, 15, 10, 30, 0, 0, time.UTC)
	var created *model.StudentObligation
	var requestedPreloads []model.PaymentProductPreload
	repo := &fakeStudentObligationRepository{
		createFn: func(ctx context.Context, studentObligation *model.StudentObligation) error {
			created = studentObligation
			return nil
		},
		existsActiveFn: func(ctx context.Context, requestedStudentID uuid.UUID, requestedPaymentProductID uuid.UUID, period time.Time) (bool, error) {
			if requestedStudentID != studentID {
				t.Fatalf("expected student ID %s, got %s", studentID, requestedStudentID)
			}
			if requestedPaymentProductID != productID {
				t.Fatalf("expected payment product ID %s, got %s", productID, requestedPaymentProductID)
			}
			expectedPeriod := time.Date(2026, time.July, 1, 0, 0, 0, 0, time.UTC)
			if !period.Equal(expectedPeriod) {
				t.Fatalf("expected normalized period %s, got %s", expectedPeriod, period)
			}
			return false, nil
		},
	}
	svc := service.NewStudentObligationService(
		&fakePaymentRepositoryManager{studentObligationRepo: repo},
		&fakeStudentService{
			firstByIDFn: func(ctx context.Context, id uuid.UUID) (*model.Student, error) {
				return makeStudent(id, tenantID), nil
			},
		},
		&fakePaymentProductService{
			findByIDFn: func(
				ctx context.Context,
				authContext *security.AuthContext,
				id uuid.UUID,
				preloads ...model.PaymentProductPreload,
			) (*model.PaymentProduct, error) {
				requestedPreloads = append([]model.PaymentProductPreload(nil), preloads...)
				return makePaymentProduct(), nil
			},
		},
		&fakeParentStudentService{},
	)
	req := &domain.StudentObligationCreate{
		StudentID:        studentID,
		PaymentProductID: productID,
		DueDate:          dueDate,
		Amount:           floatPtr(275000),
		Notes:            "July arrears",
	}

	res, err := svc.Create(context.Background(), adminAuthContext(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if created == nil {
		t.Fatalf("expected repository create to be called")
	}
	if res != created {
		t.Fatalf("expected response to be the created obligation")
	}
	if len(requestedPreloads) != 1 || requestedPreloads[0] != model.PaymentProductPreloadPaymentPolicy {
		t.Fatalf("expected payment product lookup with payment policy preload")
	}
	expectedPeriod := time.Date(2026, time.July, 1, 0, 0, 0, 0, time.UTC)
	if !created.Period.Equal(expectedPeriod) {
		t.Fatalf("expected normalized period %s, got %s", expectedPeriod, created.Period)
	}
	assertAmountEqual(t, 275000, created.OriginalAmount)
	assertAmountEqual(t, 275000, created.OutstandingAmount)
	if created.Status != model.StudentObligationStatusPending {
		t.Fatalf("expected pending status, got %s", created.Status)
	}
	if created.Notes != "July arrears" {
		t.Fatalf("expected notes to be preserved")
	}
	if created.DueDate != dueDate {
		t.Fatalf("expected due date %s, got %s", dueDate, created.DueDate)
	}
	if created.IssuedAt.IsZero() {
		t.Fatalf("expected issued at to be initialized")
	}
}

func TestStudentObligationServiceCreateRejectsInvalidRequests(t *testing.T) {
	dueDate := time.Date(2026, time.July, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name              string
		req               *domain.StudentObligationCreate
		paymentProductFn  func() (*model.PaymentProduct, error)
		studentFn         func() (*model.Student, error)
		existsFn          func() (bool, error)
		expectedErr       error
		expectedCreateHit bool
	}{
		{
			name: "due date is required",
			req: &domain.StudentObligationCreate{
				StudentID:        studentID,
				PaymentProductID: productID,
			},
			expectedErr: shared.ErrStudentObligationDueDateRequired,
		},
		{
			name: "tenant scoped payment product must exist",
			req: &domain.StudentObligationCreate{
				StudentID:        studentID,
				PaymentProductID: productID,
				DueDate:          dueDate,
			},
			paymentProductFn: func() (*model.PaymentProduct, error) {
				return nil, shared.ErrPaymentProductNotFound
			},
			expectedErr: shared.ErrPaymentProductNotFound,
		},
		{
			name: "payment product must include payment policy",
			req: &domain.StudentObligationCreate{
				StudentID:        studentID,
				PaymentProductID: productID,
				DueDate:          dueDate,
			},
			paymentProductFn: func() (*model.PaymentProduct, error) {
				product := makePaymentProduct()
				product.PaymentPolicy = nil
				return product, nil
			},
			expectedErr: shared.ErrPaymentPolicyNotFound,
		},
		{
			name: "student must belong to tenant",
			req: &domain.StudentObligationCreate{
				StudentID:        studentID,
				PaymentProductID: productID,
				DueDate:          dueDate,
			},
			paymentProductFn: func() (*model.PaymentProduct, error) {
				return makePaymentProduct(), nil
			},
			studentFn: func() (*model.Student, error) {
				return makeStudent(studentID, otherTenantID), nil
			},
			expectedErr: shared.ErrStudentNotFound,
		},
		{
			name: "amount must be positive",
			req: &domain.StudentObligationCreate{
				StudentID:        studentID,
				PaymentProductID: productID,
				DueDate:          dueDate,
				Amount:           floatPtr(0),
			},
			paymentProductFn: func() (*model.PaymentProduct, error) {
				return makePaymentProduct(), nil
			},
			studentFn: func() (*model.Student, error) {
				return makeStudent(studentID, tenantID), nil
			},
			expectedErr: shared.ErrStudentObligationAmountInvalid,
		},
		{
			name: "duplicate active obligation for same period is rejected",
			req: &domain.StudentObligationCreate{
				StudentID:        studentID,
				PaymentProductID: productID,
				DueDate:          dueDate,
			},
			paymentProductFn: func() (*model.PaymentProduct, error) {
				return makePaymentProduct(), nil
			},
			studentFn: func() (*model.Student, error) {
				return makeStudent(studentID, tenantID), nil
			},
			existsFn: func() (bool, error) {
				return true, nil
			},
			expectedErr: shared.ErrStudentObligationAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createCalls := 0
			repo := &fakeStudentObligationRepository{
				createFn: func(ctx context.Context, studentObligation *model.StudentObligation) error {
					createCalls++
					return nil
				},
				existsActiveFn: func(ctx context.Context, studentID uuid.UUID, paymentProductID uuid.UUID, period time.Time) (bool, error) {
					if tt.existsFn == nil {
						return false, nil
					}
					return tt.existsFn()
				},
			}
			productService := &fakePaymentProductService{
				findByIDFn: func(
					ctx context.Context,
					authContext *security.AuthContext,
					id uuid.UUID,
					preloads ...model.PaymentProductPreload,
				) (*model.PaymentProduct, error) {
					if tt.paymentProductFn == nil {
						return makePaymentProduct(), nil
					}
					return tt.paymentProductFn()
				},
			}
			studentService := &fakeStudentService{
				firstByIDFn: func(ctx context.Context, id uuid.UUID) (*model.Student, error) {
					if tt.studentFn == nil {
						return makeStudent(id, tenantID), nil
					}
					return tt.studentFn()
				},
			}
			svc := service.NewStudentObligationService(
				&fakePaymentRepositoryManager{studentObligationRepo: repo},
				studentService,
				productService,
				&fakeParentStudentService{},
			)

			res, err := svc.Create(context.Background(), adminAuthContext(), tt.req)

			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
			}
			if res != nil {
				t.Fatalf("expected nil response on error")
			}
			if createCalls != 0 {
				t.Fatalf("expected repository create not to be called")
			}
		})
	}
}

func TestStudentObligationServiceCreateBulkDeduplicatesStudentsAndUsesTransaction(t *testing.T) {
	dueDate := time.Date(2026, time.August, 20, 8, 0, 0, 0, time.UTC)
	transactionCalls := 0
	createBatchCalls := 0
	var capturedBatch []model.StudentObligation
	lookupCount := 0
	repo := &fakeStudentObligationRepository{
		createBatchFn: func(ctx context.Context, studentObligations []model.StudentObligation) error {
			createBatchCalls++
			capturedBatch = append([]model.StudentObligation(nil), studentObligations...)
			return nil
		},
		existsActiveFn: func(ctx context.Context, studentID uuid.UUID, paymentProductID uuid.UUID, period time.Time) (bool, error) {
			expectedPeriod := time.Date(2026, time.August, 1, 0, 0, 0, 0, time.UTC)
			if !period.Equal(expectedPeriod) {
				t.Fatalf("expected normalized period %s, got %s", expectedPeriod, period)
			}
			return false, nil
		},
	}
	var manager *fakePaymentRepositoryManager
	manager = &fakePaymentRepositoryManager{
		studentObligationRepo: repo,
		transactionFn: func(ctx context.Context, fn func(repo domain.RepositoryManagerPayment) error) error {
			transactionCalls++
			return fn(manager)
		},
	}
	svc := service.NewStudentObligationService(
		manager,
		&fakeStudentService{
			firstByIDFn: func(ctx context.Context, id uuid.UUID) (*model.Student, error) {
				lookupCount++
				return makeStudent(id, tenantID), nil
			},
		},
		&fakePaymentProductService{
			findByIDFn: func(
				ctx context.Context,
				authContext *security.AuthContext,
				id uuid.UUID,
				preloads ...model.PaymentProductPreload,
			) (*model.PaymentProduct, error) {
				return makePaymentProduct(), nil
			},
		},
		&fakeParentStudentService{},
	)
	req := &domain.StudentObligationBulkCreate{
		PaymentProductID: productID,
		StudentIDs:       []uuid.UUID{studentID, secondStudentID, studentID},
		DueDate:          dueDate,
	}

	res, err := svc.CreateBulk(context.Background(), adminAuthContext(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if transactionCalls != 1 {
		t.Fatalf("expected one transaction, got %d", transactionCalls)
	}
	if createBatchCalls != 1 {
		t.Fatalf("expected one batch create, got %d", createBatchCalls)
	}
	if lookupCount != 2 {
		t.Fatalf("expected only unique student lookups, got %d", lookupCount)
	}
	if len(capturedBatch) != 2 || len(res) != 2 {
		t.Fatalf("expected two unique obligations to be created")
	}
	if capturedBatch[0].StudentID != studentID || capturedBatch[1].StudentID != secondStudentID {
		t.Fatalf("expected student order to preserve first unique appearance")
	}
	for _, obligation := range capturedBatch {
		assertAmountEqual(t, 500000, obligation.OriginalAmount)
		assertAmountEqual(t, 500000, obligation.OutstandingAmount)
		if obligation.Status != model.StudentObligationStatusPending {
			t.Fatalf("expected pending status, got %s", obligation.Status)
		}
		if obligation.Notes != "" {
			t.Fatalf("expected bulk create to initialize empty notes")
		}
	}
}

func TestStudentObligationServiceCreateBulkValidatesRequiredFields(t *testing.T) {
	svc := service.NewStudentObligationService(
		&fakePaymentRepositoryManager{
			studentObligationRepo: &fakeStudentObligationRepository{
				createBatchFn: func(ctx context.Context, studentObligations []model.StudentObligation) error {
					t.Fatalf("batch create should not be called when request is invalid")
					return nil
				},
			},
		},
		&fakeStudentService{},
		&fakePaymentProductService{},
		&fakeParentStudentService{},
	)

	tests := []struct {
		name        string
		req         *domain.StudentObligationBulkCreate
		expectedErr error
	}{
		{
			name: "student ids are required",
			req: &domain.StudentObligationBulkCreate{
				PaymentProductID: productID,
				DueDate:          time.Date(2026, time.August, 20, 0, 0, 0, 0, time.UTC),
			},
			expectedErr: shared.ErrStudentObligationStudentIDsEmpty,
		},
		{
			name: "due date is required",
			req: &domain.StudentObligationBulkCreate{
				PaymentProductID: productID,
				StudentIDs:       []uuid.UUID{studentID},
			},
			expectedErr: shared.ErrStudentObligationDueDateRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := svc.CreateBulk(context.Background(), adminAuthContext(), tt.req)
			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
			}
			if res != nil {
				t.Fatalf("expected nil result on error")
			}
		})
	}
}

func TestStudentObligationServiceFindOutstandingByStudentIDEnforcesAccessAndBuildsResponse(t *testing.T) {
	firstObligation := makeStudentObligation(studentID)
	firstObligation.OutstandingAmount = 300000
	secondObligation := makeStudentObligation(studentID)
	secondObligation.BaseModel.ID = uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	secondObligation.OriginalAmount = 200000
	secondObligation.OutstandingAmount = 150000
	secondObligation.Status = model.StudentObligationStatusPartial

	tests := []struct {
		name              string
		authContext       *security.AuthContext
		parentLinked      bool
		expectedErr       error
		expectedRepoCalls bool
	}{
		{
			name:              "admin can read student outstanding and response computes paid amount",
			authContext:       adminAuthContext(),
			expectedRepoCalls: true,
		},
		{
			name:              "linked parent can read student outstanding",
			authContext:       parentAuthContext(),
			parentLinked:      true,
			expectedRepoCalls: true,
		},
		{
			name:        "unlinked parent is rejected",
			authContext: parentAuthContext(),
			expectedErr: shared.ErrStudentNotFound,
		},
		{
			name:              "student can read own outstanding",
			authContext:       studentAuthContext(studentID),
			expectedRepoCalls: true,
		},
		{
			name:        "student cannot read another student outstanding",
			authContext: studentAuthContext(secondStudentID),
			expectedErr: shared.ErrStudentNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findCalls := 0
			sumCalls := 0
			repo := &fakeStudentObligationRepository{
				findOutstandingByStudentIDFn: func(ctx context.Context, id uuid.UUID) ([]model.StudentObligation, error) {
					findCalls++
					return []model.StudentObligation{*firstObligation, *secondObligation}, nil
				},
				sumOutstandingByStudentIDFn: func(ctx context.Context, id uuid.UUID) (float64, error) {
					sumCalls++
					return 450000, nil
				},
			}
			svc := service.NewStudentObligationService(
				&fakePaymentRepositoryManager{studentObligationRepo: repo},
				&fakeStudentService{
					firstByIDFn: func(ctx context.Context, id uuid.UUID) (*model.Student, error) {
						return makeStudent(id, tenantID), nil
					},
				},
				&fakePaymentProductService{},
				&fakeParentStudentService{
					existsFn: func(ctx context.Context, parentUserID uuid.UUID, studentID uuid.UUID) (bool, error) {
						return tt.parentLinked, nil
					},
				},
			)

			res, err := svc.FindOutstandingByStudentID(context.Background(), tt.authContext, studentID)

			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
			}
			if tt.expectedErr != nil {
				if res != nil {
					t.Fatalf("expected nil result on access failure")
				}
				if findCalls != 0 || sumCalls != 0 {
					t.Fatalf("expected repository outstanding lookups to be skipped on access failure")
				}
				return
			}
			if res == nil {
				t.Fatalf("expected outstanding response")
			}
			if res.StudentID != studentID {
				t.Fatalf("expected student ID %s, got %s", studentID, res.StudentID)
			}
			assertAmountEqual(t, 450000, res.TotalOutstanding)
			if len(res.Obligations) != 2 {
				t.Fatalf("expected two obligations in response")
			}
			assertAmountEqual(t, 200000, res.Obligations[0].PaidAmount)
			assertAmountEqual(t, 50000, res.Obligations[1].PaidAmount)
		})
	}
}

func TestStudentObligationServiceUpdateSupportsNoopAndMergedValidation(t *testing.T) {
	current := makeStudentObligation(studentID)
	firstByIDCalls := 0
	updateCalls := 0
	repo := &fakeStudentObligationRepository{
		firstByIDFn: func(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (*model.StudentObligation, error) {
			firstByIDCalls++
			return current, nil
		},
		updateFn: func(ctx context.Context, id uuid.UUID, req *domain.StudentObligationUpdate) (int, error) {
			updateCalls++
			return 1, nil
		},
	}
	svc := service.NewStudentObligationService(
		&fakePaymentRepositoryManager{studentObligationRepo: repo},
		&fakeStudentService{
			firstByIDFn: func(ctx context.Context, id uuid.UUID) (*model.Student, error) {
				return makeStudent(id, tenantID), nil
			},
		},
		&fakePaymentProductService{},
		&fakeParentStudentService{},
	)

	noopRes, err := svc.Update(context.Background(), adminAuthContext(), obligationID, &domain.StudentObligationUpdate{})
	if err != nil {
		t.Fatalf("expected no error for noop patch, got %v", err)
	}
	if noopRes != current {
		t.Fatalf("expected noop update to return current obligation")
	}
	if updateCalls != 0 {
		t.Fatalf("expected noop update not to call repository update")
	}

	current.DueDate = time.Time{}
	res, err := svc.Update(context.Background(), adminAuthContext(), obligationID, &domain.StudentObligationUpdate{
		Notes: stringPtr("rename only"),
	})
	if !errors.Is(err, shared.ErrStudentObligationDueDateRequired) {
		t.Fatalf("expected error %v, got %v", shared.ErrStudentObligationDueDateRequired, err)
	}
	if res != nil {
		t.Fatalf("expected nil response on merged validation failure")
	}
	if updateCalls != 0 {
		t.Fatalf("expected invalid merged state not to call repository update")
	}
	if firstByIDCalls != 2 {
		t.Fatalf("expected current obligation lookup for each update attempt, got %d", firstByIDCalls)
	}
}

func TestStudentObligationServiceUpdatePersistsPatchAndReloadsEntity(t *testing.T) {
	current := makeStudentObligation(studentID)
	updatedDueDate := time.Date(2026, time.August, 10, 0, 0, 0, 0, time.UTC)
	firstByIDCalls := 0
	var capturedReq *domain.StudentObligationUpdate
	repo := &fakeStudentObligationRepository{
		firstByIDFn: func(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (*model.StudentObligation, error) {
			firstByIDCalls++
			if firstByIDCalls == 1 {
				return current, nil
			}
			updated := *current
			updated.DueDate = updatedDueDate
			updated.Notes = "Updated note"
			return &updated, nil
		},
		updateFn: func(ctx context.Context, id uuid.UUID, req *domain.StudentObligationUpdate) (int, error) {
			capturedReq = req
			return 1, nil
		},
	}
	svc := service.NewStudentObligationService(
		&fakePaymentRepositoryManager{studentObligationRepo: repo},
		&fakeStudentService{
			firstByIDFn: func(ctx context.Context, id uuid.UUID) (*model.Student, error) {
				return makeStudent(id, tenantID), nil
			},
		},
		&fakePaymentProductService{},
		&fakeParentStudentService{},
	)

	res, err := svc.Update(context.Background(), adminAuthContext(), obligationID, &domain.StudentObligationUpdate{
		DueDate: timePtr(updatedDueDate),
		Notes:   stringPtr("Updated note"),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if capturedReq == nil || capturedReq.DueDate == nil || !capturedReq.DueDate.Equal(updatedDueDate) {
		t.Fatalf("expected due date patch to be passed to repository")
	}
	if capturedReq.Notes == nil || *capturedReq.Notes != "Updated note" {
		t.Fatalf("expected notes patch to be passed to repository")
	}
	if res == nil || !res.DueDate.Equal(updatedDueDate) || res.Notes != "Updated note" {
		t.Fatalf("expected updated obligation to be reloaded")
	}
	if firstByIDCalls != 2 {
		t.Fatalf("expected obligation reload after update, got %d lookup calls", firstByIDCalls)
	}
}

func TestStudentObligationServiceDeleteRejectsAllocatedAndDeletesUnallocated(t *testing.T) {
	tests := []struct {
		name        string
		allocated   bool
		expectedErr error
	}{
		{
			name:        "allocated obligation cannot be deleted",
			allocated:   true,
			expectedErr: shared.ErrStudentObligationAllocated,
		},
		{
			name:      "unallocated obligation can be deleted",
			allocated: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deleteCalls := 0
			repo := &fakeStudentObligationRepository{
				firstByIDFn: func(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (*model.StudentObligation, error) {
					return makeStudentObligation(studentID), nil
				},
				hasPaymentAllocationsFn: func(ctx context.Context, id uuid.UUID) (bool, error) {
					return tt.allocated, nil
				},
				deleteFn: func(ctx context.Context, id uuid.UUID) (int, error) {
					deleteCalls++
					return 1, nil
				},
			}
			svc := service.NewStudentObligationService(
				&fakePaymentRepositoryManager{studentObligationRepo: repo},
				&fakeStudentService{
					firstByIDFn: func(ctx context.Context, id uuid.UUID) (*model.Student, error) {
						return makeStudent(id, tenantID), nil
					},
				},
				&fakePaymentProductService{},
				&fakeParentStudentService{},
			)

			err := svc.Delete(context.Background(), adminAuthContext(), obligationID)

			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
			}
			if tt.expectedErr != nil {
				if deleteCalls != 0 {
					t.Fatalf("expected delete not to be called when obligation is allocated")
				}
				return
			}
			if deleteCalls != 1 {
				t.Fatalf("expected delete to be called once, got %d", deleteCalls)
			}
		})
	}
}
