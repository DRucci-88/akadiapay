package auth

import (
	"context"

	"gorm.io/gorm"
)

type RepositoryManagerAuth interface {
	WithDB(db *gorm.DB) *repositoryManagerAuthImpl
	Transaction(
		ctx context.Context,
		fn func(repo *repositoryManagerAuthImpl) error,
	) error
}

type repositoryManagerAuthImpl struct {
	db *gorm.DB
}

// Constructor (Google Wire)
func NewAuthRepositoryManagerAuth(db *gorm.DB) *repositoryManagerAuthImpl {
	return &repositoryManagerAuthImpl{
		db: db,
	}
}

// Usually used internally for Transaction().
func (r *repositoryManagerAuthImpl) WithDB(db *gorm.DB) *repositoryManagerAuthImpl {
	return &repositoryManagerAuthImpl{
		db: db,
	}
}

// Execute repositories inside one database transaction.
func (r *repositoryManagerAuthImpl) Transaction(
	ctx context.Context,
	fn func(repo *repositoryManagerAuthImpl) error,
) error {

	return r.db.
		WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {

			txRepo := r.WithDB(tx)

			return fn(txRepo)
		})
}
