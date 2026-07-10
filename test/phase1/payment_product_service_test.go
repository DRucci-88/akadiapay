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

type fakePaymentProductRepository struct {
	findByIDFn func(ctx context.Context, id uuid.UUID, tenantID uuid.UUID, preloads ...model.PaymentProductPreload) (*model.PaymentProduct, error)
	createFn   func(ctx context.Context, paymentProduct *model.PaymentProduct) error
	updateFn   func(ctx context.Context, id uuid.UUID, tenantID uuid.UUID, req *domain.PaymentProductUpdate) (int, error)
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
	if f.findByIDFn == nil {
		panic("unexpected call: FindByID")
	}
	return f.findByIDFn(ctx, id, tenantID, preloads...)
}

func (f *fakePaymentProductRepository) FindByIDsIncludingDeleted(
	ctx context.Context,
	ids []uuid.UUID,
) ([]model.PaymentProduct, error) {
	panic("unexpected call: FindByIDsIncludingDeleted")
}

func (f *fakePaymentProductRepository) Create(
	ctx context.Context,
	paymentProduct *model.PaymentProduct,
) error {
	if f.createFn == nil {
		panic("unexpected call: Create")
	}
	return f.createFn(ctx, paymentProduct)
}

func (f *fakePaymentProductRepository) Update(
	ctx context.Context,
	id uuid.UUID,
	tenantID uuid.UUID,
	req *domain.PaymentProductUpdate,
) (int, error) {
	if f.updateFn == nil {
		panic("unexpected call: Update")
	}
	return f.updateFn(ctx, id, tenantID, req)
}

func TestPaymentProductServiceCreate(t *testing.T) {
	tests := []struct {
		name          string
		req           *domain.PaymentProductCreate
		policyErr     error
		expectedErr   error
		assertCreated func(t *testing.T, created *model.PaymentProduct)
	}{
		{
			name: "create product succeeds with valid payment policy and keeps revenue account mapping",
			req: &domain.PaymentProductCreate{
				PaymentPolicyID:    policyID,
				Code:               "SMAN1_SPP_JUL_2026",
				Name:               "SPP July",
				Description:        "Monthly tuition",
				RevenueAccountCode: "4999",
				RevenueAccountName: "Custom Revenue",
				Price:              500000,
				Status:             model.PaymentProductStatusActive,
			},
			assertCreated: func(t *testing.T, created *model.PaymentProduct) {
				t.Helper()
				if created.TenantID != tenantID {
					t.Fatalf("expected tenant %s, got %s", tenantID, created.TenantID)
				}
				if created.RevenueAccountCode != "4999" {
					t.Fatalf("expected revenue code to be preserved")
				}
				if created.RevenueAccountName != "Custom Revenue" {
					t.Fatalf("expected revenue name to be preserved")
				}
			},
		},
		{
			name: "create product fails if price is zero",
			req: &domain.PaymentProductCreate{
				PaymentPolicyID: policyID,
				Code:            "ZERO_PRICE",
				Name:            "Zero Price",
				Price:           0,
			},
			expectedErr: shared.ErrPaymentProductPriceInvalid,
		},
		{
			name: "create product fails if payment policy not found for tenant",
			req: &domain.PaymentProductCreate{
				PaymentPolicyID: policyID,
				Code:            "MISSING_POLICY",
				Name:            "Missing Policy",
				Price:           10000,
			},
			policyErr:   shared.ErrPaymentPolicyNotFound,
			expectedErr: shared.ErrPaymentPolicyNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var created *model.PaymentProduct
			policyRepo := &fakePaymentPolicyRepository{
				firstByIDFn: func(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (*model.PaymentPolicy, error) {
					if tt.policyErr != nil {
						return nil, tt.policyErr
					}
					return makePaymentPolicy(), nil
				},
			}
			productRepo := &fakePaymentProductRepository{
				createFn: func(ctx context.Context, paymentProduct *model.PaymentProduct) error {
					created = paymentProduct
					return nil
				},
			}
			svc := service.NewPaymentProductService(&fakePaymentRepositoryManager{
				paymentPolicyRepo:  policyRepo,
				paymentProductRepo: productRepo,
			})

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
			if tt.assertCreated != nil {
				tt.assertCreated(t, created)
			}
		})
	}
}

