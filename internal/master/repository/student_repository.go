package repository

import (
	"akadia/domain"
	"akadia/internal/shared"
	"akadia/model"
	"akadia/model/generated"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StudentRepositoryImpl struct {
	db    *gorm.DB
	query gorm.Interface[model.Student]
}

func (r *RepositoryManagerMasterImpl) Student() domain.StudentRepository {
	return &StudentRepositoryImpl{
		db:    r.db,
		query: gorm.G[model.Student](r.db),
	}
}

func (r *StudentRepositoryImpl) QueryWithPreloads(
	preloads ...model.StudentPreload,
) gorm.ChainInterface[model.Student] {
	var chain gorm.ChainInterface[model.Student] = r.query.Scopes()
	for _, preload := range preloads {
		chain = chain.Preload(string(preload), nil)
	}
	return chain
}

func (r *StudentRepositoryImpl) FirstByID(
	ctx context.Context,
	id uuid.UUID,
	preloads ...model.StudentPreload,
) (*model.Student, error) {
	student, err := r.QueryWithPreloads(preloads...).
		Where(generated.BaseModel.ID.Eq(id)).
		Where(generated.Student.Status.Eq(string(model.StudentStatusActive))).
		First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, shared.ErrStudentNotFound
	}
	return &student, err
}

func (r *StudentRepositoryImpl) FindByUserID(
	ctx context.Context,
	userID uuid.UUID,
	preloads ...model.StudentPreload,
) (*model.Student, error) {
	student, err := r.QueryWithPreloads(preloads...).
		Where(generated.Student.UserID.Eq(userID)).
		Where(generated.Student.Status.Eq(string(model.StudentStatusActive))).
		First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, shared.ErrStudentNotFound
	}
	return &student, err
}
