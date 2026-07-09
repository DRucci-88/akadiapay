package repository

import (
	"akadia/domain"
	"akadia/internal/platform/security"
	"akadia/internal/shared"
	"akadia/model"
	"akadia/model/generated"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type paymentPolicyRepositoryImpl struct {
	db    *gorm.DB
	query gorm.Interface[model.PaymentPolicy]
}

func (r *RepositoryManagerPaymentImpl) PaymentPolicy() domain.PaymentPolicyRepository {
	return &paymentPolicyRepositoryImpl{
		db:    r.db,
		query: gorm.G[model.PaymentPolicy](r.db),
	}
}

func (r *paymentPolicyRepositoryImpl) Paginate(
	ctx context.Context,
	pageable *shared.Pageable,
	chain gorm.ChainInterface[model.PaymentPolicy],
) (*shared.Page[model.PaymentPolicy], error) {
	// pageable.Normalize()

	total, err := chain.Count(ctx, "*")
	if err != nil {
		return shared.NewPageEmpty[model.PaymentPolicy](pageable), err
	}
	items, err := chain.
		Offset(pageable.Offset()).
		Limit(pageable.Limit()).
		Find(ctx)

	if err != nil {
		return shared.NewPageEmpty[model.PaymentPolicy](pageable), err
	}
	page := shared.NewPage(pageable, total, items)
	return page, nil
}

func (r *paymentPolicyRepositoryImpl) FirstByID(
	ctx context.Context,
	id uuid.UUID,
) (*model.PaymentPolicy, error) {
	paymentPolicy, err := r.query.
		Where(generated.BaseModel.ID.Eq(id)).
		First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, shared.ErrPaymentPolicyNotFound
	}
	return &paymentPolicy, err
}

func (r *paymentPolicyRepositoryImpl) FindPaginate(
	ctx context.Context,
	pageable *shared.Pageable,
	filter *domain.PaymentPolicyFilter,
	authContext *security.AuthContext,
) (*shared.Page[model.PaymentPolicy], error) {
	chain := r.query.Scopes()

	if filter == nil {
		return r.Paginate(ctx, pageable, chain)
	}

	chain = chain.Where(generated.PaymentPolicy.TenantID.Eq(authContext.TenantID))

	if filter.Keyword != nil {
		chain = chain.
			Where(generated.PaymentPolicy.Code.ILike(*filter.Keyword)).
			Or(generated.PaymentPolicy.Name.ILike(*filter.Keyword))
	}
	if filter.AllowPartial != nil {
		chain = chain.
			Where(generated.PaymentPolicy.AllowPartial.Eq(*filter.AllowPartial))
	}
	if filter.AllowOverPayment != nil {
		chain = chain.
			Where(generated.PaymentPolicy.AllowOverPayment.Eq(*filter.AllowOverPayment))
	}
	if filter.AutoCloseObligation != nil {
		chain = chain.
			Where(generated.PaymentPolicy.AutoCloseObligation.Eq(*filter.AutoCloseObligation))
	}

	return r.Paginate(ctx, pageable, chain)
}

func (r *paymentPolicyRepositoryImpl) Create(
	ctx context.Context,
	paymentPolicy *model.PaymentPolicy,
) error {
	return r.query.Create(ctx, paymentPolicy)
}

func (r *paymentPolicyRepositoryImpl) Update(
	ctx context.Context,
	id uuid.UUID,
	tenantID uuid.UUID,
	paymentPolicy *model.PaymentPolicy,
) (int, error) {
	rows, err := r.query.
		Where(generated.BaseModel.ID.Eq(id)).
		Where(generated.PaymentPolicy.TenantID.Eq(tenantID)).
		Updates(ctx, *paymentPolicy)
	if rows == 0 {
		return rows, shared.ErrPaymentPolicyNotFound
	}
	return rows, err
}
