package repository

import (
	"akadia/domain"
	"akadia/internal/shared"
	"akadia/model"
	"akadia/model/generated"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ledgerEntryRepositoryImpl struct {
	db    *gorm.DB
	query gorm.Interface[model.LedgerEntry]
}

func (r *RepositoryManagerPaymentImpl) LedgerEntry() domain.LedgerEntryRepository {
	return &ledgerEntryRepositoryImpl{
		db:    r.db,
		query: gorm.G[model.LedgerEntry](r.db),
	}
}

func (r *ledgerEntryRepositoryImpl) FindByPaymentOrderID(
	ctx context.Context,
	paymentOrderID uuid.UUID,
	tenantID uuid.UUID,
) ([]model.LedgerEntry, error) {
	paymentOrderTenantSubQuery := r.db.
		Model(&model.PaymentOrder{}).
		Select("id").
		Where("tenant_id = ?", tenantID)

	return r.query.
		Where(generated.LedgerEntry.PaymentOrderID.Eq(paymentOrderID)).
		Where("payment_order_id IN (?)", paymentOrderTenantSubQuery).
		Order("created_at ASC").
		Find(ctx)
}

func (r *ledgerEntryRepositoryImpl) ExistsByPaymentOrderID(
	ctx context.Context,
	paymentOrderID uuid.UUID,
	tenantID uuid.UUID,
) (bool, error) {
	paymentOrderTenantSubQuery := r.db.
		Model(&model.PaymentOrder{}).
		Select("id").
		Where("tenant_id = ?", tenantID)

	total, err := r.query.
		Where(generated.LedgerEntry.PaymentOrderID.Eq(paymentOrderID)).
		Where("payment_order_id IN (?)", paymentOrderTenantSubQuery).
		Count(ctx, "*")
	if err != nil {
		return false, err
	}

	return total > 0, nil
}

func (r *ledgerEntryRepositoryImpl) CreateInBatches(
	ctx context.Context,
	entries []model.LedgerEntry,
	batchSize int,
) error {
	return r.db.
		WithContext(ctx).
		CreateInBatches(&entries, batchSize).
		Error
}

func (r *ledgerEntryRepositoryImpl) FindPaginate(
	ctx context.Context,
	tenantID uuid.UUID,
	pageable *shared.Pageable,
	filter *domain.LedgerEntryFilter,
) (*shared.Page[model.LedgerEntry], error) {
	paymentOrderTenantSubQuery := r.db.
		Model(&model.PaymentOrder{}).
		Select("id").
		Where("tenant_id = ?", tenantID)

	chain := r.query.
		Scopes().
		Where("payment_order_id IN (?)", paymentOrderTenantSubQuery)

	if filter != nil {
		if filter.PaymentOrderID != nil {
			chain = chain.
				Where(generated.LedgerEntry.PaymentOrderID.Eq(*filter.PaymentOrderID))
		}
		if filter.AccountCode != nil {
			chain = chain.
				Where(generated.LedgerEntry.AccountCode.Eq(*filter.AccountCode))
		}
		if filter.EntryDateFrom != nil {
			chain = chain.
				Where(generated.LedgerEntry.EntryDate.Gte(*filter.EntryDateFrom))
		}
		if filter.EntryDateTo != nil {
			chain = chain.
				Where(generated.LedgerEntry.EntryDate.Lte(*filter.EntryDateTo))
		}
		if filter.Keyword != nil {
			keyword := "%" + *filter.Keyword + "%"
			chain = chain.Where(
				"(account_code ILIKE ? OR account_name ILIKE ? OR description ILIKE ?)",
				keyword,
				keyword,
				keyword,
			)
		}
	}

	total, err := chain.Count(ctx, "*")
	if err != nil {
		return shared.NewPageEmpty[model.LedgerEntry](pageable), err
	}

	items, err := chain.
		Order("entry_date DESC, created_at DESC").
		Offset(pageable.Offset()).
		Limit(pageable.Limit()).
		Find(ctx)
	if err != nil {
		return shared.NewPageEmpty[model.LedgerEntry](pageable), err
	}

	return shared.NewPage(pageable, total, items), nil
}
