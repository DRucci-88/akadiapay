package repository

import (
	"akadia/domain"
	"context"

	"gorm.io/gorm"
)

type RepositoryManagerMasterImpl struct {
	db *gorm.DB
}

func NewAuthRepositoryManagerMaster(db *gorm.DB) domain.RepositoryManagerMaster {
	return &RepositoryManagerMasterImpl{
		db: db,
	}
}

// Usually used internally for Transaction().
func (r *RepositoryManagerMasterImpl) WithDB(db *gorm.DB) domain.RepositoryManagerMaster {
	return &RepositoryManagerMasterImpl{
		db: db,
	}
}

// Execute repositories inside one database transaction.
func (r *RepositoryManagerMasterImpl) Transaction(
	ctx context.Context,
	fn func(repo domain.RepositoryManagerMaster) error,
) error {

	return r.db.
		WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {

			txRepo := r.WithDB(tx)

			return fn(txRepo)
		})
}
