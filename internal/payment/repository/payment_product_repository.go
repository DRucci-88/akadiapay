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

func (r *paymentProductRepositoryImpl) Paginate(
	ctx context.Context,
	pageable *shared.Pageable,
	chain gorm.ChainInterface[model.PaymentProduct],
) (*shared.Page[model.PaymentProduct], error) {
	total, err := chain.Count(ctx, "*")
	if err != nil {
		return shared.NewPageEmpty[model.PaymentProduct](pageable), err
	}

	items, err := chain.
		Offset(pageable.Offset()).
		Limit(pageable.Limit()).
		Find(ctx)
	if err != nil {
		return shared.NewPageEmpty[model.PaymentProduct](pageable), err
	}

	return shared.NewPage(pageable, total, items), nil
}

func (r *paymentProductRepositoryImpl) FindPaginate(
	ctx context.Context,
	pageable *shared.Pageable,
	filter *domain.PaymentProductFilter,
	authContext *security.AuthContext,
) (*shared.Page[model.PaymentProduct], error) {
	chain := r.query.
		Scopes().
		Where(generated.PaymentProduct.TenantID.Eq(authContext.TenantID))

	if filter == nil {
		return r.Paginate(ctx, pageable, chain)
	}

	if filter.Keyword != nil {
		chain = chain.
			Where("(code ILIKE ? OR name ILIKE ?)", *filter.Keyword, *filter.Keyword)
	}
	if filter.PaymentPolicyID != nil {
		chain = chain.
			Where(generated.PaymentProduct.PaymentPolicyID.Eq(*filter.PaymentPolicyID))
	}
	if filter.Status != nil {
		chain = chain.
			Where(generated.PaymentProduct.Status.Eq(string(*filter.Status)))
	}

	return r.Paginate(ctx, pageable, chain)
}

func (r *paymentProductRepositoryImpl) FindByID(
	ctx context.Context,
	id uuid.UUID,
	preloads ...model.PaymentProductPreload,
) (*model.PaymentProduct, error) {
	paymentProduct, err := r.QueryWithPreloads(preloads...).
		Where(generated.BaseModel.ID.Eq(id)).
		First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, shared.ErrPaymentProductNotFound
	}

	return &paymentProduct, err
}

func (r *paymentProductRepositoryImpl) Create(
	ctx context.Context,
	paymentProduct *model.PaymentProduct,
) error {
	return r.query.Create(ctx, paymentProduct)
}

func (r *paymentProductRepositoryImpl) Update(
	ctx context.Context,
	id uuid.UUID,
	tenantID uuid.UUID,
	req *domain.PaymentProductUpdate,
) (int, error) {
	updates := shared.UpdateMap{}

	updates.SetIfNotNil(generated.PaymentProduct.PaymentPolicyID.Column().Name, req.PaymentPolicyID)
	updates.SetIfNotNil(generated.PaymentProduct.Code.Column().Name, req.Code)
	updates.SetIfNotNil(generated.PaymentProduct.Name.Column().Name, req.Name)
	updates.SetIfNotNil(generated.PaymentProduct.Description.Column().Name, req.Description)
	updates.SetIfNotNil(generated.PaymentProduct.Price.Column().Name, req.Price)
	updates.SetIfNotNil(generated.PaymentProduct.Status.Column().Name, req.Status)

	rows, err := gorm.G[map[string]any](r.db).
		Where(generated.BaseModel.ID.Eq(id)).
		Where(generated.PaymentProduct.TenantID.Eq(tenantID)).
		Select("*").
		Omit("id", "created_at").
		Updates(ctx, updates)
	if rows == 0 {
		return rows, shared.ErrPaymentProductNotFound
	}

	return rows, err
}
