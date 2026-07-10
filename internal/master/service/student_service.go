package service

import (
	"akadia/domain"
	"akadia/model"
	"context"

	"github.com/google/uuid"
)

type StudentServiceImpl struct {
	repo        domain.RepositoryManagerMaster
	studentRepo domain.StudentRepository
}

func NewStudentService(repo domain.RepositoryManagerMaster) domain.StudentService {
	return &StudentServiceImpl{
		repo:        repo,
		studentRepo: repo.Student(),
	}
}

func (s *StudentServiceImpl) FirstByID(
	ctx context.Context,
	id uuid.UUID,
) (*model.Student, error) {
	return s.studentRepo.FirstByID(ctx, id, model.StudentPreloadTenant)
}

func (s *StudentServiceImpl) FindByUserID(
	ctx context.Context,
	userID uuid.UUID,
) (*model.Student, error) {
	return s.studentRepo.FindByUserID(ctx, userID, model.StudentPreloadUser, model.StudentPreloadTenant)
}
