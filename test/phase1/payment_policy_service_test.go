package phase1_test

import (
	"akadia/domain"
	service "akadia/internal/payment/service"
	"akadia/internal/platform/security"
	"akadia/internal/shared"
	"akadia/model"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type fakePaymentPolicyRepository struct {
	createFn    func(ctx context.Context, paymentPolicy *model.PaymentPolicy) error
	firstByIDFn func(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (*model.PaymentPolicy, error)
	updateFn    func(ctx context.Context, id uuid.UUID, tenantID uuid.UUID, req *domain.PaymentPolicyUpdate) (int, error)
}

func (f *fakePaymentPolicyRepository) Paginate(
	ctx context.Context,
	pageable *shared.Pageable,
	chain gorm.ChainInterface[model.PaymentPolicy],
) (*shared.Page[model.PaymentPolicy], error) {
	panic("unexpected call: Paginate")
}

func (f *fakePaymentPolicyRepository) FirstByID(
	ctx context.Context,
	id uuid.UUID,
	tenantID uuid.UUID,
) (*model.PaymentPolicy, error) {
	if f.firstByIDFn == nil {
		panic("unexpected call: FirstByID")
	}
	return f.firstByIDFn(ctx, id, tenantID)
}

func (f *fakePaymentPolicyRepository) FindPaginate(
	ctx context.Context,
	pageable *shared.Pageable,
	filter *domain.PaymentPolicyFilter,
	authContext *security.AuthContext,
) (*shared.Page[model.PaymentPolicy], error) {
	panic("unexpected call: FindPaginate")
}

func (f *fakePaymentPolicyRepository) Create(
	ctx context.Context,
	paymentPolicy *model.PaymentPolicy,
) error {
	if f.createFn == nil {
		panic("unexpected call: Create")
	}
	return f.createFn(ctx, paymentPolicy)
}

func (f *fakePaymentPolicyRepository) Update(
	ctx context.Context,
	id uuid.UUID,
	tenantID uuid.UUID,
	req *domain.PaymentPolicyUpdate,
) (int, error) {
	if f.updateFn == nil {
		panic("unexpected call: Update")
	}
	return f.updateFn(ctx, id, tenantID, req)
}

func TestPaymentPolicyServiceCreate(t *testing.T) {
	tests := []struct {
		name          string
		req           *domain.PaymentPolicyCreate
		expectedErr   error
		assertCreated func(t *testing.T, created *model.PaymentPolicy)
	}{
		{
			name: "full payment policy succeeds and zeros minimums when partial disabled",
			req: &domain.PaymentPolicyCreate{
				Code:                "FULL_PAYMENT",
				Name:                "Full Payment",
				Description:         "Must be paid in full",
				AllowPartial:        false,
				MinimumAmount:       100000,
				MinimumPercentage:   25,
				AllowOverPayment:    false,
				AutoCloseObligation: true,
			},
			assertCreated: func(t *testing.T, created *model.PaymentPolicy) {
				t.Helper()
				if created.TenantID != tenantID {
					t.Fatalf("expected tenant %s, got %s", tenantID, created.TenantID)
				}
				assertAmountEqual(t, 0, created.MinimumAmount)
				assertAmountEqual(t, 0, created.MinimumPercentage)
			},
		},
		{
			name: "partial payment policy succeeds when minimum amount valid",
			req: &domain.PaymentPolicyCreate{
				Code:                "PARTIAL_PAYMENT",
				Name:                "Partial Payment",
				AllowPartial:        true,
				MinimumAmount:       50000,
				AllowOverPayment:    false,
				AutoCloseObligation: true,
			},
			assertCreated: func(t *testing.T, created *model.PaymentPolicy) {
				t.Helper()
				assertAmountEqual(t, 50000, created.MinimumAmount)
				assertAmountEqual(t, 0, created.MinimumPercentage)
			},
		},
		{
			name: "partial payment policy fails when both minimums are zero",
			req: &domain.PaymentPolicyCreate{
				Code:         "PARTIAL_INVALID",
				Name:         "Partial Invalid",
				AllowPartial: true,
			},
			expectedErr: shared.ErrPaymentPolicyMinimumPaymentRequired,
		},
		{
			name: "policy fails when minimum amount is negative",
			req: &domain.PaymentPolicyCreate{
				Code:          "NEG_AMOUNT",
				Name:          "Negative Amount",
				AllowPartial:  true,
				MinimumAmount: -1,
			},
			expectedErr: shared.ErrPaymentPolicyMinimumAmountInvalid,
		},
		{
			name: "policy fails when minimum percentage is negative",
			req: &domain.PaymentPolicyCreate{
				Code:              "NEG_PERCENT",
				Name:              "Negative Percentage",
				AllowPartial:      true,
				MinimumPercentage: -1,
			},
			expectedErr: shared.ErrPaymentPolicyMinimumPercentageInvalid,
		},
		{
			name: "policy fails when minimum percentage is above 100",
			req: &domain.PaymentPolicyCreate{
				Code:              "HIGH_PERCENT",
				Name:              "High Percentage",
				AllowPartial:      true,
				MinimumPercentage: 101,
			},
			expectedErr: shared.ErrPaymentPolicyMinimumPercentageInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var created *model.PaymentPolicy
			repo := &fakePaymentPolicyRepository{
				createFn: func(ctx context.Context, paymentPolicy *model.PaymentPolicy) error {
					created = paymentPolicy
					return nil
				},
			}
			repositoryManager := &fakePaymentRepositoryManager{
				paymentPolicyRepo: repo,
			}

			svc := service.NewPaymentPolicyService(repositoryManager)
			res, err := svc.Create(context.Background(), adminAuthContext(), tt.req)

			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
			}
			if tt.expectedErr != nil {
				if created != nil {
					t.Fatalf("expected create not to be called")
				}
				if res != nil {
					t.Fatalf("expected nil response on error")
				}
				return
			}

			if created == nil {
				t.Fatalf("expected create to be called")
			}
			if res != created {
				t.Fatalf("expected response to be created model instance")
			}
			if tt.assertCreated != nil {
				tt.assertCreated(t, created)
			}
		})
	}
}

