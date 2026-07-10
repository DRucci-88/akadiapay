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

var paymentPolicySortableColumns = map[string]string{
	"created_at":            generated.BaseModel.CreatedAt.Column().Name,
	"updated_at":            generated.BaseModel.UpdatedAt.Column().Name,
	"code":                  generated.PaymentPolicy.Code.Column().Name,
	"name":                  generated.PaymentPolicy.Name.Column().Name,
	"minimum_amount":        generated.PaymentPolicy.MinimumAmount.Column().Name,
	"minimum_percentage":    generated.PaymentPolicy.MinimumPercentage.Column().Name,
	"allow_partial":         generated.PaymentPolicy.AllowPartial.Column().Name,
	"allow_over_payment":    generated.PaymentPolicy.AllowOverPayment.Column().Name,
	"auto_close_obligation": generated.PaymentPolicy.AutoCloseObligation.Column().Name,
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
	tenantID uuid.UUID,
) (*model.PaymentPolicy, error) {
	paymentPolicy, err := r.query.
		Where(generated.BaseModel.ID.Eq(id)).
		Where(generated.PaymentPolicy.TenantID.Eq(tenantID)).
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
	chain := r.query.
		Scopes().
		Where(generated.PaymentPolicy.TenantID.Eq(authContext.TenantID))

	if filter != nil && filter.Keyword != nil {
		chain = chain.
			Where(generated.PaymentPolicy.Code.ILike(*filter.Keyword)).
			Or(generated.PaymentPolicy.Name.ILike(*filter.Keyword))
	}
	if filter != nil && filter.AllowPartial != nil {
		chain = chain.
			Where(generated.PaymentPolicy.AllowPartial.Eq(*filter.AllowPartial))
	}
	if filter != nil && filter.AllowOverPayment != nil {
		chain = chain.
			Where(generated.PaymentPolicy.AllowOverPayment.Eq(*filter.AllowOverPayment))
	}
	if filter != nil && filter.AutoCloseObligation != nil {
		chain = chain.
			Where(generated.PaymentPolicy.AutoCloseObligation.Eq(*filter.AutoCloseObligation))
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
			paymentPolicySortableColumns,
			generated.BaseModel.CreatedAt.Column().Name,
			"DESC",
		),
	)

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
	req *domain.PaymentPolicyUpdate,
) (int, error) {
	updates := shared.UpdateMap{}

	updates.SetIfNotNil(generated.PaymentPolicy.Code.Column().Name, req.Code)
	updates.SetIfNotNil(generated.PaymentPolicy.Name.Column().Name, req.Name)
	updates.SetIfNotNil(generated.PaymentPolicy.Description.Column().Name, req.Description)
	updates.SetIfNotNil(generated.PaymentPolicy.AllowPartial.Column().Name, req.AllowPartial)
	updates.SetIfNotNil(generated.PaymentPolicy.MinimumAmount.Column().Name, req.MinimumAmount)
	updates.SetIfNotNil(generated.PaymentPolicy.MinimumPercentage.Column().Name, req.MinimumPercentage)
	updates.SetIfNotNil(generated.PaymentPolicy.AllowOverPayment.Column().Name, req.AllowOverPayment)
	updates.SetIfNotNil(generated.PaymentPolicy.AutoCloseObligation.Column().Name, req.AutoCloseObligation)

	rows, err := gorm.G[map[string]any](r.db).
		Where(generated.BaseModel.ID.Eq(id)).
		Where(generated.PaymentPolicy.TenantID.Eq(tenantID)).
		Select("*").
		Omit("id", "created_at").
		Updates(ctx, updates)
	if rows == 0 {
		return rows, shared.ErrPaymentPolicyNotFound
	}
	return rows, err
}
