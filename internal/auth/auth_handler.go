package auth

import (
	"akadia/domain"
	"akadia/internal/shared"
	"akadia/plarform/security"
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

func (h *authHandlerImpl) Profile(c *gin.Context) {
	authContextValue, exist := c.Get("auth")

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
