package domain

import (
	"context"

	"gorm.io/gorm"
)

type RepositoryManagerMaster interface {
	WithDB(db *gorm.DB) RepositoryManagerMaster
	Transaction(
		ctx context.Context,
		fn func(repo RepositoryManagerMaster) error,
	) error
	Student() StudentRepository
	Tenant() TenantRepository
	User() UserRepository
	UserTenantRole() UserTenantRoleRepository
}
