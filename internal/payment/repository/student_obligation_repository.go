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

type studentObligationRepositoryImpl struct {
	db    *gorm.DB
	query gorm.Interface[model.StudentObligation]
}

var studentObligationSortableColumns = map[string]string{
	"created_at":         generated.BaseModel.CreatedAt.Column().Name,
	"updated_at":         generated.BaseModel.UpdatedAt.Column().Name,
	"student_id":         generated.StudentObligation.StudentID.Column().Name,
	"payment_product_id": generated.StudentObligation.PaymentProductID.Column().Name,
	"period":             generated.StudentObligation.Period.Column().Name,
	"original_amount":    generated.StudentObligation.OriginalAmount.Column().Name,
	"outstanding_amount": generated.StudentObligation.OutstandingAmount.Column().Name,
	"due_date":           generated.StudentObligation.DueDate.Column().Name,
	"issued_at":          generated.StudentObligation.IssuedAt.Column().Name,
	"status":             generated.StudentObligation.Status.Column().Name,
}

func (r *RepositoryManagerPaymentImpl) StudentObligation() domain.StudentObligationRepository {
	return &studentObligationRepositoryImpl{
		db:    r.db,
		query: gorm.G[model.StudentObligation](r.db),
	}
}

func (r *studentObligationRepositoryImpl) Create(
	ctx context.Context,
	studentObligation *model.StudentObligation,
) error {
	return r.query.Create(ctx, studentObligation)
}

func (r *studentObligationRepositoryImpl) CreateBatch(
	ctx context.Context,
	studentObligations []model.StudentObligation,
) error {
	return r.db.
		WithContext(ctx).
		CreateInBatches(&studentObligations, 100).
		Error
}

func (r *studentObligationRepositoryImpl) Paginate(
	ctx context.Context,
	pageable *shared.Pageable,
	chain gorm.ChainInterface[model.StudentObligation],
) (*shared.Page[model.StudentObligation], error) {
	total, err := chain.Count(ctx, "*")
	if err != nil {
		return shared.NewPageEmpty[model.StudentObligation](pageable), err
	}

	items, err := chain.
		Offset(pageable.Offset()).
		Limit(pageable.Limit()).
		Find(ctx)
	if err != nil {
		return shared.NewPageEmpty[model.StudentObligation](pageable), err
	}

	return shared.NewPage(pageable, total, items), nil
}

func (r *studentObligationRepositoryImpl) FindPaginate(
	ctx context.Context,
	pageable *shared.Pageable,
	filter *domain.StudentObligationFilter,
	authContext *security.AuthContext,
) (*shared.Page[model.StudentObligation], error) {
	studentTenantSubQuery := r.db.
		Model(&model.Student{}).
		Select("id").
		Where("tenant_id = ?", authContext.TenantID)
	paymentProductTenantSubQuery := r.db.
		Model(&model.PaymentProduct{}).
		Select("id").
		Where("tenant_id = ?", authContext.TenantID)

	chain := r.query.
		Scopes().
		Where("student_id IN (?)", studentTenantSubQuery).
		Where("payment_product_id IN (?)", paymentProductTenantSubQuery)

	if filter != nil && filter.Keyword != nil {
		keyword := "%" + *filter.Keyword + "%"
		studentKeywordSubQuery := r.db.
			Model(&model.Student{}).
			Select("id").
			Where("tenant_id = ?", authContext.TenantID).
			Where("(full_name ILIKE ? OR nisn ILIKE ?)", keyword, keyword)
		paymentProductKeywordSubQuery := r.db.
			Model(&model.PaymentProduct{}).
			Select("id").
			Where("tenant_id = ?", authContext.TenantID).
			Where("(code ILIKE ? OR name ILIKE ?)", keyword, keyword)

		chain = chain.Where(
			"(student_id IN (?) OR payment_product_id IN (?))",
			studentKeywordSubQuery,
			paymentProductKeywordSubQuery,
		)
	}
	if filter != nil && filter.StudentID != nil {
		chain = chain.
			Where(generated.StudentObligation.StudentID.Eq(*filter.StudentID))
	}
	if filter != nil && filter.PaymentProductID != nil {
		chain = chain.
			Where(generated.StudentObligation.PaymentProductID.Eq(*filter.PaymentProductID))
	}
	if filter != nil && filter.Status != nil {
		chain = chain.
			Where(generated.StudentObligation.Status.Eq(string(*filter.Status)))
	}
	if filter != nil && filter.DueDateFrom != nil {
		chain = chain.
			Where(generated.StudentObligation.DueDate.Gte(*filter.DueDateFrom))
	}
	if filter != nil && filter.DueDateTo != nil {
		chain = chain.
			Where(generated.StudentObligation.DueDate.Lte(*filter.DueDateTo))
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
			studentObligationSortableColumns,
			generated.BaseModel.CreatedAt.Column().Name,
			"DESC",
		),
	)

	return r.Paginate(ctx, pageable, chain)
}

