//go:build wireinject
// +build wireinject

package app

import "github.com/google/wire"

func IntializedApplication() *Application {
	wire.Build(
		// App
		NewApplication,
		NewDatabase,
		LoadConfig,
		NewRouter,
	)
	return nil
}
