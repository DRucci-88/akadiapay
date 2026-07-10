package repository

import (
	"akadia/domain"
	"akadia/model"
	"akadia/model/generated"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type paymentAllocationRepositoryImpl struct {
	db    *gorm.DB
	query gorm.Interface[model.PaymentAllocation]
}

func (r *RepositoryManagerPaymentImpl) PaymentAllocation() domain.PaymentAllocationRepository {
	return &paymentAllocationRepositoryImpl{
		db:    r.db,
		query: gorm.G[model.PaymentAllocation](r.db),
	}
}

func (r *paymentAllocationRepositoryImpl) CreateBatch(
	ctx context.Context,
	paymentAllocations []model.PaymentAllocation,
) error {
	return r.db.
		WithContext(ctx).
		CreateInBatches(&paymentAllocations, 100).
		Error
}

func (r *paymentAllocationRepositoryImpl) FindByPaymentOrderID(
	ctx context.Context,
	paymentOrderID uuid.UUID,
) ([]model.PaymentAllocation, error) {
	return r.query.
		Where(generated.PaymentAllocation.PaymentOrderID.Eq(paymentOrderID)).
		Find(ctx)
}