func TestPaymentPolicyServiceUpdatePreservesRepositoryError(t *testing.T) {
	firstByIDCalls := 0
	updateErr := errors.New("update failed")
	current := makePaymentPolicy()
	repo := &fakePaymentPolicyRepository{
		firstByIDFn: func(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (*model.PaymentPolicy, error) {
			firstByIDCalls++
			return current, nil
		},
		updateFn: func(ctx context.Context, id uuid.UUID, tenantID uuid.UUID, req *domain.PaymentPolicyUpdate) (int, error) {
			return 0, updateErr
		},
	}
	svc := service.NewPaymentPolicyService(&fakePaymentRepositoryManager{
		paymentPolicyRepo: repo,
	})

	res, err := svc.Update(context.Background(), adminAuthContext(), policyID, &domain.PaymentPolicyUpdate{
		Name: stringPtr("Updated Name"),
	})

	if !errors.Is(err, updateErr) {
		t.Fatalf("expected update error %v, got %v", updateErr, err)
	}
	if res != nil {
		t.Fatalf("expected nil response on repository error")
	}
	if firstByIDCalls != 1 {
		t.Fatalf("expected one FirstByID call before update failure, got %d", firstByIDCalls)
	}
}

func TestPaymentPolicyServiceUpdateSupportsFalsePointers(t *testing.T) {
	current := makePaymentPolicy()
	firstByIDCalls := 0
	var capturedReq *domain.PaymentPolicyUpdate
	repo := &fakePaymentPolicyRepository{
		firstByIDFn: func(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (*model.PaymentPolicy, error) {
			firstByIDCalls++
			if firstByIDCalls == 1 {
				return current, nil
			}

			updated := *current
			updated.AllowPartial = false
			updated.MinimumAmount = 0
			updated.MinimumPercentage = 0
			return &updated, nil
		},
		updateFn: func(ctx context.Context, id uuid.UUID, tenantID uuid.UUID, req *domain.PaymentPolicyUpdate) (int, error) {
			capturedReq = req
			return 1, nil
		},
	}
	svc := service.NewPaymentPolicyService(&fakePaymentRepositoryManager{
		paymentPolicyRepo: repo,
	})

	res, err := svc.Update(context.Background(), adminAuthContext(), policyID, &domain.PaymentPolicyUpdate{
		AllowPartial: boolPtr(false),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if capturedReq == nil || capturedReq.AllowPartial == nil || *capturedReq.AllowPartial {
		t.Fatalf("expected AllowPartial=false to be preserved in update request")
	}
	if capturedReq.MinimumAmount == nil || capturedReq.MinimumPercentage == nil {
		t.Fatalf("expected minimum values to be explicitly zeroed when partial is disabled")
	}
	assertAmountEqual(t, 0, *capturedReq.MinimumAmount)
	assertAmountEqual(t, 0, *capturedReq.MinimumPercentage)
	if res == nil || res.AllowPartial {
		t.Fatalf("expected updated response with AllowPartial=false")
	}
}

func TestPaymentPolicyServiceUpdateValidatesMergedState(t *testing.T) {
	current := makePaymentPolicy()
	current.MinimumAmount = 0
	current.MinimumPercentage = 0
	repo := &fakePaymentPolicyRepository{
		firstByIDFn: func(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (*model.PaymentPolicy, error) {
			return current, nil
		},
		updateFn: func(ctx context.Context, id uuid.UUID, tenantID uuid.UUID, req *domain.PaymentPolicyUpdate) (int, error) {
			t.Fatalf("update should not be called when merged state is invalid")
			return 0, nil
		},
	}
	svc := service.NewPaymentPolicyService(&fakePaymentRepositoryManager{
		paymentPolicyRepo: repo,
	})

	res, err := svc.Update(context.Background(), adminAuthContext(), policyID, &domain.PaymentPolicyUpdate{
		Name: stringPtr("Rename Only"),
	})

	if !errors.Is(err, shared.ErrPaymentPolicyMinimumPaymentRequired) {
		t.Fatalf("expected merged validation error %v, got %v", shared.ErrPaymentPolicyMinimumPaymentRequired, err)
	}
	if res != nil {
		t.Fatalf("expected nil response on merged validation error")
	}
}
