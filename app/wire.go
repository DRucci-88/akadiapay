//go:build wireinject
// +build wireinject

package app

import (
	"akadia/internal/auth"
	"akadia/internal/master/repository"
	"akadia/internal/master/service"
	"akadia/plarform/middleware"

	"github.com/google/wire"
)

func IntializedApplication() *Application {

	wire.Build(
		/* APP */
		NewApplication,
		NewDatabase,
		LoadConfig,
		NewRouter,
		NewAppConfig,

		/* PLATFORM */
		middleware.NewMiddlewareManager,

		/* AUTH */
		auth.NewAuthHandler,
		auth.NewAuthService,
		// auth.NewAuthRepositoryManagerAuth,

		/* MASTER */
		repository.NewAuthRepositoryManagerMaster,

		// Service
		service.NewStudentService,
		service.NewTenantService,
		service.NewUserService,
		service.NewUserTenantRoleService,
	)
	return nil
}
