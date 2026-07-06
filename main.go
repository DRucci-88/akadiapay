package main

import (
	"akadia/app"
	"fmt"
)

func main() {
	application := app.IntializedApplication()

	config := application.Config

	application.Server.Run(fmt.Sprintf(":%d", config.APP_PORT))
}
