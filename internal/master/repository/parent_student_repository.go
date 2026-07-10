package repository

import (
	"akadia/domain"
	"akadia/model"
	"akadia/model/generated"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ParentStudentRepositoryImpl struct {
	db    *gorm.DB
	query gorm.Interface[model.ParentStudent]
}

func (r *RepositoryManagerMasterImpl) ParentStudent() domain.ParentStudentRepository {
	return &ParentStudentRepositoryImpl{
		db:    r.db,
		query: gorm.G[model.ParentStudent](r.db),
	}
}

func (r *ParentStudentRepositoryImpl) ExistsByParentUserIDAndStudentID(
	ctx context.Context,
	parentUserID uuid.UUID,
	studentID uuid.UUID,
) (bool, error) {
	total, err := r.query.
		Where(generated.ParentStudent.ParentUserID.Eq(parentUserID)).
		Where(generated.ParentStudent.StudentID.Eq(studentID)).
		Count(ctx, "*")
	if err != nil {
		return false, err
	}

	return total > 0, nil
}
