package repository

import (
	"akadia/domain"
	"akadia/internal/platform/security"
	"akadia/internal/shared"
	"akadia/model"
	"akadia/model/generated"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type paymentOrderRepositoryImpl struct {
	db    *gorm.DB
	query gorm.Interface[model.PaymentOrder]
}

var paymentOrderSortableColumns = map[string]string{
	"created_at":       generated.BaseModel.CreatedAt.Column().Name,
	"updated_at":       generated.BaseModel.UpdatedAt.Column().Name,
	"student_id":       generated.PaymentOrder.StudentID.Column().Name,
	"order_number":     generated.PaymentOrder.OrderNumber.Column().Name,
	"order_date":       generated.PaymentOrder.OrderDate.Column().Name,
	"total_amount":     generated.PaymentOrder.TotalAmount.Column().Name,
	"status":           generated.PaymentOrder.Status.Column().Name,
	"payment_method":   generated.PaymentOrder.PaymentMethod.Column().Name,
	"ledger_posted_at": generated.PaymentOrder.LedgerPostedAt.Column().Name,
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

	if filter != nil && filter.Keyword != nil {
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
	if filter != nil && filter.StudentID != nil {
		chain = chain.
			Where(generated.PaymentOrder.StudentID.Eq(*filter.StudentID))
	}
	if filter != nil && filter.Status != nil {
		chain = chain.
			Where(generated.PaymentOrder.Status.Eq(string(*filter.Status)))
	}
	if filter != nil && filter.PaymentMethod != nil {
		chain = chain.
			Where(generated.PaymentOrder.PaymentMethod.Eq(string(*filter.PaymentMethod)))
	}
	if filter != nil && filter.OrderDateFrom != nil {
		chain = chain.
			Where(generated.PaymentOrder.OrderDate.Gte(*filter.OrderDateFrom))
	}
	if filter != nil && filter.OrderDateTo != nil {
		chain = chain.
			Where(generated.PaymentOrder.OrderDate.Lte(*filter.OrderDateTo))
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
			paymentOrderSortableColumns,
			generated.BaseModel.CreatedAt.Column().Name,
			"DESC",
		),
	)

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
	tenantID uuid.UUID,
) (*model.PaymentOrder, error) {
	paymentOrder, err := r.query.
		Where(generated.BaseModel.ID.Eq(id)).
		Where(generated.PaymentOrder.TenantID.Eq(tenantID)).
		First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrPaymentOrderNotFound
		}
		return nil, err
	}

	return &paymentOrder, nil
}

func (r *paymentOrderRepositoryImpl) LockByID(
	ctx context.Context,
	id uuid.UUID,
	tenantID uuid.UUID,
) (*model.PaymentOrder, error) {
	var paymentOrder model.PaymentOrder
	if err := r.db.
		WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", id).
		Where("tenant_id = ?", tenantID).
		First(&paymentOrder).Error; err != nil {
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

func (r *paymentOrderRepositoryImpl) MarkLedgerPosted(
	ctx context.Context,
	id uuid.UUID,
	ledgerPostedAt *time.Time,
) (int, error) {
	rows, err := gorm.G[map[string]any](r.db).
		Where(generated.BaseModel.ID.Eq(id)).
		Select("*").
		Omit("id", "created_at").
		Updates(ctx, map[string]any{
			generated.PaymentOrder.LedgerPostedAt.Column().Name: ledgerPostedAt,
		})
	if rows == 0 {
		return rows, shared.ErrPaymentOrderNotFound
	}

	return rows, err
}
