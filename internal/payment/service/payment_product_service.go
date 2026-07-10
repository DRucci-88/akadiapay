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
	paymentProduct, err := s.paymentProductRepo.FindByID(ctx, id, preloads...)
	if err != nil {
		return nil, err
	}
	if paymentProduct.TenantID != authContext.TenantID {
		return nil, shared.ErrPaymentProductNotFound
	}

	return paymentProduct, nil
}

func (s *paymentProductServiceImpl) Create(
	ctx context.Context,
	authContext *security.AuthContext,
	req *domain.PaymentProductCreate,
) (*model.PaymentProduct, error) {
	if req.Price < 0 {
		return nil, shared.ErrPaymentProductPriceInvalid
	}

	status := req.Status
	if status == "" {
		status = model.PaymentProductStatusActive
	}
	if !isValidPaymentProductStatus(status) {
		return nil, shared.ErrPaymentProductStatusInvalid
	}

	paymentPolicy, err := s.paymentPolicyRepo.FirstByID(ctx, req.PaymentPolicyID)
	if err != nil {
		return nil, err
	}
	if paymentPolicy.TenantID != authContext.TenantID {
		return nil, shared.ErrPaymentPolicyNotFound
	}

	paymentProduct := &model.PaymentProduct{
		TenantID:        authContext.TenantID,
		PaymentPolicyID: req.PaymentPolicyID,
		Code:            req.Code,
		Name:            req.Name,
		Description:     req.Description,
		Price:           req.Price,
		Status:          status,
	}

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
	if req.PaymentPolicyID != nil {
		paymentPolicy, err := s.paymentPolicyRepo.FirstByID(ctx, *req.PaymentPolicyID)
		if err != nil {
			return nil, err
		}
		if paymentPolicy.TenantID != authContext.TenantID {
			return nil, shared.ErrPaymentPolicyNotFound
		}
	}
	if req.Price != nil {
		if *req.Price < 0 {
			return nil, shared.ErrPaymentProductPriceInvalid
		}
	}
	if req.Status != nil {
		if !isValidPaymentProductStatus(*req.Status) {
			return nil, shared.ErrPaymentProductStatusInvalid
		}
	}

	_, err := s.paymentProductRepo.Update(
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
