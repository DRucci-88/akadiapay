package security

import (
	"akadia/model"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	Email     string         `json:"email"`
	UserID    uuid.UUID      `json:"user_id"`
	TenantID  uuid.UUID      `json:"tenant_id"`
	StudentID *uuid.UUID     `json:"student_id,omitempty"`
	RoleCode  model.RoleCode `json:"role_code"`

	jwt.RegisteredClaims
}