func (r *studentObligationRepositoryImpl) SumOutstandingByStudentID(
	ctx context.Context,
	studentID uuid.UUID,
) (float64, error) {
	type result struct {
		Total float64
	}

	var temp result
	if err := r.db.
		WithContext(ctx).
		Model(&model.StudentObligation{}).
		Select("COALESCE(SUM(outstanding_amount), 0) AS total").
		Where("student_id = ?", studentID).
		Where("outstanding_amount > 0").
		Where("status <> ?", model.StudentObligationStatusCancelled).
		Scan(&temp).Error; err != nil {
		return 0, err
	}

	return temp.Total, nil
}

func (r *studentObligationRepositoryImpl) FindOutstandingByStudentID(
	ctx context.Context,
	studentID uuid.UUID,
) ([]model.StudentObligation, error) {
	return r.query.
		Where(generated.StudentObligation.StudentID.Eq(studentID)).
		Where(generated.StudentObligation.OutstandingAmount.Gt(0)).
		Where(generated.StudentObligation.Status.Neq(string(model.StudentObligationStatusCancelled))).
		Find(ctx)
}

func (r *studentObligationRepositoryImpl) FirstByID(
	ctx context.Context,
	id uuid.UUID,
	tenantID uuid.UUID,
) (*model.StudentObligation, error) {
	studentTenantSubQuery := r.db.
		Model(&model.Student{}).
		Select("id").
		Where("tenant_id = ?", tenantID)
	paymentProductTenantSubQuery := r.db.
		Model(&model.PaymentProduct{}).
		Select("id").
		Where("tenant_id = ?", tenantID)

	studentObligation, err := r.query.
		Where(generated.BaseModel.ID.Eq(id)).
		Where("student_id IN (?)", studentTenantSubQuery).
		Where("payment_product_id IN (?)", paymentProductTenantSubQuery).
		First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrStudentObligationNotFound
		}
		return nil, err
	}

	return &studentObligation, nil
}

func (r *studentObligationRepositoryImpl) LockByIDs(
	ctx context.Context,
	ids []uuid.UUID,
) ([]model.StudentObligation, error) {
	items := make([]model.StudentObligation, 0)
	if err := r.db.
		WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id IN ?", ids).
		Order("id ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (r *studentObligationRepositoryImpl) Update(
	ctx context.Context,
	id uuid.UUID,
	req *domain.StudentObligationUpdate,
) (int, error) {
	updates := shared.UpdateMap{}
	updates.SetIfNotNil(generated.StudentObligation.DueDate.Column().Name, req.DueDate)
	updates.SetIfNotNil(generated.StudentObligation.Notes.Column().Name, req.Notes)

	rows, err := gorm.G[map[string]any](r.db).
		Where(generated.BaseModel.ID.Eq(id)).
		Select("*").
		Omit("id", "created_at").
		Updates(ctx, updates)
	if rows == 0 {
		return rows, shared.ErrStudentObligationNotFound
	}

	return rows, err
}

func (r *studentObligationRepositoryImpl) Delete(
	ctx context.Context,
	id uuid.UUID,
) (int, error) {
	result := r.db.
		WithContext(ctx).
		Where("id = ?", id).
		Delete(&model.StudentObligation{})
	if result.RowsAffected == 0 {
		return int(result.RowsAffected), shared.ErrStudentObligationNotFound
	}

	return int(result.RowsAffected), result.Error
}

func (r *studentObligationRepositoryImpl) HasPaymentAllocations(
	ctx context.Context,
	id uuid.UUID,
) (bool, error) {
	total, err := gorm.G[model.PaymentAllocation](r.db).
		Where(generated.PaymentAllocation.StudentObligationID.Eq(id)).
		Count(ctx, "*")
	if err != nil {
		return false, err
	}

	return total > 0, nil
}

func (r *studentObligationRepositoryImpl) FindByIDs(
	ctx context.Context,
	ids []uuid.UUID,
) ([]model.StudentObligation, error) {
	return r.query.
		Where("id IN ?", ids).
		Find(ctx)
}

func (r *studentObligationRepositoryImpl) UpdateSettlement(
	ctx context.Context,
	id uuid.UUID,
	outstandingAmount float64,
	status model.StudentObligationStatus,
) (int, error) {
	rows, err := gorm.G[map[string]any](r.db).
		Where(generated.BaseModel.ID.Eq(id)).
		Select("*").
		Omit("id", "created_at").
		Updates(ctx, map[string]any{
			generated.StudentObligation.OutstandingAmount.Column().Name: outstandingAmount,
			generated.StudentObligation.Status.Column().Name:            status,
		})
	if rows == 0 {
		return rows, shared.ErrStudentObligationNotFound
	}

	return rows, err
}

func (r *studentObligationRepositoryImpl) ExistsActiveByStudentIDAndPaymentProductIDAndPeriod(
	ctx context.Context,
	studentID uuid.UUID,
	paymentProductID uuid.UUID,
	period time.Time,
) (bool, error) {
	total, err := r.query.
		Where(generated.StudentObligation.StudentID.Eq(studentID)).
		Where(generated.StudentObligation.PaymentProductID.Eq(paymentProductID)).
		Where(generated.StudentObligation.Period.Eq(period)).
		Where(generated.StudentObligation.Status.Neq(string(model.StudentObligationStatusCancelled))).
		Count(ctx, "*")
	if err != nil {
		return false, err
	}

	return total > 0, nil
}
