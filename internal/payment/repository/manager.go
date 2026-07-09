package repository

import (
	"akadia/domain"
	"context"

	"gorm.io/gorm"
)

type RepositoryManagerPaymentImpl struct {
	db *gorm.DB
}

func NewAuthRepositoryManagerPayment(db *gorm.DB) domain.RepositoryManagerPayment {
	return &RepositoryManagerPaymentImpl{
		db: db,
	}
}

// Usually used internally for Transaction().
func (r *RepositoryManagerPaymentImpl) WithDB(db *gorm.DB) domain.RepositoryManagerPayment {
	return &RepositoryManagerPaymentImpl{
		db: db,
	}
}

// Execute repositories inside one database transaction.
func (r *RepositoryManagerPaymentImpl) Transaction(
	ctx context.Context,
	fn func(repo domain.RepositoryManagerPayment) error,
) error {

	return r.db.
		WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {

			txRepo := r.WithDB(tx)

			return fn(txRepo)
		})
}
