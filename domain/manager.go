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
	ParentStudent() ParentStudentRepository
	Tenant() TenantRepository
	User() UserRepository
	UserTenantRole() UserTenantRoleRepository
}

type RepositoryManagerPayment interface {
	WithDB(db *gorm.DB) RepositoryManagerPayment
	Transaction(
		ctx context.Context,
		fn func(repo RepositoryManagerPayment) error,
	) error
	LedgerEntry() LedgerEntryRepository
	PaymentPolicy() PaymentPolicyRepository
	PaymentProduct() PaymentProductRepository
	StudentObligation() StudentObligationRepository
	PaymentAllocation() PaymentAllocationRepository
	PaymentOrder() PaymentOrderRepository
}
