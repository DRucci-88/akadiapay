package service

import (
	"akadia/domain"
	"akadia/model"
	"context"

	"github.com/google/uuid"
)

type UserTenantRoleServiceImpl struct {
	repo               domain.RepositoryManagerMaster
	userTenantRoleRepo domain.UserTenantRoleRepository
}

func NewUserTenantRoleService(repo domain.RepositoryManagerMaster) domain.UserTenantRoleService {
	return &UserTenantRoleServiceImpl{
		repo:               repo,
		userTenantRoleRepo: repo.UserTenantRole(),
	}
}

func (s *UserTenantRoleServiceImpl) FirstByUserIDAndTenantID(
	ctx context.Context,
	userID uuid.UUID,
	tenantID uuid.UUID,
) (*model.UserTenantRole, error) {
	return s.userTenantRoleRepo.FirstByUserIDAndTenantID(
		ctx,
		userID,
		tenantID,
		model.UserTenantRolePreloadRole,
	)
}

func (s *UserTenantRoleServiceImpl) FindByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]model.UserTenantRole, error) {
	return s.userTenantRoleRepo.FindByUserID(ctx, userID, model.UserTenantRolePreloadRole, model.UserTenantRolePreloadTenant)
}
