package repository

import (
	"akadia/domain"
	"akadia/model"
	"akadia/model/generated"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserTenantRoleRepositoryImpl struct {
	db    *gorm.DB
	query gorm.Interface[model.UserTenantRole]
}

func (r *RepositoryManagerMasterImpl) UserTenantRole() domain.UserTenantRoleRepository {
	return &UserTenantRoleRepositoryImpl{
		db:    r.db,
		query: gorm.G[model.UserTenantRole](r.db),
	}
}

func (r *UserTenantRoleRepositoryImpl) QueryWithPreloads(
	preloads ...model.UserTenantRolePreload,
) gorm.ChainInterface[model.UserTenantRole] {
	var chain gorm.ChainInterface[model.UserTenantRole] = r.query.Scopes()
	for _, preload := range preloads {
		chain = chain.Preload(string(preload), nil)
	}
	return chain
}

func (r *UserTenantRoleRepositoryImpl) FirstByUserIDAndTenantID(
	ctx context.Context,
	userID uuid.UUID,
	tenantID uuid.UUID,
	preloads ...model.UserTenantRolePreload,
) (*model.UserTenantRole, error) {
	userTenantRole, err := r.QueryWithPreloads(preloads...).
		Where(generated.UserTenantRole.UserID.Eq(userID)).
		Where(generated.UserTenantRole.TenantID.Eq(tenantID)).
		Where(generated.UserTenantRole.IsActive.Eq(true)).
		First(ctx)

	return &userTenantRole, err
}

func (r *UserTenantRoleRepositoryImpl) FindByUserID(
	ctx context.Context,
	userID uuid.UUID,
	preloads ...model.UserTenantRolePreload,
) ([]model.UserTenantRole, error) {
	userTenantRoleList, err := r.QueryWithPreloads(preloads...).
		Where(generated.UserTenantRole.UserID.Eq(userID)).
		Where(generated.UserTenantRole.IsActive.Eq(true)).
		Find(ctx)

	return userTenantRoleList, err
}
