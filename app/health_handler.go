package app

import (
	"akadia/model/generated"
	"net/http"

	"github.com/gin-gonic/gin"
)

// health godoc
// @Summary Health check
// @Description Returns the current application health payload and static metadata exposed by the existing health endpoint.
// @Tags Health
// @Produce json
// @Success 200 {object} domain.SwaggerHealthResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /health [get]
func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "Berjalan Perfecto",
		"aaa":    generated.PaymentPolicy.Code.Column().Name,
	})
}
