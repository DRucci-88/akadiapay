package domain

import (
	"akadia/model"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserTenantRoleService interface {
	FirstByUserIDAndTenantID(
		ctx context.Context,
		userID uuid.UUID,
		tenantID uuid.UUID,
	) (*model.UserTenantRole, error)
	FindByUserID(
		ctx context.Context,
		userID uuid.UUID,
	) ([]model.UserTenantRole, error)
}

type UserTenantRoleRepository interface {
	QueryWithPreloads(
		preloads ...model.UserTenantRolePreload,
	) gorm.ChainInterface[model.UserTenantRole]
	FirstByUserIDAndTenantID(
		ctx context.Context,
		userID uuid.UUID,
		tenantID uuid.UUID,
		preloads ...model.UserTenantRolePreload,
	) (*model.UserTenantRole, error)
	FindByUserID(
		ctx context.Context,
		userID uuid.UUID,
		preloads ...model.UserTenantRolePreload,
	) ([]model.UserTenantRole, error)
}
