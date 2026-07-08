package app

import (
	"akadia/domain"

	"github.com/gin-gonic/gin"
)

type Application struct {
	Server *gin.Engine
	Config domain.AppConfigProvider
	// CleanUp *worker.TokenCleanupWorker
}

func NewApplication(
	server *gin.Engine,
	config domain.AppConfigProvider,
	// cleanup *worker.TokenCleanupWorker,
) *Application {
	return &Application{
		Server: server,
		Config: config,
		// CleanUp: cleanup,
	}
}
