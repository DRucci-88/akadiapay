package domain

import (
	"akadia/model"
	"context"

	"github.com/google/uuid"
)

type UserService interface {
	FirstByID(
		ctx context.Context,
		id uuid.UUID,
	) (*model.User, error)
	FirstByEmail(
		ctx context.Context,
		email string,
	) (*model.User, error)
}

type UserRepository interface {
	FirstByID(
		ctx context.Context,
		id uuid.UUID,
	) (*model.User, error)
	FirstByEmail(
		ctx context.Context,
		email string,
	) (*model.User, error)
}
