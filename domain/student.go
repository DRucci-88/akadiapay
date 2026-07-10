package domain

import (
	"akadia/model"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StudentService interface {
	FirstByID(
		ctx context.Context,
		id uuid.UUID,
	) (*model.Student, error)
	FindByUserID(
		ctx context.Context,
		userID uuid.UUID,
	) (*model.Student, error)
}

type StudentRepository interface {
	QueryWithPreloads(
		preloads ...model.StudentPreload,
	) gorm.ChainInterface[model.Student]
	FirstByID(
		ctx context.Context,
		id uuid.UUID,
		preloads ...model.StudentPreload,
	) (*model.Student, error)
	FindByUserID(
		ctx context.Context,
		userID uuid.UUID,
		preloads ...model.StudentPreload,
	) (*model.Student, error)
}
