package app

import "github.com/gin-gonic/gin"

type Application struct {
	Server *gin.Engine
	Config *AppConfig
	// CleanUp *worker.TokenCleanupWorker
}

func NewApplication(
	server *gin.Engine,
	config *AppConfig,
	// cleanup *worker.TokenCleanupWorker,
) *Application {
	return &Application{
		Server: server,
		Config: config,
		// CleanUp: cleanup,
	}
}
