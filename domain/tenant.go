package domain

import (
	"akadia/model"
	"context"

	"github.com/google/uuid"
)

type TenantService interface {
	FirstByID(
		ctx context.Context,
		id uuid.UUID,
	) (*model.Tenant, error)
}

type TenantRepository interface {
	FirstByID(
		ctx context.Context,
		id uuid.UUID,
	) (*model.Tenant, error)
}
