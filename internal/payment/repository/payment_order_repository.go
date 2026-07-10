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

type paymentOrderRepositoryImpl struct {
	db    *gorm.DB
	query gorm.Interface[model.PaymentOrder]
}

func (r *RepositoryManagerPaymentImpl) PaymentOrder() domain.PaymentOrderRepository {
	return &paymentOrderRepositoryImpl{
		db:    r.db,
		query: gorm.G[model.PaymentOrder](r.db),
	}
}

func (r *paymentOrderRepositoryImpl) Paginate(
	ctx context.Context,
	pageable *shared.Pageable,
	chain gorm.ChainInterface[model.PaymentOrder],
) (*shared.Page[model.PaymentOrder], error) {
	total, err := chain.Count(ctx, "*")
	if err != nil {
		return shared.NewPageEmpty[model.PaymentOrder](pageable), err
	}

	items, err := chain.
		Offset(pageable.Offset()).
		Limit(pageable.Limit()).
		Find(ctx)
	if err != nil {
		return shared.NewPageEmpty[model.PaymentOrder](pageable), err
	}

	return shared.NewPage(pageable, total, items), nil
}

func (r *paymentOrderRepositoryImpl) FindPaginate(
	ctx context.Context,
	pageable *shared.Pageable,
	filter *domain.PaymentOrderFilter,
	authContext *security.AuthContext,
) (*shared.Page[model.PaymentOrder], error) {
	studentTenantSubQuery := r.db.
		Model(&model.Student{}).
		Select("id").
		Where("tenant_id = ?", authContext.TenantID)

	chain := r.query.
		Scopes().
		Where("student_id IN (?)", studentTenantSubQuery)

	switch authContext.RoleCode {
	case model.RoleCodeParent:
		parentStudentSubQuery := r.db.
			Model(&model.ParentStudent{}).
			Select("student_id").
			Where("parent_user_id = ?", authContext.UserID)

		chain = chain.Where("student_id IN (?)", parentStudentSubQuery)
	case model.RoleCodeStudent:
		if authContext.StudentID == nil {
			chain = chain.Where("1 = 0")
		} else {
			chain = chain.Where(generated.PaymentOrder.StudentID.Eq(*authContext.StudentID))
		}
	}

	if filter == nil {
		return r.Paginate(ctx, pageable, chain)
	}

	if filter.Keyword != nil {
		keyword := "%" + *filter.Keyword + "%"
		studentKeywordSubQuery := r.db.
			Model(&model.Student{}).
			Select("id").
			Where("tenant_id = ?", authContext.TenantID).
			Where("(full_name ILIKE ? OR nisn ILIKE ?)", keyword, keyword)

		chain = chain.Where(
			"(student_id IN (?) OR order_number ILIKE ? OR COALESCE(reference_number, '') ILIKE ? OR notes ILIKE ?)",
			studentKeywordSubQuery,
			keyword,
			keyword,
			keyword,
		)
	}
	if filter.StudentID != nil {
		chain = chain.
			Where(generated.PaymentOrder.StudentID.Eq(*filter.StudentID))
	}
	if filter.Status != nil {
		chain = chain.
			Where(generated.PaymentOrder.Status.Eq(string(*filter.Status)))
	}
	if filter.PaymentMethod != nil {
		chain = chain.
			Where(generated.PaymentOrder.PaymentMethod.Eq(string(*filter.PaymentMethod)))
	}
	if filter.OrderDateFrom != nil {
		chain = chain.
			Where(generated.PaymentOrder.OrderDate.Gte(*filter.OrderDateFrom))
	}
	if filter.OrderDateTo != nil {
		chain = chain.
			Where(generated.PaymentOrder.OrderDate.Lte(*filter.OrderDateTo))
	}

	return r.Paginate(ctx, pageable, chain)
}

func (r *paymentOrderRepositoryImpl) Create(
	ctx context.Context,
	paymentOrder *model.PaymentOrder,
) error {
	return r.query.Create(ctx, paymentOrder)
}

func (r *paymentOrderRepositoryImpl) FirstByID(
	ctx context.Context,
	id uuid.UUID,
) (*model.PaymentOrder, error) {
	paymentOrder, err := r.query.
		Where(generated.BaseModel.ID.Eq(id)).
		First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrPaymentOrderNotFound
		}
		return nil, err
	}

	return &paymentOrder, nil
}

func (r *paymentOrderRepositoryImpl) UpdateStatus(
	ctx context.Context,
	id uuid.UUID,
	status model.PaymentOrderStatus,
) (int, error) {
	rows, err := gorm.G[map[string]any](r.db).
		Where(generated.BaseModel.ID.Eq(id)).
		Select("*").
		Omit("id", "created_at").
		Updates(ctx, map[string]any{
			generated.PaymentOrder.Status.Column().Name: status,
		})
	if rows == 0 {
		return rows, shared.ErrPaymentOrderNotFound
	}

	return rows, err
}
