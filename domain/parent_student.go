package domain

import (
	"context"

	"github.com/google/uuid"
)

type ParentStudentService interface {
	ExistsByParentUserIDAndStudentID(
		ctx context.Context,
		parentUserID uuid.UUID,
		studentID uuid.UUID,
	) (bool, error)
}

type ParentStudentRepository interface {
	ExistsByParentUserIDAndStudentID(
		ctx context.Context,
		parentUserID uuid.UUID,
		studentID uuid.UUID,
	) (bool, error)
}
