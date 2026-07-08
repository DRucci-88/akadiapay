package security

import (
	"akadia/model"
	"time"

	"github.com/google/uuid"
)

type AuthContext struct {
	Email          string         `json:"email"`
	UserID         uuid.UUID      `json:"user_id"`
	TenantID       uuid.UUID      `json:"tenant_id"`
	StudentID      *uuid.UUID     `json:"student_id,omitempty"`
	RoleCode       model.RoleCode `json:"role_code"`
	Token          string         `json:"token"`
	TokenExpiredAt time.Time      `json:"token_expired_at"`
}
