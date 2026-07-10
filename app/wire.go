//go:build wireinject
// +build wireinject

package app

import (
	"akadia/internal/auth"
	repoMaster "akadia/internal/master/repository"
	serviceMaster "akadia/internal/master/service"
	handlerPayment "akadia/internal/payment/handler"
	repoPayment "akadia/internal/payment/repository"
	servicePayment "akadia/internal/payment/service"
	"akadia/internal/platform/middleware"

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
		repoMaster.NewAuthRepositoryManagerMaster,

		// Service
		serviceMaster.NewParentStudentService,
		serviceMaster.NewStudentService,
		serviceMaster.NewTenantService,
		serviceMaster.NewUserService,
		serviceMaster.NewUserTenantRoleService,

		/* PAYMENT */
		repoPayment.NewAuthRepositoryManagerPayment,

		// Handler
		handlerPayment.NewPaymentAllocationHandler,
		handlerPayment.NewPaymentOrderHandler,
		handlerPayment.NewPaymentProductHandler,
		handlerPayment.NewPaymentPolicyHandler,
		handlerPayment.NewStudentObligationHandler,

		// Service
		servicePayment.NewPaymentAllocationService,
		servicePayment.NewPaymentOrderService,
		servicePayment.NewPaymentPolicyService,
		servicePayment.NewPaymentProductService,
		servicePayment.NewStudentObligationService,
	)
	return nil
}
