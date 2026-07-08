package service

import (
	"akadia/domain"
	"akadia/model"
	"context"

	"github.com/google/uuid"
)

type TenantServiceImpl struct {
	repo       domain.RepositoryManagerMaster
	tenantRepo domain.TenantRepository
}

func NewTenantService(repo domain.RepositoryManagerMaster) domain.TenantService {
	return &TenantServiceImpl{
		repo:       repo,
		tenantRepo: repo.Tenant(),
	}
}

func (s *TenantServiceImpl) FirstByID(
	ctx context.Context,
	id uuid.UUID,
) (*model.Tenant, error) {
	return s.tenantRepo.FirstByID(ctx, id)
}
