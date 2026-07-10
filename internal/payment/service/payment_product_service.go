package service

import (
	"akadia/domain"
	"akadia/internal/platform/security"
	"akadia/internal/shared"
	"akadia/model"
	"context"

	"github.com/google/uuid"
)

type paymentProductServiceImpl struct {
	repo               domain.RepositoryManagerPayment
	paymentProductRepo domain.PaymentProductRepository
	paymentPolicyRepo  domain.PaymentPolicyRepository
}

func NewPaymentProductService(repo domain.RepositoryManagerPayment) domain.PaymentProductService {
	return &paymentProductServiceImpl{
		repo:               repo,
		paymentProductRepo: repo.PaymentProduct(),
		paymentPolicyRepo:  repo.PaymentPolicy(),
	}
}

func (s *paymentProductServiceImpl) FindPaginate(
	ctx context.Context,
	pageable *shared.Pageable,
	filter *domain.PaymentProductFilter,
	authContext *security.AuthContext,
) (*shared.Page[domain.PaymentProductResponse], error) {
	page, err := s.paymentProductRepo.FindPaginate(ctx, pageable, filter, authContext)
	if err != nil {
		return shared.NewPageEmpty[domain.PaymentProductResponse](pageable), err
	}

	return &shared.Page[domain.PaymentProductResponse]{
		Data:       domain.NewPaymentProductResponses(page.Data),
		Pagination: page.Pagination,
	}, nil
}

func (s *paymentProductServiceImpl) FindByID(
	ctx context.Context,
	authContext *security.AuthContext,
	id uuid.UUID,
	preloads ...model.PaymentProductPreload,
) (*model.PaymentProduct, error) {
	return s.paymentProductRepo.FindByID(ctx, id, authContext.TenantID, preloads...)
}

func (s *paymentProductServiceImpl) Create(
	ctx context.Context,
	authContext *security.AuthContext,
	req *domain.PaymentProductCreate,
) (*model.PaymentProduct, error) {
	if !shared.FloatGreater(req.Price, 0) {
		return nil, shared.ErrPaymentProductPriceInvalid
	}

	status := req.Status
	if status == "" {
		status = model.PaymentProductStatusActive
	}
	if !isValidPaymentProductStatus(status) {
		return nil, shared.ErrPaymentProductStatusInvalid
	}

	paymentPolicy, err := s.paymentPolicyRepo.FirstByID(ctx, req.PaymentPolicyID, authContext.TenantID)
	if err != nil {
		return nil, err
	}
	_ = paymentPolicy

	paymentProduct := &model.PaymentProduct{
		TenantID:           authContext.TenantID,
		PaymentPolicyID:    req.PaymentPolicyID,
		Code:               req.Code,
		Name:               req.Name,
		Description:        req.Description,
		RevenueAccountCode: req.RevenueAccountCode,
		RevenueAccountName: req.RevenueAccountName,
		Price:              req.Price,
		Status:             status,
	}
	normalizePaymentProductRevenueAccount(paymentProduct)

	if err := s.paymentProductRepo.Create(ctx, paymentProduct); err != nil {
		return nil, err
	}

	return paymentProduct, nil
}

func (s *paymentProductServiceImpl) Update(
	ctx context.Context,
	authContext *security.AuthContext,
	id uuid.UUID,
	req *domain.PaymentProductUpdate,
) (*model.PaymentProduct, error) {
	paymentProduct, err := s.FindByID(ctx, authContext, id)
	if err != nil {
		return nil, err
	}

	merged := *paymentProduct
	if req.PaymentPolicyID != nil {
		paymentPolicy, err := s.paymentPolicyRepo.FirstByID(ctx, *req.PaymentPolicyID, authContext.TenantID)
		if err != nil {
			return nil, err
		}
		_ = paymentPolicy
		merged.PaymentPolicyID = *req.PaymentPolicyID
	}
	if req.Code != nil {
		merged.Code = *req.Code
	}
	if req.Name != nil {
		merged.Name = *req.Name
	}
	if req.Description != nil {
		merged.Description = *req.Description
	}
	if req.RevenueAccountCode != nil {
		merged.RevenueAccountCode = *req.RevenueAccountCode
	}
	if req.RevenueAccountName != nil {
		merged.RevenueAccountName = *req.RevenueAccountName
	}
	if req.Price != nil {
		merged.Price = *req.Price
	}
	if req.Status != nil {
		merged.Status = *req.Status
	}
	normalizePaymentProductRevenueAccount(&merged)

	if !shared.FloatGreater(merged.Price, 0) {
		return nil, shared.ErrPaymentProductPriceInvalid
	}
	if !isValidPaymentProductStatus(merged.Status) {
		return nil, shared.ErrPaymentProductStatusInvalid
	}

	req.RevenueAccountCode = &merged.RevenueAccountCode
	req.RevenueAccountName = &merged.RevenueAccountName

	_, err = s.paymentProductRepo.Update(
		ctx,
		id,
		authContext.TenantID,
		req,
	)
	if err != nil {
		return nil, err
	}

	return s.FindByID(ctx, authContext, id)
}

func isValidPaymentProductStatus(status model.PaymentProductStatus) bool {
	switch status {
	case model.PaymentProductStatusActive, model.PaymentProductStatusInactive:
		return true
	default:
		return false
	}
}
