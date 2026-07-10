package repository

import (
	"akadia/domain"
	"akadia/internal/platform/security"
	"akadia/internal/shared"
	"akadia/model"
	"akadia/model/generated"
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type studentObligationRepositoryImpl struct {
	db    *gorm.DB
	query gorm.Interface[model.StudentObligation]
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
	if filter.StudentID != nil {
		chain = chain.
			Where(generated.StudentObligation.StudentID.Eq(*filter.StudentID))
	}
	if filter.PaymentProductID != nil {
		chain = chain.
			Where(generated.StudentObligation.PaymentProductID.Eq(*filter.PaymentProductID))
	}
	if filter.Status != nil {
		chain = chain.
			Where(generated.StudentObligation.Status.Eq(string(*filter.Status)))
	}
	if filter.DueDateFrom != nil {
		chain = chain.
			Where(generated.StudentObligation.DueDate.Gte(*filter.DueDateFrom))
	}
	if filter.DueDateTo != nil {
		chain = chain.
			Where(generated.StudentObligation.DueDate.Lte(*filter.DueDateTo))
	}

	return r.Paginate(ctx, pageable, chain)
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
