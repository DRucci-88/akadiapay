package service

import (
	"akadia/domain"
	"akadia/model"
	"context"

	"github.com/google/uuid"
)

type paymentProductServiceImpl struct {
	repo               domain.RepositoryManagerPayment
	paymentProductRepo domain.PaymentProductRepository
}

func NewPaymentProductService(repo domain.RepositoryManagerPayment) domain.PaymentProductService {
	return &paymentProductServiceImpl{
		repo:               repo,
		paymentProductRepo: repo.PaymentProduct(),
	}
}

func (s paymentProductServiceImpl) FindByID(
	ctx context.Context,
	id uuid.UUID,
	preloads ...model.PaymentProductPreload,
) (*model.PaymentProduct, error) {
	return s.paymentProductRepo.FindByID(ctx, id, preloads...)
}
