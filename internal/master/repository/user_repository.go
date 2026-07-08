package repository

import (
	"akadia/domain"
	"akadia/internal/shared"
	"akadia/model"
	"akadia/model/generated"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db    *gorm.DB
	query gorm.Interface[model.User]
}

func (r *RepositoryManagerMasterImpl) User() domain.UserRepository {
	return &UserRepositoryImpl{
		db:    r.db,
		query: gorm.G[model.User](r.db),
	}
}

func (r *UserRepositoryImpl) FirstByID(
	ctx context.Context,
	id uuid.UUID,
) (*model.User, error) {
	user, err := r.query.
		Where(generated.BaseModel.ID.Eq(id)).
		Where(generated.User.Status.Eq(model.UserStatusActive)).
		First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, shared.ErrUserNotFound
	}
	return &user, err
}

func (r *UserRepositoryImpl) FirstByEmail(
	ctx context.Context,
	email string,
) (*model.User, error) {
	user, err := r.query.
		Where(generated.User.Email.Eq(email)).
		Where(generated.User.Status.Eq(model.UserStatusActive)).
		First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, shared.ErrUserNotFound
	}
	return &user, err
}
