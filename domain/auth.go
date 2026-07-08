package domain

import (
	"akadia/model"
	"akadia/plarform/security"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler interface {
	Login(c *gin.Context)
	Profile(c *gin.Context)
}

type AuthService interface {
	Login(
		ctx context.Context,
		req *AuthLoginRequest,
	) ([]AuthLoginResponse, error)
	Profile(
		ctx context.Context,
		authContext *security.AuthContext,
	) (*AuthProfileResponse, error)
}

type AuthLoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthLoginResponse struct {
	Token      string     `json:"token"`
	UserID     uuid.UUID  `json:"user_id"`
	RoleCode   string     `json:"role_code"`
	StudentID  *uuid.UUID `json:"student_id,omitempty"`
	TenantID   uuid.UUID  `json:"tenant_id"`
	TenantCode string     `json:"tenant_code"`
	TenantName string     `json:"tenant_name"`
	IsDefault  bool       `json:"is_default"`
}

type AuthProfileResponse struct {
	UserID      uuid.UUID      `json:"user_id"`
	RoleCode    model.RoleCode `json:"role_code"`
	Email       string         `json:"email"`
	TenantID    uuid.UUID      `json:"tenant_id"`
	TenantCode  string         `json:"tenant_code"`
	TenantName  string         `json:"tenant_name"`
	FullName    *string        `json:"full_name,omitempty"`
	DisplayName string         `json:"display_name"`
}
