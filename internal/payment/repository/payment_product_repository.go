package repository

import (
	"akadia/domain"
	"akadia/internal/shared"
	"akadia/model"
	"akadia/model/generated"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type paymentProductRepositoryImpl struct {
	db    *gorm.DB
	query gorm.Interface[model.PaymentProduct]
}

func (r *RepositoryManagerPaymentImpl) PaymentProduct() domain.PaymentProductRepository {
	return &paymentProductRepositoryImpl{
		db:    r.db,
		query: gorm.G[model.PaymentProduct](r.db),
	}
}

func (r *paymentProductRepositoryImpl) QueryWithPreloads(
	preloads ...model.PaymentProductPreload,
) gorm.ChainInterface[model.PaymentProduct] {
	var chain gorm.ChainInterface[model.PaymentProduct] = r.query.Scopes()
	for _, preload := range preloads {
		chain = chain.Preload(string(preload), nil)
	}
	return chain
}

func (r *paymentProductRepositoryImpl) FindByID(
	ctx context.Context,
	id uuid.UUID,
	preloads ...model.PaymentProductPreload,
) (*model.PaymentProduct, error) {
	paymentProduct, err := r.QueryWithPreloads(preloads...).
		Where(generated.BaseModel.ID.Eq(id)).
		Where(generated.PaymentProduct.Status.Eq(string(model.PaymentProductStatusActive))).
		First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, shared.ErrPaymentProductNotFound
	}

	return &paymentProduct, err
}
