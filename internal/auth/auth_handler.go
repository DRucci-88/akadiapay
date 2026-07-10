package auth

import (
	"akadia/domain"
	"akadia/internal/platform/security"
	"akadia/internal/shared"
	"net/http"

	"github.com/gin-gonic/gin"
)

type authHandlerImpl struct {
	authService domain.AuthService
}

func NewAuthHandler(authService domain.AuthService) domain.AuthHandler {
	return &authHandlerImpl{
		authService: authService,
	}
}

// Login godoc
// @Summary Login and get tenant access tokens
// @Description Authenticates a user and returns JWT-based tenant and role access contexts for the workspaces available to that user.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body domain.AuthLoginRequest true "Login request"
// @Success 200 {object} domain.SwaggerAuthLoginResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 422 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /auth/login [post]
func (h *authHandlerImpl) Login(c *gin.Context) {
	var req domain.AuthLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	res, err := h.authService.Login(c.Request.Context(), &req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

// Profile godoc
// @Summary Get authenticated profile
// @Description Returns the active JWT-derived user, role, tenant, and optional student context for the current access token.
// @Tags Auth
// @Produce json
// @Success 200 {object} domain.SwaggerAuthProfileResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Security BearerAuth
// @Router /auth/profile [get]
func (h *authHandlerImpl) Profile(c *gin.Context) {
	authContextValue, exist := c.Get(domain.ContextKeyAuth)

	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": shared.ErrAuthUnauthorized,
		})
		return
	}

	authContext := authContextValue.(*security.AuthContext)

	res, err := h.authService.Profile(c.Request.Context(), authContext)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}
