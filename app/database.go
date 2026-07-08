package app

import (
	"akadia/domain"
	"akadia/model"
	"akadia/seed"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(
	appConfig domain.AppConfigProvider,
) *gorm.DB {
	// dsn := "postgres://postgres:12345678@localhost:5432/ujian2_rematch?sslmode=disable"
	dsn := appConfig.DB_DSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("Database Not Connected " + err.Error())
	}

	schemas := []string{
		model.SchemaMaster,
		model.SchemaPayment,
	}

	for _, schema := range schemas {
		if err := db.Exec(
			fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schema),
		).Error; err != nil {
			panic("Schema Failed Create" + err.Error())
		}
	}

	if err := db.AutoMigrate(
	// &model.User{},
	// &model.Attendance{},
	// &model.Department{},
	// &model.Employee{},
	// &model.Leave{},
	// &model.Position{},
	// &model.Salary{},
	// &model.BlackListedToken{},
	); err != nil {
		panic("Auto Migrate Failed " + err.Error())
	}

	if err := seed.Run(db); err != nil {
		panic("Seeder FAILED" + err.Error())
	}

	return db
}
