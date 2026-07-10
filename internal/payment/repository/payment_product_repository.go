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

var paymentProductSortableColumns = map[string]string{
	"created_at":           generated.BaseModel.CreatedAt.Column().Name,
	"updated_at":           generated.BaseModel.UpdatedAt.Column().Name,
	"code":                 generated.PaymentProduct.Code.Column().Name,
	"name":                 generated.PaymentProduct.Name.Column().Name,
	"price":                generated.PaymentProduct.Price.Column().Name,
	"status":               generated.PaymentProduct.Status.Column().Name,
	"payment_policy_id":    generated.PaymentProduct.PaymentPolicyID.Column().Name,
	"revenue_account_code": generated.PaymentProduct.RevenueAccountCode.Column().Name,
	"revenue_account_name": generated.PaymentProduct.RevenueAccountName.Column().Name,
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

	if filter != nil && filter.Keyword != nil {
		chain = chain.
			Where("(code ILIKE ? OR name ILIKE ?)", *filter.Keyword, *filter.Keyword)
	}
	if filter != nil && filter.PaymentPolicyID != nil {
		chain = chain.
			Where(generated.PaymentProduct.PaymentPolicyID.Eq(*filter.PaymentPolicyID))
	}
	if filter != nil && filter.Status != nil {
		chain = chain.
			Where(generated.PaymentProduct.Status.Eq(string(*filter.Status)))
	}

	var sortBy *string
	var order *string
	if filter != nil {
		sortBy = filter.SortBy
		order = filter.Order
	}
	chain = chain.Order(
		shared.ResolveSortClause(
			sortBy,
			order,
			paymentProductSortableColumns,
			generated.BaseModel.CreatedAt.Column().Name,
			"DESC",
		),
	)

	return r.Paginate(ctx, pageable, chain)
}

func (r *paymentProductRepositoryImpl) FindByID(
	ctx context.Context,
	id uuid.UUID,
	tenantID uuid.UUID,
	preloads ...model.PaymentProductPreload,
) (*model.PaymentProduct, error) {
	paymentProduct, err := r.QueryWithPreloads(preloads...).
		Where(generated.BaseModel.ID.Eq(id)).
		Where(generated.PaymentProduct.TenantID.Eq(tenantID)).
		First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, shared.ErrPaymentProductNotFound
	}

	return &paymentProduct, err
}

func (r *paymentProductRepositoryImpl) FindByIDsIncludingDeleted(
	ctx context.Context,
	ids []uuid.UUID,
) ([]model.PaymentProduct, error) {
	items := make([]model.PaymentProduct, 0)
	if err := r.db.
		WithContext(ctx).
		Unscoped().
		Preload(string(model.PaymentProductPreloadPaymentPolicy)).
		Where("id IN ?", ids).
		Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
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
	updates.SetIfNotNil(generated.PaymentProduct.RevenueAccountCode.Column().Name, req.RevenueAccountCode)
	updates.SetIfNotNil(generated.PaymentProduct.RevenueAccountName.Column().Name, req.RevenueAccountName)
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
