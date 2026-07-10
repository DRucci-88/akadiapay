package service

import (
	"akadia/domain"
	"context"

	"github.com/google/uuid"
)

type ParentStudentServiceImpl struct {
	repo              domain.RepositoryManagerMaster
	parentStudentRepo domain.ParentStudentRepository
}

func NewParentStudentService(repo domain.RepositoryManagerMaster) domain.ParentStudentService {
	return &ParentStudentServiceImpl{
		repo:              repo,
		parentStudentRepo: repo.ParentStudent(),
	}
}

func (s *ParentStudentServiceImpl) ExistsByParentUserIDAndStudentID(
	ctx context.Context,
	parentUserID uuid.UUID,
	studentID uuid.UUID,
) (bool, error) {
	return s.parentStudentRepo.ExistsByParentUserIDAndStudentID(ctx, parentUserID, studentID)
}