func TestPaymentProductServiceUpdateSupportsPatchPointers(t *testing.T) {
	current := makePaymentProduct()
	firstByIDCalls := 0
	var capturedReq *domain.PaymentProductUpdate
	policyRepo := &fakePaymentPolicyRepository{
		firstByIDFn: func(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (*model.PaymentPolicy, error) {
			return makePaymentPolicy(), nil
		},
	}
	productRepo := &fakePaymentProductRepository{
		findByIDFn: func(ctx context.Context, id uuid.UUID, tenantID uuid.UUID, preloads ...model.PaymentProductPreload) (*model.PaymentProduct, error) {
			firstByIDCalls++
			if firstByIDCalls == 1 {
				return current, nil
			}

			updated := *current
			updated.Status = model.PaymentProductStatusInactive
			updated.RevenueAccountCode = "4998"
			updated.RevenueAccountName = "Updated Revenue"
			return &updated, nil
		},
		updateFn: func(ctx context.Context, id uuid.UUID, tenantID uuid.UUID, req *domain.PaymentProductUpdate) (int, error) {
			capturedReq = req
			return 1, nil
		},
	}
	svc := service.NewPaymentProductService(&fakePaymentRepositoryManager{
		paymentPolicyRepo:  policyRepo,
		paymentProductRepo: productRepo,
	})
	status := model.PaymentProductStatusInactive
	req := &domain.PaymentProductUpdate{
		Status:             &status,
		RevenueAccountCode: stringPtr("4998"),
		RevenueAccountName: stringPtr("Updated Revenue"),
	}

	res, err := svc.Update(context.Background(), adminAuthContext(), productID, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if capturedReq == nil || capturedReq.Status == nil || *capturedReq.Status != model.PaymentProductStatusInactive {
		t.Fatalf("expected inactive status to be forwarded in update request")
	}
	if capturedReq.RevenueAccountCode == nil || *capturedReq.RevenueAccountCode != "4998" {
		t.Fatalf("expected revenue account code patch to be preserved")
	}
	if capturedReq.RevenueAccountName == nil || *capturedReq.RevenueAccountName != "Updated Revenue" {
		t.Fatalf("expected revenue account name patch to be preserved")
	}
	if res == nil || res.Status != model.PaymentProductStatusInactive {
		t.Fatalf("expected updated product response")
	}
}

func TestPaymentProductServiceUpdateValidatesMergedState(t *testing.T) {
	current := makePaymentProduct()
	current.Price = 0
	productRepo := &fakePaymentProductRepository{
		findByIDFn: func(ctx context.Context, id uuid.UUID, tenantID uuid.UUID, preloads ...model.PaymentProductPreload) (*model.PaymentProduct, error) {
			return current, nil
		},
		updateFn: func(ctx context.Context, id uuid.UUID, tenantID uuid.UUID, req *domain.PaymentProductUpdate) (int, error) {
			t.Fatalf("update should not be called when merged state is invalid")
			return 0, nil
		},
	}
	svc := service.NewPaymentProductService(&fakePaymentRepositoryManager{
		paymentPolicyRepo:  &fakePaymentPolicyRepository{},
		paymentProductRepo: productRepo,
	})

	res, err := svc.Update(context.Background(), adminAuthContext(), productID, &domain.PaymentProductUpdate{
		Description: stringPtr("Rename only"),
	})

	if !errors.Is(err, shared.ErrPaymentProductPriceInvalid) {
		t.Fatalf("expected merged validation error %v, got %v", shared.ErrPaymentProductPriceInvalid, err)
	}
	if res != nil {
		t.Fatalf("expected nil response on merged validation error")
	}
}

func TestPaymentProductServiceFindByIDUsesTenantIsolation(t *testing.T) {
	productRepo := &fakePaymentProductRepository{
		findByIDFn: func(ctx context.Context, id uuid.UUID, tenantID uuid.UUID, preloads ...model.PaymentProductPreload) (*model.PaymentProduct, error) {
			if tenantID != adminAuthContext().TenantID {
				t.Fatalf("expected tenant isolation tenant %s, got %s", adminAuthContext().TenantID, tenantID)
			}
			return makePaymentProduct(), nil
		},
	}
	svc := service.NewPaymentProductService(&fakePaymentRepositoryManager{
		paymentPolicyRepo:  &fakePaymentPolicyRepository{},
		paymentProductRepo: productRepo,
	})

	res, err := svc.FindByID(context.Background(), adminAuthContext(), productID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res == nil || res.ID != productID {
		t.Fatalf("expected payment product response for requested ID")
	}
}
