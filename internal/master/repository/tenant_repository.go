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

type TenantRepositoryImpl struct {
	db    *gorm.DB
	query gorm.Interface[model.Tenant]
}

func (r *RepositoryManagerMasterImpl) Tenant() domain.TenantRepository {
	return &TenantRepositoryImpl{
		db:    r.db,
		query: gorm.G[model.Tenant](r.db),
	}
}

func (r *TenantRepositoryImpl) FirstByID(
	ctx context.Context,
	id uuid.UUID,
) (*model.Tenant, error) {
	tenant, err := r.query.
		Where(generated.BaseModel.ID.Eq(id)).
		Where(generated.Tenant.Status.Eq(string(model.TenantStatusActive))).
		First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, shared.ErrTenantNotFound
	}
	return &tenant, err
}
