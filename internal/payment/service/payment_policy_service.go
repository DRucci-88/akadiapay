package service

import (
	"akadia/domain"
	"akadia/internal/platform/security"
	"akadia/internal/shared"
	"akadia/model"
	"context"

	"github.com/google/uuid"
)

type paymentPolicyServiceImpl struct {
	repo              domain.RepositoryManagerPayment
	paymentPolicyRepo domain.PaymentPolicyRepository
}

func NewPaymentPolicyService(repo domain.RepositoryManagerPayment) domain.PaymentPolicyService {
	return &paymentPolicyServiceImpl{
		repo:              repo,
		paymentPolicyRepo: repo.PaymentPolicy(),
	}
}

func (s *paymentPolicyServiceImpl) FirstByID(
	ctx context.Context,
	id uuid.UUID,
) (*model.PaymentPolicy, error) {
	return s.paymentPolicyRepo.FirstByID(ctx, id)
}

func (s *paymentPolicyServiceImpl) FindPaginate(
	ctx context.Context,
	pageable *shared.Pageable,
	filter *domain.PaymentPolicyFilter,
	authContext *security.AuthContext,
) (*shared.Page[domain.PaymentPolicyResponse], error) {
	page, err := s.paymentPolicyRepo.FindPaginate(ctx, pageable, filter, authContext)
	if err != nil {
		return shared.NewPageEmpty[domain.PaymentPolicyResponse](pageable), err
	}

	return &shared.Page[domain.PaymentPolicyResponse]{
		Data:       domain.NewPaymentPolicyResponses(page.Data),
		Pagination: page.Pagination,
	}, nil
}

func (s *paymentPolicyServiceImpl) Create(
	ctx context.Context,
	authContext *security.AuthContext,
	req *domain.PaymentPolicyCreate,
) (*model.PaymentPolicy, error) {
	if !req.AllowPartial {
		req.MinimumAmount = 0
		req.MinimumPercentage = 0
	}

	if req.MinimumAmount < 0 {
		return nil, shared.ErrPaymentPolicyMinimumAmountInvalid
	}

	if req.MinimumPercentage < 0 || req.MinimumPercentage > 100 {
		return nil, shared.ErrPaymentPolicyMinimumPercentageInvalid
	}

	if req.AllowPartial &&
		req.MinimumAmount == 0 &&
		req.MinimumPercentage == 0 {
		return nil, shared.ErrPaymentPolicyMinimumPaymentRequired
	}

	paymentPolicy := &model.PaymentPolicy{
		TenantID:            authContext.TenantID,
		Code:                req.Code,
		Name:                req.Name,
		Description:         req.Description,
		AllowPartial:        req.AllowPartial,
		MinimumAmount:       req.MinimumAmount,
		MinimumPercentage:   req.MinimumPercentage,
		AllowOverPayment:    req.AllowOverPayment,
		AutoCloseObligation: req.AutoCloseObligation,
	}

	if err := s.paymentPolicyRepo.Create(ctx, paymentPolicy); err != nil {
		return nil, err
	}
	return paymentPolicy, nil
}

// FULL PAYMENT
// {
//   "code": "FULL_PAYMENT",
//   "name": "Full Payment",
//   "description": "Must be paid in full",
//   "allow_partial": false,
//   "allow_over_payment": false,
//   "auto_close_obligation": true
// }

// PARTIAL PAYMENT
// {
//   "code": "PARTIAL_PAYMENT",
//   "name": "Partial Payment",
//   "description": "Can be paid gradually",
//   "allow_partial": true,
//   "minimum_amount": 100000,
//   "minimum_percentage": 20,
//   "allow_over_payment": false,
//   "auto_close_obligation": true
// }

func (s *paymentPolicyServiceImpl) Update(
	ctx context.Context,
	authContext *security.AuthContext,
	id uuid.UUID,
	req *domain.PaymentPolicyUpdate,
) (*model.PaymentPolicy, error) {
	_, err := s.paymentPolicyRepo.Update(
		ctx,
		id,
		authContext.TenantID,
		req,
	)

	paymentPolicy, err := s.FirstByID(ctx, id)

	return paymentPolicy, err
}
