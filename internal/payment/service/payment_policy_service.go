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
	authContext *security.AuthContext,
	id uuid.UUID,
) (*model.PaymentPolicy, error) {
	return s.paymentPolicyRepo.FirstByID(ctx, id, authContext.TenantID)
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
	normalizePaymentPolicyForValidation(paymentPolicy)
	if err := validatePaymentPolicyState(paymentPolicy); err != nil {
		return nil, err
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
	paymentPolicy, err := s.FirstByID(ctx, authContext, id)
	if err != nil {
		return nil, err
	}

	merged := *paymentPolicy
	if req.Code != nil {
		merged.Code = *req.Code
	}
	if req.Name != nil {
		merged.Name = *req.Name
	}
	if req.Description != nil {
		merged.Description = *req.Description
	}
	if req.AllowPartial != nil {
		merged.AllowPartial = *req.AllowPartial
	}
	if req.MinimumAmount != nil {
		merged.MinimumAmount = *req.MinimumAmount
	}
	if req.MinimumPercentage != nil {
		merged.MinimumPercentage = *req.MinimumPercentage
	}
	if req.AllowOverPayment != nil {
		merged.AllowOverPayment = *req.AllowOverPayment
	}
	if req.AutoCloseObligation != nil {
		merged.AutoCloseObligation = *req.AutoCloseObligation
	}

	normalizePaymentPolicyForValidation(&merged)
	if err := validatePaymentPolicyState(&merged); err != nil {
		return nil, err
	}

	if !merged.AllowPartial {
		zero := 0.0
		req.MinimumAmount = &zero
		req.MinimumPercentage = &zero
	}

	_, err = s.paymentPolicyRepo.Update(
		ctx,
		id,
		authContext.TenantID,
		req,
	)
	if err != nil {
		return nil, err
	}

	paymentPolicy, err = s.FirstByID(ctx, authContext, id)

	return paymentPolicy, err
}

func normalizePaymentPolicyForValidation(
	paymentPolicy *model.PaymentPolicy,
) {
	if !paymentPolicy.AllowPartial {
		paymentPolicy.MinimumAmount = 0
		paymentPolicy.MinimumPercentage = 0
	}
}

func validatePaymentPolicyState(
	paymentPolicy *model.PaymentPolicy,
) error {
	if shared.FloatLess(paymentPolicy.MinimumAmount, 0) {
		return shared.ErrPaymentPolicyMinimumAmountInvalid
	}

	if shared.FloatLess(paymentPolicy.MinimumPercentage, 0) ||
		shared.FloatGreater(paymentPolicy.MinimumPercentage, 100) {
		return shared.ErrPaymentPolicyMinimumPercentageInvalid
	}

	if paymentPolicy.AllowPartial &&
		shared.FloatIsZero(paymentPolicy.MinimumAmount) &&
		shared.FloatIsZero(paymentPolicy.MinimumPercentage) {
		return shared.ErrPaymentPolicyMinimumPaymentRequired
	}

	return nil
}
