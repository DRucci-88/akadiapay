package service

import (
	"akadia/domain"
	"akadia/model"
	"context"

	"github.com/google/uuid"
)

type UserServiceImpl struct {
	repo     domain.RepositoryManagerMaster
	userRepo domain.UserRepository
}

func NewUserService(repo domain.RepositoryManagerMaster) domain.UserService {
	return &UserServiceImpl{
		repo:     repo,
		userRepo: repo.User(),
	}
}

func (s *UserServiceImpl) FirstByID(
	ctx context.Context,
	id uuid.UUID,
) (*model.User, error) {
	return s.userRepo.FirstByID(ctx, id)
}

func (s *UserServiceImpl) FirstByEmail(
	ctx context.Context,
	email string,
) (*model.User, error) {
	return s.userRepo.FirstByEmail(ctx, email)
}
