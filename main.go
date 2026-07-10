package main

import (
	"akadia/app"
	"fmt"
)

// @title AkadiaPay API
// @version 1.0
// @description AkadiaPay is a school financial platform API for authentication, payment configuration, student billing, payment processing, allocation, and financial ledger.
// @description
// @description Business process coverage:
// @description - Authentication and tenant workspace access
// @description - Payment policy configuration
// @description - Payment product configuration
// @description - Student obligation billing
// @description - Parent/student outstanding bill visibility
// @description - Payment order processing
// @description - Payment allocation
// @description - Financial ledger posting and reporting
// @termsOfService http://swagger.io/terms/
// @contact.name AkadiaPay API Support
// @contact.email support@akadiapay.local
// @license.name MIT
// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token. Example: "Bearer eyJhbGciOiJIUzI1NiIs..."
func main() {
	application := app.IntializedApplication()

	config := application.Config

	application.Server.Run(fmt.Sprintf(":%d", config.APP_PORT()))
}
